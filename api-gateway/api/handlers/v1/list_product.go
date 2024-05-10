package v1

import (
	"api-gateway/api/models"
	"api-gateway/genproto/product_service"
	"api-gateway/genproto/user_service"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Security 		BearerAuth
// @Summary 		Get Product Orders
// @Description 	This API for getting orders of product
// @Tags			list with product
// @Accept 			json
// @Produce 		json
// @Param			id path string true "Product ID"
// @Success 		200 {object} []models.Order
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure			403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/product/orders/{id} [GET]
func (h *HandlerV1) GetProductOrders(c *gin.Context) {
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

	productID := c.Param("id")

	listOrders, err := h.Service.ProductService().GetProductOrders(ctx, &product_service.GetWithID{
		Id: productID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	var orders []models.Order
	for _, order := range listOrders.Orders {
		user, err := h.Service.UserService().GetUser(ctx, &user_service.Filter{
			Filter: map[string]string{
				"id": order.UserId,
			},
		})
		if err != nil {
			log.Println("user deleted", order.UserId)
			continue
		}
		product, err := h.Service.ProductService().GetProductByID(ctx, &product_service.GetWithID{
			Id: order.ProductId,
		})
		if err != nil {
			log.Println("product deleted", order.ProductId)
			continue
		}
		orders = append(orders, models.Order{
			ID:     order.Id,
			Status: order.Status,
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

	c.JSON(http.StatusOK, models.ListOrder{
		Orders: orders,
		Total:  listOrders.TotalCount,
	})
}

// @Security 		BearerAuth
// @Summary 		Get Product Comments
// @Description 	This API for getting comments of product
// @Tags			list with product
// @Accept 			json
// @Produce 		json
// @Param			id path string true "Product ID"
// @Success 		200 {object} []models.Comment
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure			403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/product/comments/{id} [GET]
func (h *HandlerV1) GetProductComments(c *gin.Context) {
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

	productID := c.Param("id")

	listComments, err := h.Service.ProductService().GetProductComments(ctx, &product_service.GetWithID{
		Id: productID,
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
// @Summary 		Get Product Likes
// @Description 	This API for getting likes of product
// @Tags			list with product
// @Accept 			json
// @Produce 		json
// @Param			id path string true "Product ID"
// @Success 		200 {object} []models.Order
// @Failure 		400 {object} models.Like
// @Failure 		401 {object} models.Error
// @Failure			403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/product/likes/{id} [GET]
func (h *HandlerV1) GetProductLikes(c *gin.Context) {
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

	productID := c.Param("id")

	listLikes, err := h.Service.ProductService().GetProductLikes(ctx, &product_service.GetWithID{
		Id: productID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	var likes []models.Like
	for _, like := range listLikes.Likes {
		user, err := h.Service.UserService().GetUser(ctx, &user_service.Filter{
			Filter: map[string]string{
				"id": like.UserId,
			},
		})
		if err != nil {
			log.Println("user deleted", like.UserId)
			continue
		}
		product, err := h.Service.ProductService().GetProductByID(ctx, &product_service.GetWithID{
			Id: like.ProductId,
		})
		if err != nil {
			log.Println("product deleted", like.ProductId)
			continue
		}
		likes = append(likes, models.Like{
			ID: like.Id,
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

	c.JSON(http.StatusOK, models.ListLike{
		Likes: likes,
		Total: listLikes.TotalCount,
	})
}

// @Security 		BearerAuth
// @Summary 		Get Product Stars
// @Description 	This API for getting stars of product
// @Tags			list with product
// @Accept 			json
// @Produce 		json
// @Param			id path string true "Product ID"
// @Success 		200 {object} []models.Star
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure			403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/product/stars/{id} [GET]
func (h *HandlerV1) GetProductStars(c *gin.Context) {
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

	productID := c.Param("id")

	listStars, err := h.Service.ProductService().GetProductStars(ctx, &product_service.GetWithID{
		Id: productID,
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
