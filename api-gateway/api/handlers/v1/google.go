package v1

import (
	userserviceproto "api-gateway/genproto/user_service"
	"api-gateway/internal/pkg/token"
	"context"
	"encoding/json"
	"api-gateway/api/models"
	"api-gateway/internal/pkg/config"
	regtool "api-gateway/internal/pkg/regtool"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary 		Google Login
// @Description 	Redirects the user to Google's OAuth 2.0 consent page
// @Tags 			oauth
// @Success 		303 {string} string "Redirect"
// @Router 			/v1/google/login [GET]
func (h *HandlerV1) GoogleLogin(c *gin.Context) {
	googleConfig := config.SetupConfig()
	url := googleConfig.AuthCodeURL("RandomState")
	c.Redirect(http.StatusSeeOther, url)

	c.JSON(303, url)
}

// @Summary 		Handle Google callback
// @Description 	Handles the callback from Google OAuth 2.0,
// @Description 	exchanges code for token and retrieves user info
// @Tags 			oauth
// @Param 			state query string true "OAuth State"
// @Param 			code query string true "OAuth Code"
// @Success 		200 {string} models.LoginResp "User info"
// @Failure 		400 {string} string "Bad Request"
// @Failure 		401 {string} string "Unauthorized"
// @Failure 		500 {string} string "Internal Server Error"
// @Router 			/v1/google/callback [GET]
func (h *HandlerV1) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

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

	if state != "RandomState" {
		c.JSON(http.StatusUnauthorized, models.MisMatchMessage)
		return
	}
	if code == "" {
		c.JSON(http.StatusBadRequest, models.MisCode)
		return
	}

	googleConfig := config.SetupConfig()

	tokens, err := googleConfig.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.TokenExchangeMessage)
		log.Println(err.Error())
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tokens.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailedGetUserInfo)
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailedGetUserInfo)
		log.Println(err.Error())
		return
	}
	var body models.GoogleUser

	err = json.Unmarshal(userData, &body)
	if err != nil {
		c.JSON(http.StatusSeeOther, models.Error{
			Message: models.InternalMessage,
		})
	}

	id := uuid.New().String()
	hashpassword, err := regtool.HashPassword("salom")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}
	h.RefreshToken = token.JWTHandler{
		Sub:        id,
		Role:       "user",
		SigningKey: h.Config.Token.SignInKey,
		Log:        h.Logger,
		Email:      body.Email,
	}

	access, refresh, err := h.RefreshToken.GenerateAuthJWT()

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}

	status, err := h.Service.UserService().UniqueEmail(ctx, &userserviceproto.IsUnique{
		Email: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.EmailUsedMessage,
		})
		log.Println(err.Error())
		return
	}
	if status.Status {
		filter := map[string]string{
			"email": body.Email,
		}
		responesUser, err := h.Service.UserService().GetUser(ctx, &userserviceproto.Filter{
			Filter: filter,
		})
		if err != nil {
			c.JSON(http.StatusNotFound, models.Error{
				Message: models.NotFoundMessage,
			})
			log.Println(err.Error())
			return
		}
		c.JSON(http.StatusOK, models.LoginResp{
			Id:          responesUser.Id,
			FirstName:   responesUser.FirstName,
			LastName:    responesUser.LastName,
			Email:       responesUser.Email,
			PhoneNumber: responesUser.PhoneNumber,
			Gender:      responesUser.Gender,
			Role:        responesUser.Role,
			Refresh:     refresh,
			Access:      access,
		})
		return
	}

	Resp, err := h.Service.UserService().CreateUser(ctx, &userserviceproto.User{
		Id:        id,
		FirstName: body.GivenName,
		LastName:  body.FamilyName,
		Email:     body.Email,
		Password:  hashpassword,
		Gender:    "male",
		Role:      "user",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.LoginResp{
		Id:        Resp.Guid,
		FirstName: body.GivenName,
		LastName:  body.FamilyName,
		Email:     body.Email,
		Gender:    "male",
		Role:      "user",
		Refresh:   refresh,
		Access:    access,
	})
}
