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
	"google.golang.org/protobuf/encoding/protojson"
)

// @Security 		BearerAuth
// @Summary 		Create Product
// @Description 	This API for create a new product
// @Tags 			products
// @Produce 		json
// @Accept 			json
// @Param 			product body models.ProductReq true "Create Product Model"
// @Success			201 {object} models.ProductCreateResponse
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/product [POST]
func (h *HandlerV1) CreateProduct(c *gin.Context) {
	var (
		body        models.ProductReq
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

	createdProductResponse, err := h.Service.ProductService().CreateProduct(ctx, &product_service.Product{
		Name:        body.Name,
		Description: body.Description,
		Category:    body.Category,
		MadeIn:      body.MadeIn,
		Color:       body.Color,
		Count:       body.Count,
		Cost:        float32(body.Cost),
		Discount:    float32(body.Discount),
		AgeMin:      body.AgeMin,
		AgeMax:      body.AgeMax,
		ForGender:   body.ForGender,
		Size_:       body.Size,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, models.ProductCreateResponse{
		ProductID: createdProductResponse.Id,
	})
}

// @Security 		BearerAuth
// @Summary 		Update Product
// @Description 	This API for updating a product
// @Tags 			products
// @Produce 		json
// @Accept 			json
// @Param 			product body models.Product true "Update Product Model"
// @Success			200 {object} models.Product
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/product [PUT]
func (h *HandlerV1) UpdateProduct(c *gin.Context) {
	var (
		body        models.Product
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

	updatedProduct, err := h.Service.ProductService().UpdateProduct(ctx, &product_service.Product{
		Id:          body.ID,
		Name:        body.Name,
		Description: body.Description,
		Category:    body.Category,
		MadeIn:      body.MadeIn,
		Color:       body.Color,
		Count:       body.Count,
		Cost:        float32(body.Cost),
		Discount:    float32(body.Discount),
		AgeMin:      body.AgeMin,
		AgeMax:      body.AgeMax,
		ForGender:   body.ForGender,
		Size_:       body.Size,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.Product{
		ID:          updatedProduct.Id,
		Name:        updatedProduct.Name,
		Category:    updatedProduct.Category,
		Description: updatedProduct.Description,
		MadeIn:      updatedProduct.MadeIn,
		Color:       updatedProduct.Color,
		Size:        updatedProduct.Size_,
		Count:       updatedProduct.Count,
		Cost:        float64(updatedProduct.Cost),
		Discount:    float64(updatedProduct.Discount),
		AgeMin:      updatedProduct.AgeMin,
		AgeMax:      updatedProduct.AgeMax,
		ForGender:   updatedProduct.ForGender,
	})
}

// @Security 		BearerAuth
// @Summary 		Delete Product
// @Description 	This API for deleting a product with product_id
// @Tags 			products
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Product ID"
// @Success			200 {object} bool
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/product/{id} [DELETE]
func (h *HandlerV1) DeleteProduct(c *gin.Context) {
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

	productID := c.Param("id")

	responseStatus, err := h.Service.ProductService().DeleteProduct(ctx, &product_service.GetWithID{
		Id: productID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, responseStatus.Status)
}

// @Security 		BearerAuth
// @Summary 		Get Product
// @Description 	This API for getting a product with product_id
// @Tags 			products
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Product ID"
// @Success			200 {object} models.Product
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/product/{id} [GET]
func (h *HandlerV1) GetProduct(c *gin.Context) {
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

	productID := c.Param("id")

	product, err := h.Service.ProductService().GetProductByID(ctx, &product_service.GetWithID{
		Id: productID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.Product{
		ID:          productID,
		Name:        product.Name,
		Category:    product.Category,
		Description: product.Description,
		MadeIn:      product.MadeIn,
		Color:       product.Color,
		Size:        product.Size_,
		Count:       product.Count,
		Cost:        float64(product.Cost),
		Discount:    float64(product.Discount),
		AgeMin:      product.AgeMin,
		AgeMax:      product.AgeMax,
		ForGender:   product.ForGender,
	})
}

// @Security 		BearerAuth
// @Summary 		List Products
// @Description 	This API for getting list of products
// @Tags 			products
// @Produce 		json
// @Accept 			json
// @Param 			page query uint64 true "Page"
// @Param 			limit query uint64 true "Limit"
// @Param 			name query string false "Product Name"
// @Success			200 {object} models.ListProduct
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/products [GET]
func (h *HandlerV1) ListProducts(c *gin.Context) {
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
	name := c.Query("name")

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

	listProducts, err := h.Service.ProductService().GetAllProducts(ctx, &product_service.ListProductRequest{
		Page:  uint64(pageInt),
		Limit: uint64(limitInt),
		Name:  name,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	var products []models.Product
	for _, product := range listProducts.Products {
		products = append(products, models.Product{
			ID:          product.Id,
			Name:        product.Name,
			Category:    product.Category,
			Description: product.Description,
			MadeIn:      product.MadeIn,
			Color:       product.Color,
			Size:        product.Size_,
			Count:       product.Count,
			Cost:        float64(product.Cost),
			Discount:    float64(product.Discount),
			AgeMin:      product.AgeMin,
			AgeMax:      product.AgeMax,
			ForGender:   product.ForGender,
		})
	}

	response := models.ListProduct{}
	response.Products = products
	if len(products) == 0 {
		response.Total = 0
	} else if len(products) < int(listProducts.TotalCount) {
		response.Total = uint64(len(products))
	} else {
		response.Total = listProducts.TotalCount
	}

	c.JSON(http.StatusOK, response)
}
