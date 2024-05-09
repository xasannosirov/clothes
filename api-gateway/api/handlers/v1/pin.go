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
// @Summary 		Like Product
// @Description 	This API for save likes a product by user
// @Tags 			pin
// @Produce 		json
// @Accept 			json
// @Param 			like body models.LikeReq true "Like Product Model"
// @Success			201 {object} bool
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/like-product [POST]
func (h *HandlerV1) LikeProduct(c *gin.Context) {
	var (
		body        models.OrderReq
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

	likeResponse, err := h.Service.ProductService().LikeProduct(ctx, &product_service.Like{
		ProductId: body.ProductID,
		UserId:    body.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, likeResponse.Status)
}

// @Security 		BearerAuth
// @Summary 		Save Product
// @Description 	This API for saving a order by user
// @Tags			pin
// @Produce 		json
// @Accept 			json
// @Param 			save body models.SaveReq true "Save Product Model"
// @Success			201 {object} bool
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/save-product [POST]
func (h *HandlerV1) SaveProduct(c *gin.Context) {
	var (
		body        models.SaveReq
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

	saveResponse, err := h.Service.ProductService().SaveProduct(ctx, &product_service.Save{
		ProductId: body.ProductID,
		UserId:    body.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, saveResponse.Status)
}

// @Security 		BearerAuth
// @Summary 		Star Product
// @Description 	This API for valuation a product with star by user
// @Tags 			pin
// @Produce 		json
// @Accept 			json
// @Param 			star body models.StarReq true "Star Product Model"
// @Success			201 {object} bool
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/star-product [POST]
func (h *HandlerV1) StarToProduct(c *gin.Context) {
	var (
		body        models.StarReq
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

	starProduct, err := h.Service.ProductService().StarProduct(ctx, &product_service.Star{
		ProductId: body.ProductID,
		UserId:    body.UserID,
		Star:      int64(body.Star),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, starProduct.Status)
}

// @Security 		BearerAuth
// @Summary 		Comment to Product
// @Description 	This API for write a comment to product by user
// @Tags 			pin
// @Produce 		json
// @Accept 			json
// @Param 			comment body models.CommentReq true "Comment Product Model"
// @Success			201 {object} bool
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/comment-product [POST]
func (h *HandlerV1) CommentToProduct(c *gin.Context) {
	var (
		body        models.CommentReq
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

	commentProduct, err := h.Service.ProductService().CommentToProduct(ctx, &product_service.Comment{
		ProductId: body.ProductID,
		UserId:    body.UserID,
		Comment:   body.Comment,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, commentProduct.Status)
}
