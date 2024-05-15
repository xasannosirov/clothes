package v1

import (
	"api-gateway/api/models"
	userproto "api-gateway/genproto/user_service"
	regtool "api-gateway/internal/pkg/regtool"
	"api-gateway/internal/pkg/validation"
	"context"
	"log"
	"net/http"
	"strconv"
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
// @Success 		201 {object} models.UserCreateResponse
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/user [POST]
func (h *HandlerV1) CreateUser(c *gin.Context) {
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

	body.Email, err = validation.EmailValidation(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
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
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Email already used",
		})
		return
	}

	hashpassword, err := regtool.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	userServiceCreateResponse, err := h.Service.UserService().CreateUser(ctx, &userproto.User{
		Id:        uuid.New().String(),
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  hashpassword,
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

	c.JSON(http.StatusCreated, models.UserCreateResponse{
		UserID: userServiceCreateResponse.Guid,
	})
}

// @Security  		BearerAuth
// @Summary   		Update User
// @Description 	Api for update a user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			user body models.WorkerPut true "Update User Model"
// @Success 		200 {object} models.User
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/user [PUT]
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
	filter := map[string]string{
		"id": body.Id,
	}
	user, err := h.Service.UserService().GetUser(ctx, &userproto.Filter{Filter: filter})
	if err != nil{
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
	}

	if user.Email != body.Email{
		status, err := h.Service.UserService().UniqueEmail(ctx, &userproto.IsUnique{
			Email: body.Email,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Message: err.Error(),
			})
			log.Println(err.Error())
			return
		}
		if status.Status {
			c.JSON(http.StatusBadRequest, models.Error{
				Message: "email already used",
			})
			return
		}
	}

	
	if body.PhoneNumber != "" {
		status := validation.PhoneUz(body.PhoneNumber)
		if !status {
			c.JSON(http.StatusBadRequest, models.Error{
				Message: "phone number is invalid",
			})
			log.Println("phone number is invalid")
			return
		}
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
		Role:        "worker",
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
// @Failure 		404 {object} models.Error
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

	user, err := h.Service.UserService().GetUser(ctx, &userproto.Filter{
		Filter: map[string]string{
			"id": userID,
		},
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	if user.Role == "admin" {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Wrong request",
		})
		return
	}

	response, err := h.Service.UserService().DeleteUser(ctx, &userproto.UserWithGUID{
		Guid: userID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	} else if !response.Status {
		c.JSON(http.StatusNotFound, models.Error{
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
// @Failure 		404 {object} models.Error
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
		c.JSON(http.StatusNotFound, models.Error{
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
	})
}

// @Security  		BearerAuth
// @Summary   		Get  Delete User
// @Description 	Api for getting a deleted user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "User ID"
// @Success 		200 {object} models.User
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/del/user/{id} [GET]
func (h *HandlerV1) GetDelUser(c *gin.Context) {
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
	response, err := h.Service.UserService().GetUserDelete(ctx, &userproto.Filter{
		Filter: filter,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
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
	})
}

// @Security  		BearerAuth
// @Summary   		List User
// @Description 	Api for getting list user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			page query uint64 true "Page"
// @Param 			limit query uint64 true "Limit"
// @Success 		200 {object} models.ListUser
// @Failure 		404 {object} models.Error
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
		Role:  "user",
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
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
		})
	}

	c.JSON(http.StatusOK, models.ListUser{
		User:  users,
		Total: listUsers.TotalCount,
	})
}

// @Security  		BearerAuth
// @Summary   		Create Worker
// @Description 	Api for create a new worker
// @Tags 			workers
// @Accept 			json
// @Produce 		json
// @Param 			worker body models.WorkerPost true "Create Worker Model"
// @Success 		201 {object} models.UserCreateResponse
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/worker [POST]
func (h *HandlerV1) CreateWorker(c *gin.Context) {
	var (
		body        models.WorkerPost
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

	body.Email, err = validation.EmailValidation(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.InternalMessage,
		})
		log.Println(err.Error())
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
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Eamil already used",
		})
		return
	}

	st := validation.PhoneUz(body.PhoneNumber)
	if !st {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "phone number is invalid",
		})
		log.Println("phone number is invalid")
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

	userServiceCreateResponse, err := h.Service.UserService().CreateUser(ctx, &userproto.User{
		Id:          uuid.New().String(),
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Email:       body.Email,
		Password:    hashPassword,
		PhoneNumber: body.PhoneNumber,
		Gender:      body.Gender,
		Role:        "worker",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, models.UserCreateResponse{
		UserID: userServiceCreateResponse.Guid,
	})
}

