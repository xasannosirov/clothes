package v1

import (
	_ "api-gateway/api/docs"
	"context"
	"net/http"
	"time"

	"api-gateway/api/models"
	pbu "api-gateway/genproto/user_service"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// GetCategories
// @Security ApiKeyAuth
// @Router /v1/users/create [post]
// @Summary Get categories
// @Description Get categories
// @Tags user
// @Accept json
// @Produce json
// @Param User body models.CreateUser true "createUserModel"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
func (h HandlerV1) CreateUser(c *gin.Context) {

	var (
		body        models.CreateUser
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration((7)))
	defer cancel()
	response, err := h.Service.UserService().CreateUser(ctx, &pbu.User{
		FirstName: body.Name,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
		Role:      "user",
		Refresh:   "dasdfasdfdfdfasdfasdfda;sldkfja;sdflkajjsdf",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	respProduct := pbu.User{
		Id: response.Id,
		// LastName:  response.LastName,
		// FirstName: response.Name,
		// Email:     response.Email,
		// Password:  response.Password,
		Role: "user",
		// Refresh:   response.RefreshToken,
		// CreatedAt: response.CreatedAt,
		// UpdatedAt: response.UpdatedAt,
	}

	c.JSON(http.StatusCreated, respProduct)
}
