package v1

import (
	_ "api-gateway/api/docs"
	"api-gateway/api/models"
	pb "api-gateway/genproto/user_service"
	regtool "api-gateway/internal/pkg/regtool"
	validation "api-gateway/internal/pkg/validation"
	"api-gateway/internal/usecase/refresh_token"
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

// register register users
// @Summary RegisterUser
// @Description Api for register user
// @Tags registration
// @Accept json
// @Produce json
// @Param User body models.UserCreateReq true "RegisterUser"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users/register [post]
func (h *HandlerV1) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	var (
		body        models.UserCreateReq
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}
	valid := govalidator.IsEmail(body.Email)
	if !valid {
		c.JSON(http.StatusBadRequest, models.Error{
			Message : "Bad email",
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

	IsPhoneNumber := validation.PhoneUz(body.PhoneNumber)
	if !IsPhoneNumber {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Incorrect telephone number!!!",
		})
		log.Println(err)
		return
	}

	IsName := validation.NameValiddation(body.FirstName)
	if !IsName {
		c.JSON(http.StatusBadRequest, models.Error{
			Message : "The name must contain only letters",
		})
		log.Println("The name must contain only letters")
		return
	}

	IsSurname := validation.NameValiddation(body.LastName)
	if !IsSurname {
		c.JSON(http.StatusBadRequest, models.Error{
			Message : "The name must contain only letters",
		})
		log.Println("The name must contain only letters")
		return
	}

	gender := validation.GenderValidation(body.Gender)
	if !gender {
		c.JSON(http.StatusBadRequest, models.Error{
			Message : "Gender must contain only male or female",
		})
		log.Println("Gender must contain only male or female")
		return
	}

	exists, err := h.Service.UserService().UniqueEmail(ctx, &pb.IsUnique{
		Email: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message : err.Error(),
		})
		log.Println(err)
		return
	}

	if exists.Status {
		c.JSON(http.StatusConflict, models.Error{
			Message : "This email already in use:",
		})
		return
	}

	radomNumber, err := regtool.SendCodeGmail(body.Email, "ClothesStore\n", "./internal/pkg/regtool/emailotp.html", h.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return 
	}

	err = h.redisStorage.Set(ctx, radomNumber, body, time.Second*300)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message : err.Error(),
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Response: "We have sent otp your email",
	})
}

// register verify users
// @Summary RegisterUser
// @Description Api for verify user
// @Tags registration
// @Accept json
// @Produce json
// @Param email query string true "email"
// @Param code query string true "code"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users/verify [post]
func (h *HandlerV1) Verify(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	code := c.Query("code")
	email := c.Query("email")

	userData, err := h.redisStorage.Get(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message : err.Error(),
		})
		log.Println(err)
		return
	}
	var user models.User

	err = json.Unmarshal(userData, &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}
	err = h.redisStorage.Del(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	if user.Email != email {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "The email did not match ",
		})
		log.Println(err)
		return
	}

	id := uuid.NewString()

	h.RefreshToken = refresh_token.JWTHandler{
		Sub:        id,
		Role:       "user",
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

	hashPassword, err := regtool.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusBadGateway, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	claims, err := refresh_token.ExtractClaim(access, []byte(h.Config.Token.SignInKey))
	if err != nil {
		c.JSON(http.StatusBadGateway, models.Error{
			Message: err.Error(),
		})
	}

	_, err = h.Service.UserService().CreateUser(ctx, &pb.User{
		Id:          id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Password:    hashPassword,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		Age:         user.Age,
		Refresh:     refresh,
		Role:        cast.ToString(claims["role"]),
	})

	respUser := &models.UserResponse{
		Id:          id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Password:    hashPassword,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		Age:         user.Age,
		Role:        cast.ToString(claims["role"]),
		Access:      access,
		Refresh:     refresh,
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message : err.Error(),
		})
		log.Println(err)
		return
	}
	c.JSON(http.StatusCreated, respUser)

}