// @Security  		BearerAuth
// @Summary   		Update Worker
// @Description 	Api for update a user
// @Tags 			workers
// @Accept 			json
// @Produce 		json
// @Param 			worker body models.WorkerPut true "Update Worker Model"
// @Success 		200 {object} models.WorkerPut
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/worker [PUT]
func (h *HandlerV1) UpdateWorker(c *gin.Context) {
	var (
		body        models.WorkerPut
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
	filter := map[string]string{
		"id": body.ID,
	}
	user, err := h.Service.UserService().GetUser(ctx, &userproto.Filter{Filter: filter})
	if err != nil{
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
	}

	if user.Email != body.Email{
		status, err := h.Service.UserService().UniqueEmail(ctx, &userproto.IsUnique{
			Email: body.Email,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Message: err.Error(),
			})
			log.Println(err.Error())
			return
		}
		if status.Status {
			c.JSON(http.StatusBadRequest, models.Error{
				Message: "email already used",
			})
			return
		}
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
		Id:          body.ID,
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Email:       body.Email,
		Password:    body.Password,
		PhoneNumber: body.PhoneNumber,
		Gender:      body.Gender,
		Age:         body.Age,
		Role:        "worker",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.WorkerPut{
		ID:          updatedUser.Id,
		FirstName:   updatedUser.FirstName,
		LastName:    updatedUser.LastName,
		Email:       updatedUser.Email,
		PhoneNumber: updatedUser.PhoneNumber,
		Gender:      updatedUser.Gender,
		Age:         updatedUser.Age,
	})
}

// @Security  		BearerAuth
// @Summary   		Delete Worker
// @Description 	Api for delete a worker
// @Tags 			workers
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Worker ID"
// @Success 		200 {object} bool
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/worker/{id} [DELETE]
func (h *HandlerV1) DeleteWorker(c *gin.Context) {
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

	user, err := h.Service.UserService().GetUser(ctx, &userproto.Filter{
		Filter: map[string]string{
			"id": userID,
		},
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		return
	}
	if user.Role == "admin" {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Bad request",
		})
		return
	}

	response, err := h.Service.UserService().DeleteUser(ctx, &userproto.UserWithGUID{
		Guid: userID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	} else if !response.Status {
		c.JSON(http.StatusNotFound, models.Error{
			Message: "Server error",
		})
		log.Println(response.Status)
		return
	}

	c.JSON(http.StatusOK, true)
}

// @Security  		BearerAuth
// @Summary   		Get Worker
// @Description 	Api for getting a worker
// @Tags 			workers
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Worker ID"
// @Success 		200 {object} models.User
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/worker/{id} [GET]
func (h *HandlerV1) GetWorker(c *gin.Context) {
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
		c.JSON(http.StatusNotFound, models.Error{
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
	})
}

// @Security  		BearerAuth
// @Summary   		List Worker
// @Description 	Api for getting list worker
// @Tags 			workers
// @Accept 			json
// @Produce 		json
// @Param 			page query uint64 true "Page"
// @Param 			limit query uint64 true "Limit"
// @Success 		200 {object} models.ListUser
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/workers [GET]
func (h *HandlerV1) ListWorker(c *gin.Context) {
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
		Role:  "worker",
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
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
		})
	}

	c.JSON(http.StatusOK, models.ListUser{
		User:  users,
		Total: listUsers.TotalCount,
	})
}

// @Security  		BearerAuth
// @Summary   		Test
// @Description 	API get user_id
// @Tags 			test
// @Accept 			json
// @Produce 		json
// @Success 		200 {object} string
// @Failure 		400 {object} models.Error
// @Router 			/v1/test [GET]
func (h *HandlerV1) Test(c *gin.Context) {
	userID, statusCode := GetIdFromToken(c.Request, &h.Config)

	if statusCode == 401 {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Access token expired",
		})
		return
	}

	c.JSON(http.StatusOK, userID)
}
