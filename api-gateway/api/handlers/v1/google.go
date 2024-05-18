package v1

import (
	"api-gateway/api/models"
	userproto "api-gateway/genproto/user_service"
	"api-gateway/internal/pkg/config"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	regtool "api-gateway/internal/pkg/regtool"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	tokens "api-gateway/internal/pkg/token"

)

// GoogleLogin godoc
// @Summary Redirect to Google for login
// @Description Redirects the user to Google's OAuth 2.0 consent page
// @Tags auth
// @Success 303 {string} string "Redirect"
// @Router /v1/google/login [get]
func (h *HandlerV1)GoogleLogin(c *gin.Context) {
	googleConfig := config.SetupConfig()
	url := googleConfig.AuthCodeURL("RandomState")
	c.Redirect(http.StatusSeeOther, url)
	c.JSON(303, url)
}

// GoogleCallback godoc
// @Summary Handle Google callback
// @Description Handles the callback from Google OAuth 2.0, exchanges code for token and retrieves user info
// @Tags auth
// @Param state query string true "OAuth State"
// @Param code query string true "OAuth Code"
// @Success 200 {string} models.LoginResp "User info"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/google/callback [get]
func (h *HandlerV1)GoogleCallback(c *gin.Context) {
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
	query := c.Request.URL.Query()
	state := query.Get("state")
	if state != "RandomState" {
		c.String(http.StatusUnauthorized, "state mismatch")
		return
	}

	code := query.Get("code")
	if code == "" {
		c.String(http.StatusBadRequest, "missing code")
		return
	}

	googleConfig := config.SetupConfig()

	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("token exchange error:", err)
		c.String(http.StatusInternalServerError, "token exchange failed: %v", err)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to get user info: %v", err)
		return
	}
	defer resp.Body.Close()

	userData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to read user info: %v", err)
		return
	}
	var body models.GoogleUser

	err = json.Unmarshal(userData, &body)
	if err != nil{
		c.JSON(303, models.Error{
			Message: err.Error(),
		})
	}
	id := uuid.New().String()
	hashpassword, err := regtool.HashPassword(body.VerifiedEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	h.RefreshToken = tokens.JWTHandler{
		Sub:        id,
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

	status, err := h.Service.UserService().UniqueEmail(ctx, &userproto.IsUnique{
		Email: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		return
	}
	if status.Status {
		filter := map[string]string{
			"email": body.Email,
		}
		responesUser, err := h.Service.UserService().GetUser(ctx, &userproto.Filter{
			Filter: filter,
		})
		if err != nil{
			c.JSON(400, models.Error{
				Message: err.Error(),
			})
		}
		c.JSON(http.StatusOK, models.LoginResp{
			Id: responesUser.Id,
			FirstName: responesUser.FirstName,
			LastName:  responesUser.LastName,
			Email: responesUser.Email,
			Password: responesUser.Password,
			PhoneNumber: responesUser.PhoneNumber,
			Gender: responesUser.Gender,
			Role: responesUser.Role,
			Refresh: refresh,
			Access: access,
		})
		return
	}
	


	Resp, err := h.Service.UserService().CreateUser(ctx, &userproto.User{
		Id:        id,
		FirstName: body.GivenName,
		LastName:  body.FamilyName,
		Email:     body.Email,
		Password:  hashpassword,
		Gender:    "erkak",
		Role:      "user",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}



	c.JSON(http.StatusOK, models.LoginResp{
		Id: Resp.Guid,
		FirstName: body.GivenName,
		LastName: body.FamilyName,
		Email: body.Email,
		Password: hashpassword,
		Gender: "erkak",
		Role: "user",
		Refresh: refresh,
		Access: access,
		
	})
}