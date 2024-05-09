package v1

import (
	"api-gateway/api/models"
	userproto "api-gateway/genproto/user_service"
	"api-gateway/internal/pkg/validation"
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Security  		BearerAuth
// @Summary   		Create User
// @Description 	Api for create a new user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			user body models.UserRegister true "Create User Model"
// @Success 		201 {object} string
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/user/ [POST]
func (h *HandlerV1) CreateUser(c *gin.Context) {
	var (
		body        models.UserRegister
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
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

	body.Role = strings.ToLower(body.Role)
	if !(body.Role == "user") || !(body.Role == "worker") {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Invalid role",
		})
		log.Println(body.Role)
		return
	}

	body.Email, err = validation.EmailValidation(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
		return
	}

	userServiceCreateResponse, err := h.Service.UserService().CreateUser(ctx, &userproto.User{
		Id:        uuid.New().String(),
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
		Gender:    body.Gender,
		Role:      body.Role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, userServiceCreateResponse.Guid)
}

// @Security  		BearerAuth
// @Summary   		Update User
// @Description 	Api for update a user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			user body models.User true "Create User Model"
// @Success 		200 {object} models.User
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/user/ [PUT]
func (h *HandlerV1) UpdateUser(c *gin.Context) {
	var (
		body        models.User
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
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

	body.Email, err = validation.EmailValidation(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	status := validation.PhoneUz(body.PhoneNumber)
	if !status {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "phone number is invalid",
		})
		log.Println("phone number is invalid")
		return
	}

	updatedUser, err := h.Service.UserService().UpdateUser(ctx, &userproto.User{
		Id:          body.Id,
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Email:       body.Email,
		Password:    body.Password,
		PhoneNumber: body.PhoneNumber,
		Gender:      body.Gender,
		Age:         body.Age,
		Role:        body.Role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.User{
		Id:          updatedUser.Id,
		FirstName:   updatedUser.FirstName,
		LastName:    updatedUser.LastName,
		Email:       updatedUser.Email,
		PhoneNumber: updatedUser.PhoneNumber,
		Gender:      updatedUser.Gender,
		Age:         updatedUser.Age,
		Role:        updatedUser.Role,
	})
}

// @Security  		BearerAuth
// @Summary   		Delete User
// @Description 	Api for delete a user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "User ID"
// @Success 		200 {object} bool
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/user/{id} [DELETE]
func (h *HandlerV1) DeleteUser(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	userID := c.Param("id")
	response, err := h.Service.UserService().DeleteUser(ctx, &userproto.UserWithGUID{
		Guid: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	} else if !response.Status {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: "Server error",
		})
		log.Println(response.Status)
		return
	}

	c.JSON(http.StatusOK, true)
}

// @Security  		BearerAuth
// @Summary   		Get User
// @Description 	Api for getting a user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "User ID"
// @Success 		200 {object} models.User
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/user/{id} [GET]
func (h *HandlerV1) GetUser(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	userID := c.Param("id")
	filter := map[string]string{
		"id": userID,
	}
	response, err := h.Service.UserService().GetUser(ctx, &userproto.Filter{
		Filter: filter,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.User{
		Id:        userID,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		Email:     response.Email,
		Password:  response.Password,
		Gender:    response.Gender,
		Age:       response.Age,
		Role:      response.Role,
	})
}

// @Security  		BearerAuth
// @Summary   		Update User
// @Description 	Api for getting list user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			page query uint64 true "Page"
// @Param 			limit query uint64 true "Limit"
// @Param 			role query string true "Role"
// @Success 		200 {object} []models.User
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/users [GET]
func (h *HandlerV1) ListUsers(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	page := c.Query("page")
	limit := c.Query("limit")
	role := c.Query("role")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	listUsers, err := h.Service.UserService().GetAllUsers(ctx, &userproto.ListUserRequest{
		Page:  int64(pageInt),
		Limit: int64(limitInt),
		Role:  role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	var users []models.User
	for _, user := range listUsers.Users {
		users = append(users, models.User{
			Id:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			Gender:    user.Gender,
			Age:       user.Age,
			Role:      user.Role,
		})
	}

	c.JSON(http.StatusOK, models.ListUser{
		User:  users,
		Total: listUsers.TotalCount,
	})
}
