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
// @Summary 		Recommandation Products
// @Description 	This API for searching products with name
// @Tags			another
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

	var products []models.Product
	for _, product := range listProducts.Products {
		products = append(products, models.Product{
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

	c.JSON(http.StatusOK, models.ListProduct{
		Products: products,
		Total:    listProducts.TotalCount,
	})
}

// @Security 		BearerAuth
// @Summary 		Search Products
// @Description 	This API for searching products with name
// @Tags			another
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

	var products []models.Product
	for _, product := range listProducts.Products {
		products = append(products, models.Product{
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

	c.JSON(http.StatusOK, models.ListProduct{
		Products: products,
		Total:    listProducts.TotalCount,
	})
}

// @Security 		BearerAuth
// @Summary 		List Disable Products
// @Description 	This API for getting list of disable product
// @Tags			another
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

	var orders []models.Order
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
		orders = append(orders, models.Order{
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

	c.JSON(http.StatusOK, models.ListOrder{
		Orders: orders,
		Total:  listProducts.TotalCount,
	})
}
