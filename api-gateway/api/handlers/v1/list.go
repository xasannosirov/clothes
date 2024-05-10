package v1

import (
	"api-gateway/api/models"
	"api-gateway/genproto/product_service"
	"api-gateway/genproto/user_service"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Security 		BearerAuth
// @Summary 		List Commnets
// @Description 	This API for getting list of comments
// @Tags			list
// @Accept 			json
// @Produce 		json
// @Param			page query uint64 true "Page"
// @Param			limit query uint64 true "Limit"
// @Success 		200 {object} []models.Comment
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure			403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/comments [GET]
func (h *HandlerV1) GetAllComments(c *gin.Context) {
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

	page := c.Query("page")
	limit := c.Query("limit")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	listComments, err := h.Service.ProductService().GetAllComments(ctx, &product_service.ListRequest{
		Page:  int64(pageInt),
		Limit: int64(limitInt),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	var comments []models.Comment
	for _, comment := range listComments.Comments {
		user, err := h.Service.UserService().GetUser(ctx, &user_service.Filter{
			Filter: map[string]string{
				"id": comment.UserId,
			},
		})
		if err != nil {
			log.Println("user deleted", comment.UserId)
			continue
		}
		product, err := h.Service.ProductService().GetProductByID(ctx, &product_service.GetWithID{
			Id: comment.ProductId,
		})
		if err != nil {
			log.Println("product deleted", comment.ProductId)
			continue
		}
		comments = append(comments, models.Comment{
			ID:      comment.Id,
			Comment: comment.Comment,
			User: models.User{
				Id:          user.Id,
				FirstName:   user.FirstName,
				LastName:    user.LastName,
				Email:       user.Email,
				Password:    user.Password,
				PhoneNumber: user.PhoneNumber,
				Gender:      user.Gender,
				Age:         user.Age,
				Refresh:     user.Refresh,
			},
			Product: models.Product{
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
			},
		})
	}

	c.JSON(http.StatusOK, models.ListComment{
		Comments: comments,
		Total:    listComments.TotalCount,
	})
}

// @Security 		BearerAuth
// @Summary 		List Stars
// @Description 	This API for getting list of stars
// @Tags			list
// @Accept 			json
// @Produce 		json
// @Param			page query uint64 true "Page"
// @Param			limit query uint64 true "Limit"
// @Success 		200 {object} []models.Star
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure			403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/stars [GET]
func (h *HandlerV1) GetAllStars(c *gin.Context) {
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

	page := c.Query("page")
	limit := c.Query("limit")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	listStars, err := h.Service.ProductService().GetAllStars(ctx, &product_service.ListRequest{
		Page:  int64(pageInt),
		Limit: int64(limitInt),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	var stars []models.Star
	for _, star := range listStars.Stars {
		user, err := h.Service.UserService().GetUser(ctx, &user_service.Filter{
			Filter: map[string]string{
				"id": star.UserId,
			},
		})
		if err != nil {
			log.Println("user deleted", star.UserId)
			continue
		}
		product, err := h.Service.ProductService().GetProductByID(ctx, &product_service.GetWithID{
			Id: star.ProductId,
		})
		if err != nil {
			log.Println("product deleted", star.ProductId)
			continue
		}
		stars = append(stars, models.Star{
			ID:   star.Id,
			Star: int16(star.Star),
			User: models.User{
				Id:          user.Id,
				FirstName:   user.FirstName,
				LastName:    user.LastName,
				Email:       user.Email,
				Password:    user.Password,
				PhoneNumber: user.PhoneNumber,
				Gender:      user.Gender,
				Age:         user.Age,
				Refresh:     user.Refresh,
			},
			Product: models.Product{
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
			},
		})
	}

	c.JSON(http.StatusOK, models.ListStar{
		Stars: stars,
		Totol: listStars.TotalCount,
	})
}
