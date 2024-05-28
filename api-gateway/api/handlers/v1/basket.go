package v1

import (
	"api-gateway/api/models"
	"api-gateway/genproto/media_service"
	"api-gateway/genproto/product_service"
	"api-gateway/internal/pkg/regtool"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Security 		BearerAuth
// @Summary 		Save To Basket
// @Description 	This API for create a new basket for product
// @Tags 			basket
// @Produce 		json
// @Accept 			json
// @Param 			order body models.BasketCeateReq true "Create Basket Model"
// @Success			201 {object} models.CreateResponse
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/basket [POST]
func (h *HandlerV1) SaveToBasket(c *gin.Context) {
	var (
		body        models.BasketCeateReq
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
	userId, statusCode := regtool.GetIdFromToken(c.Request, &h.Config)
	if statusCode != 0 {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "oops something went wrong",
		})
	}
	basket, err := h.Service.ProductService().SaveToBasket(ctx, &product_service.BasketCreateReq{
		ProductId: body.ProductId,
		UserId:    userId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.CreateResponse{
		UserID: basket.Id,
	})
}


// @Security 		BearerAuth
// @Summary 		Get Basket
// @Description 	This API for getting a Basket with id
// @Tags 			basket
// @Produce 		json
// @Accept 			json
// @Param 			page query uint64 true "Page"
// @Param 			limit query uint64 true "Limit"
// @Param 			id query string true "Basket ID"
// @Success			200 {object} models.Basket
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/basket [GET]
func (h *HandlerV1) GetBasketProduct(c *gin.Context) {
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

	id := c.Query("id")
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

	basket, err := h.Service.ProductService().GetBasket(ctx, &product_service.BasketGetReq{
		UserId: id,
		Page: int64(pageInt),
		Limit: int64(limitInt),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	products := []models.Product{}
	for _, product := range basket.Product {

		media, err := h.Service.MediaService().Get(ctx, &media_service.MediaWithID{
			Id: product.Id,
		})
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: err.Error(),
			})
			log.Println(err.Error())
			return
		}
		var imagesURL []string
		for _, imageUrl := range media.Images {
			imagesURL = append(imagesURL, imageUrl.ImageUrl)
		}



		products = append(products, models.Product{
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
			ImageURL:    imagesURL,
		})
	}

	c.JSON(http.StatusOK, models.Basket{
		UserId:    basket.UserId,
		ProductId: products,
		TotalCount: basket.TotalCount,
	})
}
