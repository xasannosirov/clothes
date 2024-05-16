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
// @Summary 		Save To Basket
// @Description 	This API for create a new basket for product
// @Tags 			basket
// @Produce 		json
// @Accept 			json
// @Param 			order body models.BasketCeateReq true "Create Basket Model"
// @Success			201 {object} models.Basket
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
	userId, statusCode := GetIdFromToken(c.Request, &h.Config)
	if statusCode != 0 {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "oops something went wrong",
		})
	}
	basket, err := h.Service.ProductService().SaveToBasket(ctx, &product_service.Basket{
		Id:   uuid.NewString(),
		ProductId: body.ProductId,
		UserId: userId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.Basket{
		Id:   basket.Id,
		UserId: userId,
		ProductId: body.ProductId,
	})
}


// @Security 		BearerAuth
// @Summary 		Delete Basket 
// @Description 	This API for delete a basket with id
// @Tags 			basket
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Basket ID"
// @Success			201 {object} bool
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		404 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/basket/{id} [DELETE]
func (h *HandlerV1) DeleteFromBasket(c *gin.Context) {
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

	id := c.Param("id")
	filter := map[string]string{
		"id": id,
	}
	status, err := h.Service.ProductService().DeleteFromBasket(ctx, &product_service.RequestBasket{
		Filter: filter,
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
// @Summary 		Get Basket
// @Description 	This API for getting a Basket with id
// @Tags 			basket
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Basket ID"
// @Success			201 {object} models.Basket
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/basket/{id} [GET]
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

	id := c.Param("id")
	filter := map[string]string{
		"id": id,
	}
	basket, err := h.Service.ProductService().GetBasketProduct(ctx, &product_service.RequestBasket{
		Filter: filter,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.Basket{
		Id: basket.Id,
		UserId: basket.UserId,
		ProductId: basket.ProductId,
	})
}

// @Security 		BearerAuth
// @Summary 		List Basket 
// @Description 	This API for getting Basket's list
// @Tags 			basket
// @Produce 		json
// @Accept 			json
// @Param 			page query uint64 true "Page"
// @Param 			limit query uint64 true "Limit"
// @Success			201 {object} string
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/baskets [GET]
func (h *HandlerV1) GetBasketProducts(c *gin.Context) {
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

	basketsList, err := h.Service.ProductService().GetBasketProducts(ctx, &product_service.ListBasketRequest{
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

	var baskets []models.Basket
	for _, basket := range basketsList.Baskets {
		baskets = append(baskets, models.Basket{
			Id: basket.Id,
			UserId: basket.UserId,
			ProductId: basket.ProductId,
		})
	}

	c.JSON(http.StatusOK, models.ListBasket{
		Baskets: baskets,
		Total:     basketsList.TotalCount,
	})
}
