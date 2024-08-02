package v1

import (
	_ "api-gateway/api/docs"
	"api-gateway/api/models"
	userserviceproto "api-gateway/genproto/user_service"
	"api-gateway/internal/pkg/regtool"
	"api-gateway/internal/pkg/token"
	"api-gateway/internal/pkg/validation"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Summary 		Register
// @Description 	Api for register user
// @Tags 			auth
// @Accept 			json
// @Produce 		json
// @Param 			User body models.UserRegister true "Register User"
// @Success 		200 {object} string
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
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		log.Println(err.Error())
		return
	}

	body.Email, err = validation.EmailValidation(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.InvalidEmailMessage,
		})
		log.Println(err.Error())
		return
	}

	status := validation.PasswordValidation(body.Password)
	if !status {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WeakPasswordMessage,
		})
		log.Println(models.WeakPasswordMessage, body.Password)
		return
	}

	body.Gender = strings.ToLower(body.Gender)
	if body.Gender != "male" && body.Gender != "female" {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.InvalidGenderMessage,
		})
		log.Println(body.Gender, models.InvalidGenderMessage)
		return
	}

	exists, err := h.Service.UserService().UniqueEmail(ctx, &userserviceproto.IsUnique{
		Email: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}
	if exists.Status {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.EmailUsedMessage,
		})
		log.Println(models.EmailUsedMessage, body.Email)
		return
	}

	randomNumber, err := regtool.SendCodeGmail(body.Email, "Clothes Store\n", "./internal/pkg/regtool/emailotp.html", h.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(models.InternalMessage)
		return
	}

	err = h.redisStorage.Set(ctx, randomNumber, body, time.Minute*5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.SentOTPMessage)
}

// @Summary Verify User
// @Description Api for verify user
// @Tags auth
// @Accept json
// @Produce json
// @Param email query string true "email"
// @Param code query string true "code"
// @Success 201 {object} models.User
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/verify/register [post]
func (h *HandlerV1) VerifyRegister(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	code := c.Query("code")
	email := c.Query("email")

	userData, err := h.redisStorage.Get(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}
	var user models.UserRegister

	err = json.Unmarshal(userData, &user)
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

	h.RefreshToken = token.JWTHandler{
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

	claims, err := token.ExtractClaim(access, []byte(h.Config.Token.SignInKey))
	if err != nil {
		c.JSON(http.StatusBadGateway, models.Error{
			Message: err.Error(),
		})
	}

	_, err = h.Service.UserService().CreateUser(ctx, &userserviceproto.User{
		Id:        id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashPassword,
		Gender:    user.Gender,
		Refresh:   refresh,
		Role:      cast.ToString(claims["role"]),
	})

	respUser := &models.User{
		Id:        id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Gender:    user.Gender,
		Access:    access,
		Refresh:   refresh,
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}
	c.JSON(http.StatusCreated, respUser)

}

// @Summary 		Login
// @Description 	Api for user user
// @Tags 			auth
// @Accept 			json
// @Produce 		json
// @Param 			login body models.Login true "Login Model"
// @Success 		200 {object} models.LoginResp
// @Failure 		400 {object} models.Error
// @Failure 		404 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/login [POST]
func (h *HandlerV1) Login(c *gin.Context) {
	var body models.Login

	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		log.Println(err.Error())
		return
	}

	response, err := h.Service.UserService().GetUser(
		ctx, &userserviceproto.Filter{
			Filter: map[string]string{
				"email": body.Email,
			},
		})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: models.NotFoundMessage,
		})
		log.Println(err.Error())
		return
	}

	if !(regtool.CheckHashPassword(body.Password, response.Password)) {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.IncorrectPasswordMessage,
		})
		return
	}

	h.RefreshToken = token.JWTHandler{
		Sub:        response.Id,
		Role:       response.Role,
		SigningKey: h.Config.Token.SignInKey,
		Log:        h.Logger,
		Email:      response.Email,
	}

	access, refresh, err := h.RefreshToken.GenerateAuthJWT()

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}

	respUser := &models.LoginResp{
		Id:          response.Id,
		FirstName:   response.FirstName,
		LastName:    response.LastName,
		Email:       response.Email,
		PhoneNumber: response.PhoneNumber,
		Gender:      response.Gender,
		Age:         response.Age,
		Role:        response.Role,
		Refresh:     refresh,
		Access:      access,
	}
	_, err = h.Service.UserService().UpdateRefresh(ctx, &userserviceproto.RefreshRequest{
		UserId:       response.Id,
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, respUser)
}

// @Summary 		Forgot Password
// @Description 	Api for sending otp
// @Tags 			auth
// @Accept 			json
// @Produce 		json
// @Param 			email path string true "Email"
// @Success 		200 {object} string
// @Failure 		400 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/forgot/{email} [POST]
func (h *HandlerV1) Forgot(c *gin.Context) {
	email := c.Param("email")

	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	email, err = validation.EmailValidation(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.InvalidEmailMessage,
		})
		log.Println(err.Error())
		return
	}

	user, err := h.Service.UserService().GetUser(ctx, &userserviceproto.Filter{
		Filter: map[string]string{
			"email": email,
		},
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: models.NotFoundMessage,
		})
		log.Println(err.Error())
		return
	}

	radomNumber, err := regtool.SendCodeGmail(user.Email, "ClothesStore\n", "./internal/pkg/regtool/forgotpassword.html", h.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	if err := h.redisStorage.Set(ctx, radomNumber, cast.ToString(email), time.Minute*5); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.SentOTPMessage)
}

