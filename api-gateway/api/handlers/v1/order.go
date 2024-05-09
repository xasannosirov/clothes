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
// @Summary 		Create Order
// @Description 	This API for create a new order
// @Tags 			orders
// @Produce 		json
// @Accept 			json
// @Param 			order body models.OrderReq true "Create Order Model"
// @Success			201 {object} string
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/order [POST]
func (h *HandlerV1) CreateOrder(c *gin.Context) {
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

	createdOrderResponse, err := h.Service.ProductService().CreateOrder(ctx, &product_service.Order{
		ProductId: body.ProductID,
		UserId:    body.UserID,
		Status:    "ordered",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	product, err := h.Service.ProductService().GetProductByID(ctx, &product_service.GetWithID{
		Id: body.ProductID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	_, err = h.Service.ProductService().UpdateProduct(ctx, &product_service.Product{
		Id:          body.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		MadeIn:      product.MadeIn,
		Color:       product.Color,
		Count:       product.Count - 1,
		Cost:        product.Cost,
		Discount:    product.Discount,
		AgeMin:      product.AgeMin,
		AgeMax:      product.AgeMax,
		ForGender:   product.ForGender,
		Size_:       product.Size_,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, createdOrderResponse.Id)
}

// @Security 		BearerAuth
// @Summary 		Cancel Order
// @Description 	This API for cancelling a order
// @Tags 			orders
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Order ID"
// @Success			200 {object} string
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/order/{id} [DELETE]
func (h *HandlerV1) CancelOrder(c *gin.Context) {
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

	orderId := c.Param("id")

	responseStatus, err := h.Service.ProductService().CancelOrder(ctx, &product_service.GetWithID{
		Id: orderId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, responseStatus.Status)
}

// @Security 		BearerAuth
// @Summary 		Get Order
// @Description 	This API for getting a order with id
// @Tags 			orders
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Order ID"
// @Success			200 {object} models.Order
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/order/{id} [GET]
func (h *HandlerV1) GetOrder(c *gin.Context) {
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

	orderId := c.Param("id")

	order, err := h.Service.ProductService().GetOrderByID(ctx, &product_service.GetWithID{
		Id: orderId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	user, err := h.Service.UserService().GetUser(ctx, &user_service.Filter{
		Filter: map[string]string{
			"id": order.UserId,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	product, err := h.Service.ProductService().GetProductByID(ctx, &product_service.GetWithID{
		Id: order.ProductId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, models.Order{
		ID: order.Id,
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
		User: models.User{
			Id:          user.Id,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			Password:    user.Password,
			PhoneNumber: user.PhoneNumber,
			Gender:      user.Gender,
			Age:         user.Age,
			Role:        user.Role,
			Refresh:     user.Refresh,
		},
		Status: order.Status,
	})
}

// @Security 		BearerAuth
// @Summary 		List Orders
// @Description 	This API for getting list of orders
// @Tags 			orders
// @Produce 		json
// @Accept 			json
// @Param 			page query uint64 true "Page"
// @Param			limit query uint64 true "Limit"
// @Success			200 {object} []models.Order
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/orders [GET]
func (h *HandlerV1) ListOrders(c *gin.Context) {
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

	listOrders, err := h.Service.ProductService().GetAllOrders(ctx, &product_service.ListRequest{
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

	var orders []models.Order
	for _, order := range listOrders.Orders {
		user, err := h.Service.UserService().GetUser(ctx, &user_service.Filter{
			Filter: map[string]string{
				"id": order.UserId,
			},
		})
		if err != nil {
			continue
		}
		product, err := h.Service.ProductService().GetProductByID(ctx, &product_service.GetWithID{
			Id: order.ProductId,
		})
		if err != nil {
			continue
		}
		orders = append(orders, models.Order{
			ID: order.Id,
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
			User: models.User{
				Id:          user.Id,
				FirstName:   user.FirstName,
				LastName:    user.LastName,
				Email:       user.Email,
				Password:    user.Password,
				PhoneNumber: user.PhoneNumber,
				Gender:      user.Gender,
				Age:         user.Age,
				Role:        user.Role,
				Refresh:     user.Refresh,
			},
			Status: order.Status,
		})
	}

	c.JSON(http.StatusOK, models.ListOrder{
		Orders: orders,
		Total:  listOrders.TotalCount,
	})
}
