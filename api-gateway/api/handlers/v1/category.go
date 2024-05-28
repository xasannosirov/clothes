package v1

import (
	"api-gateway/api/models"
	"api-gateway/genproto/media_service"
	"api-gateway/genproto/product_service"
	"api-gateway/internal/pkg/query_parameter"
	"api-gateway/internal/pkg/regtool"
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

	listCategories, err := h.Service.ProductService().ListCategories(ctx, &product_service.ListRequest{
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

// @Security 		BearerAuth
// @Summary 		Search Category
// @Description 	This api search products with category
// @Tags			category
// @Accept 			application/json
// @Produce 		applocation/json
// @Param			page query uint64 true "Page"
// @Param			limit query uint64 true "Limit"
// @Param			name query string true "Category Name"
// @Success 		200 {object} models.ListProduct
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		404 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router			/v1/category/search [GET]
func (h *HandlerV1) SearchCategory(c *gin.Context) {
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

	categoryName := c.Query("name")
	parameters := query_parameter.New(c.Request.URL.Query())

	products, err := h.Service.ProductService().SearchCategory(ctx, &product_service.SearchRequest{
		Page:  parameters.GetPage(),
		Limit: parameters.GetLimit(),
		Params: map[string]string{
			"name": categoryName,
		},
	})

	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: models.NotFoundMessage,
		})
		log.Println(err.Error())
		return
	}

	token := c.Request.Header.Get("Authorization")
	var response models.ListProduct
	if token != "" {
		for _, product := range products.Products {
			media, err := h.Service.MediaService().Get(ctx, &media_service.MediaWithID{
				Id: product.Id,
			})
			if err != nil {
				log.Println(err.Error())
			}

			var imagesURL []string
			for _, imageUrl := range media.Images {
				imagesURL = append(imagesURL, imageUrl.ImageUrl)
			}

			userId, statusCode := regtool.GetIdFromToken(c.Request, &h.Config)
			if statusCode != 0 {
				c.JSON(http.StatusBadRequest, models.Error{
					Message: "oops something went wrong",
				})
			}

			likeStatus, err := h.Service.ProductService().IsUnique(ctx, &product_service.IsUniqueReq{
				TableName: "wishlist",
				UserId:    userId,
				ProductId: product.Id,
			})
			if err != nil {
				c.JSON(http.StatusBadRequest, models.Error{
					Message: err.Error(),
				})
				log.Println(err.Error())
				return
			}
			basketStatus, err := h.Service.ProductService().IsUnique(ctx, &product_service.IsUniqueReq{
				TableName: "basket",
				UserId:    userId,
				ProductId: product.Id,
			})
			if err != nil {
				c.JSON(http.StatusBadRequest, models.Error{
					Message: err.Error(),
				})
				log.Println(err.Error())
				return
			}

			response.Products = append(response.Products, models.Product{
				ID:          product.Id,
				Name:        product.Name,
				Category:    product.Category,
				Description: product.Description,
				MadeIn:      product.MadeIn,
				Color:       product.Color,
				Size:        product.ProductSize,
				Count:       product.Count,
				Cost:        float64(product.Cost),
				Discount:    float64(product.Discount),
				AgeMin:      product.AgeMin,
				AgeMax:      product.AgeMax,
				ForGender:   product.ForGender,
				Liked:       likeStatus.Status,
				Basket:      basketStatus.Status,
				ImageURL:    imagesURL,
			})
		}
	} else {
		for _, product := range products.Products {
			media, err := h.Service.MediaService().Get(ctx, &media_service.MediaWithID{
				Id: product.Id,
			})
			if err != nil {
				log.Println(err.Error())
			}

			var imagesURL []string
			for _, imageUrl := range media.Images {
				imagesURL = append(imagesURL, imageUrl.ImageUrl)
			}

			response.Products = append(response.Products, models.Product{
				ID:          product.Id,
				Name:        product.Name,
				Category:    product.Category,
				Description: product.Description,
				MadeIn:      product.MadeIn,
				Color:       product.Color,
				Size:        product.ProductSize,
				Count:       product.Count,
				Cost:        float64(product.Cost),
				Discount:    float64(product.Discount),
				AgeMin:      product.AgeMin,
				AgeMax:      product.AgeMax,
				ForGender:   product.ForGender,
				Liked:       false,
				Basket:      false,
				ImageURL:    imagesURL,
			})
		}
	}
	response.Total = products.TotalCount

	c.JSON(http.StatusOK, response)
}
