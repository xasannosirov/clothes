package v1

import (
	"api-gateway/api/models"
	"api-gateway/genproto/product_service"
	"api-gateway/internal/pkg/regtool"
	"context"
	"log"
	"net/http"
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
// @Summary 		Delete Basket
// @Description 	This API for delete a basket with id
// @Tags 			basket
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Product ID"
// @Success			200 {object} bool
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

	userId, statusCode := regtool.GetIdFromToken(c.Request, &h.Config)
	if statusCode != 0 {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "oops something went wrong",
		})
	}
	id := c.Param("id")

	status, err := h.Service.ProductService().DeleteFromBasket(ctx, &product_service.DeleteBasket{
		UserId: userId,
		ProductId: id,
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
// @Success			200 {object} models.Basket
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
	basket, err := h.Service.ProductService().GetBasket(ctx, &product_service.GetWithID{
		Id: id,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.Basket{
		UserId:    basket.UserId,
		ProductId: basket.ProductId,
		Count:     int64(basket.Count),
	})
}