// Login  users
// @Summary LoginUser
// @Description Api for user user
// @Tags registration
// @Accept json
// @Produce json
// @Param email query string true "email"
// @Param password query string true "password"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users/login [post]
func (h *HandlerV1) LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	email := c.Query("email")
	password := c.Query("password")

	filter := map[string]string{
		"email": email,
	}

	response, err := h.Service.UserService().GetUser(
		ctx, &pb.Filter{
			Filter: filter,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	if !(regtool.CheckHashPassword(password, response.Password)) {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: "Noto'gri parol!!!!",
		})
		return
	}

	h.RefreshToken = refresh_token.JWTHandler{
		Sub:        response.Id,
		Role:       response.Role,
		SigningKey: h.Config.Token.SignInKey,
		Log:        h.Logger,
		Email:      response.Email,
	}

	access, refresh, err := h.RefreshToken.GenerateAuthJWT()

	if err != nil {
		c.JSON(http.StatusConflict, models.Error{
		   Message : err.Error(),
		})
		log.Println(err)
		return
	}
	respUser := &models.UserResponse{
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
		c.JSON(http.StatusConflict, models.Error{
		   Message : err.Error(),
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusCreated, respUser)

}

// update acces token
// @Summary LoginUser
// @Description Api for updated acces token
// @Tags registration
// @Accept json
// @Produce json
// @Param refreshToken query string true "Refresh Token"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users/token [post]
func (h *HandlerV1) Token(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	RToken := c.Query("refreshToken")

	filter := map[string]string{
		"refresh_token": RToken,
	}
	user, err := h.Service.UserService().GetUser(ctx, &pb.Filter{
		Filter: filter,
	})

	if err != nil {
		c.JSON(500, models.Error{
		   Message : err.Error(),
		})
		log.Println(err)
		return
	}

	resclaim, err := refresh_token.ExtractClaim(user.Refresh, []byte(h.Config.Token.SignInKey))
	if err != nil {
		c.JSON(500, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}
	Now_time := time.Now().Unix()
	exp := (resclaim["exp"])
	if exp.(float64)-float64(Now_time) > 0 {
		h.RefreshToken = refresh_token.JWTHandler{
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

		respUser := &models.UserResponse{
			Id:          user.Id,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			Password:    user.Password,
			PhoneNumber: user.PhoneNumber,
			Gender:      user.Gender,
			Age:         user.Age,
			Role:        user.Role,
			Refresh:     refresh,
			Access:      access,
		}

		c.JSON(http.StatusCreated, respUser)
	} else {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "login",
		})
		log.Println(err)
		return
	}
}

// Forget  password
// @Summary ForgetPassword
// @Description Api for sending otp
// @Tags registration
// @Accept json
// @Produce json
// @Param email query string true "email"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users/forgetpassword [post]
func (h *HandlerV1) ForgetPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(7))
	defer cancel()

	email := c.Query("email")

	status, err := h.Service.UserService().UniqueEmail(ctx, &pb.IsUnique{
		Email: email,
	})

	if err != nil {
		c.JSON(500, models.Error{
		    Message: err.Error(),
		})
		log.Println(err)
		return

	}

	if !status.Status {
		c.JSON(500, models.Error{
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
	if err := h.redisStorage.Set(ctx, email, cast.ToString(radomNumber), time.Second*300); err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardErrorModel{
			Error: models.Error{
				Message: err.Error(),
			},
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Response: "We have sent otp your email",
	})
}

// verify forget password
// @Summary VerifyForgetPassword
// @Description Api for verify forget password
// @Tags registration
// @Accept json
// @Produce json
// @Param User body models.VerifyForgetPassword true "VerifyForgotPassword"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users/verify/forgetpassword [post]
func (h *HandlerV1) VerifyForgetPassword(c *gin.Context) {
	var (
		body models.VerifyForgetPassword
	)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(7))
	defer cancel()
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
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
	otp, err := h.redisStorage.Get(ctx, body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	err = h.redisStorage.Del(ctx, body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}
	if string(otp)[1:len(string(otp))-1] != body.Otp {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "The entered code did not match",
		})
		return
	}

	hashPassword, err := regtool.HashPassword(body.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadGateway, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	filter := map[string]string{
		"email": body.Email,
	}
	respUser, err := h.Service.UserService().GetUser(ctx, &pb.Filter{
		Filter: filter,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	h.Service.UserService().UpdatePassword(ctx, &pb.UpdatePasswordRequest{
		UserId:      respUser.Id,
		NewPassword: hashPassword,
	})
	c.JSON(http.StatusCreated, models.Response{
		Response: "Your password succesfully updated",
	})
}
