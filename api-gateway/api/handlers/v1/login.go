package v1

import (
	_ "api-gateway/api/docs"
	"api-gateway/api/models"
	pb "api-gateway/genproto/user_service"
	regtool "api-gateway/internal/pkg/regtool"
	tokens "api-gateway/internal/pkg/token"
	validation "api-gateway/internal/pkg/validation"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	govalidator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Summary 		Register
// @Description 	Api for register user
// @Tags 			registration
// @Accept 			json
// @Produce 		json
// @Param 			User body models.UserRegister true "Register User"
// @Success 		201 {object} models.User
// @Failure 		400 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/register [POST]
func (h *HandlerV1) Register(c *gin.Context) {
	var (
		body        models.UserRegister
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	valid := govalidator.IsEmail(body.Email)
	if !valid {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Bad email",
		})
		log.Println(err)
		return
	}

	body.Email, err = validation.EmailValidation(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	status := validation.PasswordValidation(body.Password)
	if !status {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Password should be 8-20 characters long and contain at least one lowercase letter, one uppercase letter, and one digit",
		})
		log.Println(err)
		return
	}

	exists, err := h.Service.UserService().UniqueEmail(ctx, &pb.IsUnique{
		Email: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}
	if exists.Status {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "This email already in use:",
		})
		return
	}

	hashPassword, err := regtool.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	createdUser, err := h.Service.UserService().CreateUser(ctx, &pb.User{
		Id:        uuid.NewString(),
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  hashPassword,
		Gender:    body.Gender,
		Role:      "user",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	h.RefreshToken = tokens.JWTHandler{
		Sub:        createdUser.Guid,
		Role:       "user",
		SigningKey: h.Config.Token.SignInKey,
		Log:        h.Logger,
		Email:      body.Email,
	}

	access, refresh, err := h.RefreshToken.GenerateAuthJWT()

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	responseStatus, err := h.Service.UserService().UpdateRefresh(ctx, &pb.RefreshRequest{
		UserId:       createdUser.Guid,
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	if !responseStatus.Status {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: "Error happened",
		})
		log.Println("Server error updating refresh in register")
		return
	}

	c.JSON(http.StatusCreated, models.User{
		Id:        createdUser.Guid,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  hashPassword,
		Gender:    body.Gender,
		Role:      "user",
		Refresh:   refresh,
		Access:    access,
	})
}

// @Summary 		Login
// @Description 	Api for user user
// @Tags 			registration
// @Accept 			json
// @Produce 		json
// @Param 			login body models.Login true "Login Model"
// @Success 		200 {object} models.User
// @Failure 		400 {object} models.Error
// @Failure 		404 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/login [POST]
func (h *HandlerV1) Login(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var body models.Login

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	filter := map[string]string{
		"email": body.Email,
	}

	response, err := h.Service.UserService().GetUser(
		ctx, &pb.Filter{
			Filter: filter,
		})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	if !(regtool.CheckHashPassword(body.Password, response.Password)) {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Incorrect Password",
		})
		return
	}

	h.RefreshToken = tokens.JWTHandler{
		Sub:        response.Id,
		Role:       response.Role,
		SigningKey: h.Config.Token.SignInKey,
		Log:        h.Logger,
		Email:      response.Email,
	}

	access, refresh, err := h.RefreshToken.GenerateAuthJWT()

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	respUser := &models.User{
		Id:          response.Id,
		FirstName:   response.FirstName,
		LastName:    response.LastName,
		Email:       response.Email,
		Password:    response.Password,
		PhoneNumber: response.PhoneNumber,
		Gender:      response.Gender,
		Age:         response.Age,
		Role:        response.Role,
		Refresh:     refresh,
		Access:      access,
	}
	_, err = h.Service.UserService().UpdateRefresh(ctx, &pb.RefreshRequest{
		UserId:       response.Id,
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, respUser)
}

