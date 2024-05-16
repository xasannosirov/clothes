package v1

import (
	"api-gateway/api/models"
	"api-gateway/genproto/product_service"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Security 		BearerAuth
// @Summary 		Create Category
// @Description 	This API for create a new category for product
// @Tags 			category
// @Produce 		json
// @Accept 			json
// @Param 			order body models.CategoryReq true "Create Category Model"
// @Success			201 {object} models.Category
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/category [POST]
func (h *HandlerV1) CreateCategory(c *gin.Context) {
	var (
		body        models.CategoryReq
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

	category, err := h.Service.ProductService().CreateCategory(ctx, &product_service.Category{
		Id:   uuid.NewString(),
		Name: body.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, models.Category{
		ID:   category.Id,
		Name: category.Name,
	})
}

// @Security 		BearerAuth
// @Summary 		Update Category
// @Description 	This API for update a category
// @Tags 			category
// @Produce 		json
// @Accept 			json
// @Param 			order body models.Category true "Create Category Model"
// @Success			200 {object} models.Category
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/category [PUT]
func (h *HandlerV1) UpdateCategory(c *gin.Context) {
	var (
		body        models.Category
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

	category, err := h.Service.ProductService().UpdateCategory(ctx, &product_service.Category{
		Id:   body.ID,
		Name: body.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.Category{
		ID:   category.Id,
		Name: category.Name,
	})
}

// @Security 		BearerAuth
// @Summary 		Delete Category
// @Description 	This API for delete a category with id
// @Tags 			category
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Category ID"
// @Success			200 {object} bool
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		404 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/category/{id} [DELETE]
func (h *HandlerV1) DeleteCategory(c *gin.Context) {
	var (
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

	categoryID := c.Param("id")

	status, err := h.Service.ProductService().DeleteCategory(ctx, &product_service.GetWithID{
		Id: categoryID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, status.Status)
}

// @Security 		BearerAuth
// @Summary 		Get Category
// @Description 	This API for getting a category with id
// @Tags 			category
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Category ID"
// @Success			200 {object} models.Category
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/category/{id} [GET]
func (h *HandlerV1) GetCategory(c *gin.Context) {
	var (
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

	categoryID := c.Param("id")

	category, err := h.Service.ProductService().GetCategory(ctx, &product_service.GetWithID{
		Id: categoryID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.Category{
		ID:   categoryID,
		Name: category.Name,
	})
}

// @Security 		BearerAuth
// @Summary 		List Category
// @Description 	This API for getting categories
// @Tags 			category
// @Produce 		json
// @Accept 			json
// @Param 			page query uint64 true "Page"
// @Param 			limit query uint64 true "Limit"
// @Success			200 {object} string
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/categories [GET]
func (h *HandlerV1) ListCategory(c *gin.Context) {
	var (
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

	listCategories, err := h.Service.ProductService().GetAllCategory(ctx, &product_service.ListRequest{
		Page:  int64(pageInt),
		Limit: int64(limitInt),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	var categories []models.Category
	for _, category := range listCategories.Categories {
		categories = append(categories, models.Category{
			ID:   category.Id,
			Name: category.Name,
		})
	}

	c.JSON(http.StatusOK, models.ListCategory{
		Categories: categories,
		Total:      listCategories.TotalCount,
	})
}
