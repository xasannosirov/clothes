package v1

import (
	"api-gateway/api/models"
	"api-gateway/genproto/media_service"
	"api-gateway/genproto/product_service"
	"api-gateway/internal/pkg/query_parameter"
	"context"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Security 		BearerAuth
// @Summary 		Create Product
// @Description 	This API for create a new product
// @Tags 			products
// @Produce 		json
// @Accept 			json
// @Param 			product body models.ProductReq true "Create Product Model"
// @Success			201 {object} models.ProductCreateResponse
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

	for i := 0; i < len(body.Color); i++ {
		isAlpha := func (s string) bool {
			for _, char := range s {
				if !unicode.IsLetter(char) {
					return false
				}
			}
			return true
		}(body.Color[i])

		if len(body.Color[i]) < 3 || !isAlpha{
			c.JSON(http.StatusBadRequest, models.Error{
				Message: "color must be more than 3 only letters long",
			})
			log.Println("color must be more than 3 only letters long")
			return
		}
	}

	for i := 0; i < len(body.Size); i++ {
		romanRegex := `^(?i)(M{0,4}(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3}))$`
		matched, err := func(s string) (matched bool, err error) {
			re, err := regexp.Compile(romanRegex)
			if err != nil {
				return false, err
			}
			return re.MatchString(s), nil
		}(body.Size[i])
		if err != nil {
			return
		}
		if !matched {
			c.JSON(http.StatusBadRequest, models.Error{
				Message: "product size should be in Roman numerals only",
			})
			log.Println("product size should be in Roman numerals only")
			return
		}
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
		ProductSize: body.Size,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, models.ProductCreateResponse{
		ProductID: createdProductResponse.Id,
	})
}

// @Security 		BearerAuth
// @Summary 		Update Product
// @Description 	This API for updating a product
// @Tags 			products
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
		ProductSize: body.Size,
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
		Size:        updatedProduct.ProductSize,
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
// @Tags 			products
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Product ID"
// @Success			200 {object} bool
// @Failure 		404 {object} models.Error
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
		c.JSON(http.StatusInternalServerError, models.Error{
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
		c.JSON(http.StatusNotFound, models.Error{
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
// @Tags 			products
// @Produce 		json
// @Accept 			json
// @Param 			id path string true "Product ID"
// @Success			200 {object} models.Product
// @Failure 		404 {object} models.Error
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
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	productID := c.Param("id")

	product, err := h.Service.ProductService().GetProduct(ctx, &product_service.GetWithID{
		Id: productID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	media, err := h.Service.MediaService().Get(ctx, &media_service.MediaWithID{
		Id: productID,
	})
	
	if err != nil {
		log.Println(err.Error())
	}
	var imagesURL []string
	for _, imageUrl := range media.Images{
		imagesURL = append(imagesURL, imageUrl.ImageUrl)
	}

	if len(media.Images) == 0 {
		media.Images = append(media.Images, &media_service.Media{
			ImageUrl: "",
		})
	}

	c.JSON(http.StatusOK, models.Product{
		ID:          productID,
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
		ImageURL:    imagesURL,
	})
}

// @Security 		BearerAuth
// @Summary 		List Products
// @Description 	This API for getting list of products
// @Tags 			products
// @Produce 		json
// @Accept 			json
// @Param 			page query uint64 true "Page"
// @Param 			limit query uint64 true "Limit"
// @Param 			name query string false "Product Name"
// @Success			200 {object} models.ListProduct
// @Failure 		404 {object} models.Error
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
	name := c.Query("name")

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

	listProducts := &product_service.ListProduct{}
	if name == "" {
		listProducts, err = h.Service.ProductService().ListProducts(ctx, &product_service.ListRequest{
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
	} else {
		listProducts, err = h.Service.ProductService().SearchProduct(ctx, &product_service.SearchRequest{
			Page:  uint64(pageInt),
			Limit: uint64(limitInt),
			Params: map[string]string{
				"name": name,
			},
		})
		if err != nil {
			c.JSON(http.StatusNotFound, models.Error{
				Message: err.Error(),
			})
			log.Println(err.Error())
			return
		}
	}

	var products []models.Product
	for _, product := range listProducts.Products {
		media, err := h.Service.MediaService().Get(ctx, &media_service.MediaWithID{
			Id: product.Id,
		})
		if err != nil {
			log.Println(err.Error())
		}

		var imagesURL []string
		for _, imageUrl := range media.Images {
			imagesURL = append(imagesURL, imageUrl.ImageUrl)
		}

		if len(media.Images) == 0 {
			media.Images = append(media.Images, &media_service.Media{
				ImageUrl: "",
			})
		}

		products = append(products, models.Product{
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
			ImageURL:    imagesURL,
		})
	}

	response := models.ListProduct{}
	response.Products = products
	response.Total = listProducts.TotalCount

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		Discount Products
// @Description 	This API returns discount products
// @Tags			products
// @Accept 			json
// @Produce 		json
// @Param			page query uint64 true "Page"
// @Param			limit query uint64 true "Limit"
// @Success			200 {object} models.ListProduct
// @Failure			401 {object} models.Error
// @Failure			403 {object} models.Error
// @Failure			404 {object} models.Error
// @Failure			500 {object} models.Error
// @Router			/v1/products/discount [GET]
func (h *HandlerV1) GetDicountProducts(c *gin.Context) {
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

	parameters := query_parameter.New(c.Request.URL.Query())

	products, err := h.Service.ProductService().GetDiscountProducts(ctx, &product_service.ListRequest{
		Page:  int64(parameters.GetPage()),
		Limit: int64(parameters.GetLimit()),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: models.NotFoundMessage,
		})
		log.Println(err.Error())
		return
	}

	var response models.ListProduct
	for _, serviceProduct := range products.Products {
		media, err := h.Service.MediaService().Get(ctx, &media_service.MediaWithID{
			Id: serviceProduct.Id,
		})
		if err != nil {
			log.Println(err.Error())
		}

		var imagesURL []string
		for _, imageUrl := range media.Images {
			imagesURL = append(imagesURL, imageUrl.ImageUrl)
		}

		if len(media.Images) == 0 {
			media.Images = append(media.Images, &media_service.Media{
				ImageUrl: "",
			})
		}
		response.Products = append(response.Products, models.Product{
			ID:          serviceProduct.Id,
			Name:        serviceProduct.Name,
			Category:    serviceProduct.Category,
			Description: serviceProduct.Description,
			MadeIn:      serviceProduct.MadeIn,
			Color:       serviceProduct.Color,
			Size:        serviceProduct.ProductSize,
			Count:       serviceProduct.Count,
			Cost:        float64(serviceProduct.Cost),
			Discount:    float64(serviceProduct.Discount),
			AgeMin:      serviceProduct.AgeMin,
			AgeMax:      serviceProduct.AgeMax,
			ForGender:   serviceProduct.ForGender,
			ImageURL:    imagesURL,
		})
	}
	response.Total = products.TotalCount

	c.JSON(http.StatusOK, response)
}