// @Summary 		Forget Password
// @Description 	Api for sending otp
// @Tags 			registration
// @Accept 			json
// @Produce 		json
// @Param 			email path string true "Email"
// @Success 		200 {object} string
// @Failure 		400 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/forgot/{email} [POST]
func (h *HandlerV1) Forgot(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	email := c.Param("email")

	email, err = validation.EmailValidation(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	status, err := h.Service.UserService().UniqueEmail(ctx, &pb.IsUnique{
		Email: email,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	if !status.Status {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "This user is not registered",
		})
		return
	}

	radomNumber, err := regtool.SendCodeGmail(email, "ClothesStore\n", "./internal/pkg/regtool/forgotpassword.html", h.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	if err := h.redisStorage.Set(ctx, radomNumber, cast.ToString(email), time.Second*300); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, "We have sent otp your email")
}

// @Summary 		Verify OTP
// @Description 	Api for verify user
// @Tags 			registration
// @Accept 			json
// @Produce 		json
// @Param 			email query string true "Email"
// @Param 			otp query string true "OTP"
// @Success 		200 {object} bool
// @Failure 		400 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/verify [POST]
func (h *HandlerV1) Verify(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	otp := c.Query("otp")
	email := c.Query("email")

	userData, err := h.redisStorage.Get(ctx, otp)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err)
		return
	}
	var redisEmail string

	err = json.Unmarshal(userData, &redisEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err)
		return
	}

	if redisEmail != email {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "The email did not match",
		})
		log.Println("The email did not match")
		return
	}

	c.JSON(http.StatusCreated, true)
}

// @Summary 		Reset Password
// @Description 	Api for reset password
// @Tags 			registration
// @Accept 			json
// @Produce 		json
// @Param 			User body models.ResetPassword true "Reset Password"
// @Success 		200 {object} bool
// @Failure 		400 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/reset-password [PUT]
func (h *HandlerV1) ResetPassword(c *gin.Context) {
	var (
		body models.ResetPassword
	)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(7))
	defer cancel()

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	status := validation.PasswordValidation(body.NewPassword)
	if !status {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Password should be 8-20 characters long and contain at least one lowercase letter, one uppercase letter, one symbol, and one digit",
		})
		log.Println(err)
		return
	}

	user, err := h.Service.UserService().GetUser(ctx, &pb.Filter{
		Filter: map[string]string{
			"email": body.Email,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	hashPassword, err := regtool.HashPassword(body.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadGateway, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err)
		return
	}

	responseStatus, err := h.Service.UserService().UpdatePassword(ctx, &pb.UpdatePasswordRequest{
		UserId:      user.Id,
		NewPassword: hashPassword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	if !responseStatus.Status {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: "Password doesn't updated",
		})
		log.Println("Password doesn't updated")
		return
	}

	c.JSON(http.StatusCreated, true)
}

// @Summary 		New Token
// @Description 	Api for updated acces token
// @Tags 			registration
// @Accept 			json
// @Produce 		json
// @Param 			refresh path string true "Refresh Token"
// @Success 		200 {object} models.TokenResp
// @Failure 		400 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		409 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/token/{refresh} [GET]
func (h *HandlerV1) Token(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	RToken := c.Param("refresh")

	user, err := h.Service.UserService().GetUser(ctx, &pb.Filter{
		Filter: map[string]string{
			"refresh": RToken,
		},
	})

	if err != nil {
		c.JSON(500, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err)
		return
	}

	resclaim, err := tokens.ExtractClaim(RToken, []byte(h.Config.Token.SignInKey))
	if err != nil {
		c.JSON(500, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}
	Now_time := time.Now().Unix()
	exp := (resclaim["exp"])
	if exp.(float64)-float64(Now_time) > 0 {
		h.RefreshToken = tokens.JWTHandler{
			Sub:        user.Id,
			Role:       user.Role,
			SigningKey: h.Config.Token.SignInKey,
			Log:        h.Logger,
			Email:      user.Email,
		}

		access, refresh, err := h.RefreshToken.GenerateAuthJWT()
		if err != nil {
			c.JSON(http.StatusConflict, models.Error{
				Message: err.Error(),
			})
			log.Println(err)
			return
		}

		_, err = h.Service.UserService().UpdateRefresh(ctx, &pb.RefreshRequest{
			UserId:       user.Id,
			RefreshToken: refresh,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Message: err.Error(),
			})
			log.Println(err)
			return
		}

		respUser := &models.TokenResp{
			ID:      user.Id,
			Role:    user.Role,
			Refresh: refresh,
			Access:  access,
		}

		c.JSON(http.StatusCreated, respUser)
	} else {
		c.JSON(http.StatusUnauthorized, models.Error{
			Message: "refresh token expired",
		})
		log.Println("refresh token expired")
		return
	}
}
