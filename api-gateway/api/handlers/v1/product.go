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
// @Summary 		Create Product
// @Description 	This API for create a new product
// @Tags 			products and orders
// @Produce 		json
// @Accept 			json
// @Param 			product body models.ProductReq true "Create Product Model"
// @Success			201 {object} string
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/product [POST]
func (h *HandlerV1) CreateProduct(c *gin.Context) {
	var (
		body        models.ProductReq
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

	createdProductResponse, err := h.Service.ProductService().CreateProduct(ctx, &product_service.Product{
		Name:        body.Name,
		Description: body.Description,
		Category:    body.Category,
		MadeIn:      body.MadeIn,
		Color:       body.Color,
		Count:       body.Count,
		Cost:        float32(body.Cost),
		Discount:    float32(body.Discount),
		AgeMin:      body.AgeMin,
		AgeMax:      body.AgeMax,
		ForGender:   body.ForGender,
		Size_:       body.Size,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, createdProductResponse.Id)
}

// @Security 		BearerAuth
// @Summary 		Update Product
// @Description 	This API for updating a product
// @Tags 			products and orders
// @Produce 		json
// @Accept 			json
// @Param 			product body models.Product true "Update Product Model"
// @Success			200 {object} models.Product
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/product [PUT]
func (h *HandlerV1) UpdateProduct(c *gin.Context) {
	var (
		body        models.Product
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

	updatedProduct, err := h.Service.ProductService().UpdateProduct(ctx, &product_service.Product{
		Id:          body.ID,
		Name:        body.Name,
		Description: body.Description,
		Category:    body.Category,
		MadeIn:      body.MadeIn,
		Color:       body.Color,
		Count:       body.Count,
		Cost:        float32(body.Cost),
		Discount:    float32(body.Discount),
		AgeMin:      body.AgeMin,
		AgeMax:      body.AgeMax,
		ForGender:   body.ForGender,
		Size_:       body.Size,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.Product{
		ID:          updatedProduct.Id,
		Name:        updatedProduct.Name,
		Category:    updatedProduct.Category,
		Description: updatedProduct.Description,
		MadeIn:      updatedProduct.MadeIn,
		Color:       updatedProduct.Color,
		Size:        updatedProduct.Size_,
		Count:       updatedProduct.Count,
		Cost:        float64(updatedProduct.Cost),
		Discount:    float64(updatedProduct.Discount),
		AgeMin:      updatedProduct.AgeMin,
		AgeMax:      updatedProduct.AgeMax,
		ForGender:   updatedProduct.ForGender,
	})
}

// @Security 		BearerAuth
// @Summary 		Delete Product
// @Description 	This API for deleting a product with product_id
// @Tags 			products and orders
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Product ID"
// @Success			200 {object} bool
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/product/{id} [DELETE]
func (h *HandlerV1) DeleteProduct(c *gin.Context) {
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

	responseStatus, err := h.Service.ProductService().DeleteProduct(ctx, &product_service.GetWithID{
		Id: productID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, responseStatus.Status)
}

// @Security 		BearerAuth
// @Summary 		Get Product
// @Description 	This API for getting a product with product_id
// @Tags 			products and orders
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Product ID"
// @Success			200 {object} models.Product
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/product/{id} [GET]
func (h *HandlerV1) GetProduct(c *gin.Context) {
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

	product, err := h.Service.ProductService().GetProductByID(ctx, &product_service.GetWithID{
		Id: productID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.Product{
		ID:          productID,
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

// @Security 		BearerAuth
// @Summary 		List Products
// @Description 	This API for getting list of products
// @Tags 			products and orders
// @Produce 		json
// @Accept 			json
// @Param 			page query uint64 true "Page"
// @Param 			limit query uint64 true "Limit"
// @Success			200 {object} models.Product
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Faulure 		500 {object} models.Error
// @Router 			/v1/products [GET]
func (h *HandlerV1) ListProducts(c *gin.Context) {
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

	listProducts, err := h.Service.ProductService().GetAllProducts(ctx, &product_service.ListRequest{
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

	var response []models.Product
	for _, product := range listProducts.Products {
		response = append(response, models.Product{
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

	c.JSON(http.StatusCreated, response)
}

// @Security 		BearerAuth
// @Summary 		Create Order
// @Description 	This API for create a new order
// @Tags 			products and orders
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

	c.JSON(http.StatusCreated, createdOrderResponse.Id)
}

// @Security 		BearerAuth
// @Summary 		Cancel Order
// @Description 	This API for cancelling a order
// @Tags 			products and orders
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
// @Tags 			products and orders
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
// @Tags 			products and orders
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

	orders, err := h.Service.ProductService().GetAllOrders(ctx, &product_service.ListRequest{
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

	var response []models.Order
	for _, order := range orders.Orders {
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
		response = append(response, models.Order{
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

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		Like Product
// @Description 	This API for save likes a product by user
// @Tags 			products and orders
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
// @Tags			products and orders
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
// @Tags 			products and orders
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
// @Tags 			products and orders
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

// @Security 		BearerAuth
// @Summary 		List Commnets
// @Description 	This API for getting list of comments
// @Tags			products and orders
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

	var response []models.Comment
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
		response = append(response, models.Comment{
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
				Role:        user.Role,
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

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		List Stars
// @Description 	This API for getting list of stars
// @Tags			products and orders
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

	var response []models.Star
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
		response = append(response, models.Star{
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
				Role:        user.Role,
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

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		Get Product Orders
// @Description 	This API for getting orders of product
// @Tags			products and orders
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

	var response []models.Order
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
		response = append(response, models.Order{
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
				Role:        user.Role,
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

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		Get Product Comments
// @Description 	This API for getting comments of product
// @Tags			products and orders
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

	var response []models.Comment
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
		response = append(response, models.Comment{
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
				Role:        user.Role,
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

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		Get Product Likes
// @Description 	This API for getting likes of product
// @Tags			products and orders
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

	var response []models.Like
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
		response = append(response, models.Like{
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
				Role:        user.Role,
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

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		Get Product Stars
// @Description 	This API for getting stars of product
// @Tags			products and orders
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

	var response []models.Star
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
		response = append(response, models.Star{
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
				Role:        user.Role,
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

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		Get User Saved Products
// @Description 	This API for getting saved product of user
// @Tags			products and orders
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

	var response []models.Product
	for _, product := range listProducts.Products {
		response = append(response, models.Product{
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

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		Get User Likes Products
// @Description 	This API for getting like products of user
// @Tags			products and orders
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

	var response []models.Product
	for _, product := range listProducts.Products {
		response = append(response, models.Product{
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

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		Get User Ordered Products
// @Description 	This API for getting ordered products of user
// @Tags			products and orders
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

	var response []models.Product
	for _, product := range listProducts.Products {
		response = append(response, models.Product{
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

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		Recommandation Products
// @Description 	This API for searching products with name
// @Tags			products and orders
// @Accept 			json
// @Produce 		json
// @Param			gender query string true "Gender"
// @Param			age query uint64 true "Age"
// @Success 		200 {object} []models.Product
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure			403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/recommendation [GET]
func (h *HandlerV1) RecommendProducts(c *gin.Context) {
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

	gender := c.Query("gender")
	age := c.Query("age")
	ageInt, err := strconv.Atoi(age)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	listProducts, err := h.Service.ProductService().Recommendation(ctx, &product_service.Recom{
		Gender: gender,
		Age:    int64(ageInt),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	var response []models.Product
	for _, product := range listProducts.Products {
		response = append(response, models.Product{
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

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		Search Products
// @Description 	This API for searching products with name
// @Tags			products and orders
// @Accept 			json
// @Produce 		json
// @Param			name path string true "Product Name"
// @Success 		200 {object} []models.Product
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure			403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/search/{name} [GET]
func (h *HandlerV1) SearchProduct(c *gin.Context) {
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

	name := c.Param("name")

	listProducts, err := h.Service.ProductService().SearchProduct(ctx, &product_service.Filter{
		Name: name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	var response []models.Product
	for _, product := range listProducts.Products {
		response = append(response, models.Product{
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

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		List Disable Products
// @Description 	This API for getting list of disable product
// @Tags			products and orders
// @Accept 			json
// @Produce 		json
// @Param			page query uint64 true "Page"
// @Param			limit query uint64 true "Limit"
// @Success 		200 {object} []models.Order
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure			403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/disable-orders [GET]
func (h *HandlerV1) GetDisableProducts(c *gin.Context) {
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

	listProducts, err := h.Service.ProductService().GetDisableProducts(ctx, &product_service.ListRequest{
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

	var response []models.Order
	for _, order := range listProducts.Orders {
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
		response = append(response, models.Order{
			ID: order.Id,
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

	c.JSON(http.StatusOK, response)
}
