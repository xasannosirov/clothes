package v1

import (
	"api-gateway/api/models"
	"api-gateway/genproto/product_service"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Security 		BearerAuth
// @Summary 		Get User Saved Products
// @Description 	This API for getting saved product of user
// @Tags			list with user
// @Accept 			json
// @Produce 		json
// @Param			id path string true "User ID"
// @Success 		200 {object} []models.Product
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure			403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/user/save/{id} [GET]
func (h *HandlerV1) GetUserSavedProducts(c *gin.Context) {
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

	listProducts, err := h.Service.ProductService().GetSavedProductsByUserID(ctx, &product_service.GetWithUserID{
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
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

	c.JSON(http.StatusOK, models.ListProduct{
		Products: products,
		Total:    listProducts.TotalCount,
	})
}

// @Security 		BearerAuth
// @Summary 		Get User Likes Products
// @Description 	This API for getting like products of user
// @Tags			list with user
// @Accept 			json
// @Produce 		json
// @Param			id path string true "User ID"
// @Success 		200 {object} []models.Product
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure			403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/user/likes/{id} [GET]
func (h *HandlerV1) GetUserLikesProducts(c *gin.Context) {
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

	listProducts, err := h.Service.ProductService().GetWishlistByUserID(ctx, &product_service.GetWithUserID{
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
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

	c.JSON(http.StatusOK, models.ListProduct{
		Products: products,
		Total:    listProducts.TotalCount,
	})
}

// @Security 		BearerAuth
// @Summary 		Get User Ordered Products
// @Description 	This API for getting ordered products of user
// @Tags			list with user
// @Accept 			json
// @Produce 		json
// @Param			id path string true "User ID"
// @Success 		200 {object} []models.Product
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure			403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/user/orders/{id} [GET]
func (h *HandlerV1) GetUserOrderedProducts(c *gin.Context) {
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

	listProducts, err := h.Service.ProductService().GetOrderedProductsByUserID(ctx, &product_service.GetWithUserID{
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
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

	c.JSON(http.StatusOK, models.ListProduct{
		Products: products,
		Total:    listProducts.TotalCount,
	})
}
