package v1

import (
	"api-gateway/api/models"
	"api-gateway/genproto/media_service"
	"api-gateway/genproto/product_service"
	"api-gateway/internal/pkg/regtool"
	"api-gateway/internal/pkg/validation"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
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
// @Success			201 {object} bool
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
			Message: "you needs register or login",
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

	c.JSON(http.StatusOK, basket.Status)
}

// @Security 		BearerAuth
// @Summary 		Get Basket
// @Description 	This API for getting a Basket with id
// @Tags 			basket
// @Produce 		json
// @Accept 			json
// @Param 			id query string false "User ID"
// @Success			200 {object} []models.Product
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/user-baskets [GET]
func (h *HandlerV1) GetUserBaskets(c *gin.Context) {
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

	userID := c.Query("id")

	var status int
	if !validation.ValidateUUID(userID) {
		userID, status = regtool.GetIdFromToken(c.Request, &h.Config)
		if status != 0 {
			c.JSON(http.StatusUnauthorized, models.Error{
				Message: models.TokenInvalidMessage,
			})
			log.Println(models.TokenInvalidMessage)
			return
		}
	}

	products, err := h.Service.ProductService().GetUserBaskets(ctx, &product_service.GetWithID{
		Id: userID,
	})
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, []models.Product{})
		log.Println(err.Error())
		return
	}

	var response []models.Product
	for _, product := range products.Products {
		media, err := h.Service.MediaService().Get(ctx, &media_service.MediaWithID{
			Id: product.Id,
		})
		if err != nil {
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
		likeStatus, err := h.Service.ProductService().IsUnique(ctx, &product_service.IsUniqueReq{
			TableName: "wishlist",
			UserId:    userID,
			ProductId: product.Id,
		})
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				likeStatus = &product_service.MoveResponse{
					Status: false,
				}
			} else {
				c.JSON(http.StatusBadRequest, models.Error{
					Message: err.Error(),
				})
				log.Println(err.Error())
				return
			}
		}

		response = append(response, models.Product{
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
			Basket:      true,
			Liked:       likeStatus.Status,
			ImageURL:    imagesURL,
		})
	}

	c.JSON(http.StatusOK, response)
}
