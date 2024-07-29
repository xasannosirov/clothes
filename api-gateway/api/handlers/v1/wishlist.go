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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Security 		BearerAuth
// @Summary 		Like Product
// @Description 	This API for save likes a product by user
// @Tags 			wishlist
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Product ID"
// @Success			201 {object} bool
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/like/{id} [POST]
func (h *HandlerV1) LikeProduct(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	productID := c.Param("id")

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

	userId, statusCode := regtool.GetIdFromToken(c.Request, &h.Config)
	if statusCode != 0 {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "oops something went wrong",
		})
	}

	likeResponse, err := h.Service.ProductService().LikeProduct(ctx, &product_service.Like{
		ProductId: productID,
		UserId:    userId,
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
// @Summary 		User Wishlist
// @Description 	This API for getting wishlist for user
// @Tags 			wishlist
// @Produce 		json
// @Accept 			json
// @Param 			page query uint64 true "Page"
// @Param 			limit query uint64 true "Limit"
// @Success			200 {object} models.ListProduct
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		404 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/wishlist [GET]
func (h *HandlerV1) UserWishlist(c *gin.Context) {
	userID, statusCode := regtool.GetIdFromToken(c.Request, &h.Config)
	if statusCode == http.StatusUnauthorized {
		c.JSON(http.StatusUnauthorized, models.Error{
			Message: models.NoAccessMessage,
		})
		return
	}
	parameters := query_parameter.New(c.Request.URL.Query())

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

	products, err := h.Service.ProductService().UserWishlist(ctx, &product_service.SearchRequest{
		Page:  parameters.GetPage(),
		Limit: parameters.GetLimit(),
		Params: map[string]string{
			"user_id": userID,
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
			var imagesURL []string
			if err != nil {
				log.Println(err.Error())
			} else {
				for _, imageUrl := range media.Images {
					imagesURL = append(imagesURL, imageUrl.ImageUrl)
				}
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
			basketStatus, err := h.Service.ProductService().IsUnique(ctx, &product_service.IsUniqueReq{
				TableName: "basket",
				UserId:    userId,
				ProductId: product.Id,
			})
			if err != nil {
				if strings.Contains(err.Error(), "no rows") {
					basketStatus = &product_service.MoveResponse{
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
			var imagesURL []string
			if err != nil {
				log.Println(err.Error())
			} else {
				for _, imageUrl := range media.Images {
					imagesURL = append(imagesURL, imageUrl.ImageUrl)
				}
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