// @Summary 		Verify OTP
// @Description 	Api for verify user
// @Tags 			auth
// @Accept 			json
// @Produce 		json
// @Param 			email query string true "Email"
// @Param 			otp query string true "OTP"
// @Success 		200 {object} bool
// @Failure 		400 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/verify [POST]
func (h *HandlerV1) Verify(c *gin.Context) {
	otp := c.Query("otp")
	email := c.Query("email")

	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	userData, err := h.redisStorage.Get(ctx, otp)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}
	var redisEmail string

	err = json.Unmarshal(userData, &redisEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		log.Println(err.Error())
		return
	}

	if redisEmail != email {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotMatchOTP,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}

// @Summary 		Reset Password
// @Description 	Api for reset password
// @Tags 			auth
// @Accept 			json
// @Produce 		json
// @Param 			User body models.ResetPassword true "Reset Password"
// @Success 		200 {object} bool
// @Failure 		400 {object} models.Error
// @Failure 		404 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/reset-password [PUT]
func (h *HandlerV1) ResetPassword(c *gin.Context) {
	var (
		body models.ResetPassword
	)
	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		log.Println(err.Error())
		return
	}

	status := validation.PasswordValidation(body.NewPassword)
	if !status {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WeakPasswordMessage,
		})
		log.Println(models.WeakPasswordMessage, body.NewPassword)
		return
	}

	user, err := h.Service.UserService().GetUser(ctx, &userserviceproto.Filter{
		Filter: map[string]string{
			"email": body.Email,
		},
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: models.NotFoundMessage,
		})
		log.Println(err.Error())
		return
	}

	hashPassword, err := regtool.HashPassword(body.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}

	responseStatus, err := h.Service.UserService().UpdatePassword(ctx, &userserviceproto.UpdatePasswordRequest{
		UserId:      user.Id,
		NewPassword: hashPassword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}
	if !responseStatus.Status {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.NotUpdatedMessage,
		})
		log.Println(models.NotUpdatedMessage)
		return
	}

	c.JSON(http.StatusCreated, true)
}

// @Security		BearerAuth
// @Summary 		Update Password
// @Description		This API for update password with
// @Tags			auth
// @Param			info body models.UpdatePassword true "Update Password"
// @Success			200 {object} string
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		404 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router			/v1/update-password [PUT]
func (h *HandlerV1) UpdatePassword(c *gin.Context) {
	var info models.UpdatePassword

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

	userId, statusCode := regtool.GetIdFromToken(c.Request, &h.Config)

	if statusCode == http.StatusUnauthorized {
		c.JSON(http.StatusUnauthorized, models.Error{
			Message: models.NoAccessMessage,
		})
		return
	}

	err = c.ShouldBindJSON(&info)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		log.Println(err.Error())
		return
	}

	userInfo, err := h.Service.UserService().GetUser(ctx, &userserviceproto.Filter{
		Filter: map[string]string{
			"id": userId,
		},
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: models.NotFoundMessage,
		})
		log.Println(err.Error())
		return
	}

	if !(regtool.CheckHashPassword(info.PresetPassword, userInfo.Password)) {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.IncorrectPasswordMessage,
		})
		return
	}
	if info.ConfirmPassword != info.NewPassword {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.IncorrectPasswordMessage,
		})
		return
	}

	hashConfirmPassword, err := regtool.HashPassword(info.ConfirmPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}

	respStatus, err := h.Service.UserService().UpdatePassword(ctx, &userserviceproto.UpdatePasswordRequest{
		UserId:      userId,
		NewPassword: hashConfirmPassword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}
	if !respStatus.Status {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessUpdatedPassword)
}

// @Summary 		New Token
// @Description 	Api for updated acces token
// @Tags 			auth
// @Accept 			json
// @Produce 		json
// @Param 			refresh path string true "Refresh Token"
// @Success 		200 {object} models.TokenResp
// @Failure 		400 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/token/{refresh} [GET]
func (h *HandlerV1) NewToken(c *gin.Context) {
	RToken := c.Param("refresh")

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

	if RToken == "" {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "refresh token is required",
		})
		log.Println("refresh token is required")
		return
	}

	claims, err := token.ExtractClaim(RToken, []byte(h.Config.Token.SignInKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "refresh token is expired",
		})
		log.Println(err.Error())
		return
	}

	userRefreshTokenSub := cast.ToString(claims["sub"])

	user, err := h.Service.UserService().GetUser(ctx, &userserviceproto.Filter{
		Filter: map[string]string{
			"id": userRefreshTokenSub,
		},
	})

	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: models.NotFoundMessage,
		})
		log.Println(err.Error())
		return
	}

	nowTime := time.Now().Unix()
	exp := claims["exp"]
	if exp.(float64)-float64(nowTime) > 0 {
		h.RefreshToken = token.JWTHandler{
			Sub:        user.Id,
			Role:       user.Role,
			SigningKey: h.Config.Token.SignInKey,
			Log:        h.Logger,
			Email:      user.Email,
		}

		access, refresh, err := h.RefreshToken.GenerateAuthJWT()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			log.Println(err)
			return
		}

		_, err = h.Service.UserService().UpdateRefresh(ctx, &userserviceproto.RefreshRequest{
			UserId:       user.Id,
			RefreshToken: refresh,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			log.Println(err.Error())
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
			Message: models.TokenInvalidMessage,
		})
		return
	}
}
