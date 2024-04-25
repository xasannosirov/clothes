package postgresql

import (
	"context"
	"log"
	"product-service/internal/entity"
	"product-service/internal/infrastructure/repository"
	"testing"
	"time"

	pbp "product-service/genproto/product_service"
	"product-service/internal/pkg/config"
	db "product-service/internal/pkg/postgres"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ProductRepositorySuiteTest struct {
	suite.Suite
	Repository repository.Product
}

func (p *ProductRepositorySuiteTest) SetupSuite() {
	pgPoll, err := db.New(config.New())
	if err != nil {
		log.Fatal("Error while connecting to database with suite test")
		return
	}

	p.Repository = NewProductsRepo(pgPoll)
}

func (p *ProductRepositorySuiteTest) TestCreateProduct() {
	productReq := &pbp.Product{
		Id:             uuid.NewString(),
		Name:           "test",
		Description:    "test",
		Category:       "test",
		MadeIn:         "test",
		Color:          "test",
		Cost:           1,
		Count:          1,
		Discount:       1,
		AgeMin:         1,
		AgeMax:         1,
		TemperatureMin: 1,
		TemperatureMax: 1,
		ForGender:      "test",
		Size_:          1,
	}

	createProduct, err := p.Repository.CreateProduct(context.Background(), &entity.Product{
		Id:             productReq.Id,
		Name:           productReq.Name,
		Description:    productReq.Description,
		Category:       productReq.Category,
		MadeIn:         productReq.MadeIn,
		Color:          productReq.Color,
		Cost:           productReq.Cost,
		Count:          productReq.Count,
		Discount:       productReq.Discount,
		AgeMin:         productReq.AgeMin,
		AgeMax:         productReq.AgeMax,
		TemperatureMin: productReq.TemperatureMin,
		TemperatureMax: productReq.TemperatureMax,
		ForGender:      productReq.ForGender,
		Size:           productReq.Size_,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	p.Suite.NoError(err)
	p.Suite.NotNil(createProduct)
	p.Suite.Equal(productReq.Id, createProduct.Id)
	p.Suite.Equal(productReq.Description, createProduct.Description)
	p.Suite.Equal(productReq.Name, createProduct.Name)
	p.Suite.Equal(productReq.Category, createProduct.Category)
	p.Suite.Equal(productReq.MadeIn, createProduct.MadeIn)
	p.Suite.Equal(productReq.Color, createProduct.Color)
	p.Suite.Equal(productReq.Cost, createProduct.Cost)
	p.Suite.Equal(productReq.Count, createProduct.Count)
	p.Suite.Equal(productReq.Discount, createProduct.Discount)
	p.Suite.Equal(productReq.AgeMin, createProduct.AgeMin)
	p.Suite.Equal(productReq.AgeMax, createProduct.AgeMax)
	p.Suite.Equal(productReq.TemperatureMin, createProduct.TemperatureMin)
	p.Suite.Equal(productReq.TemperatureMax, createProduct.TemperatureMax)
	p.Suite.Equal(productReq.ForGender, createProduct.ForGender)
	p.Suite.Equal(productReq.Size_, createProduct.Size)
	p.Suite.NotNil(createProduct.CreatedAt)
	p.Suite.NotNil(createProduct.UpdatedAt)
}

func (p *ProductRepositorySuiteTest) TestGetProduct() {
	filter := make(map[string]string)
	filter["id"] = "fba19244-0b60-4a0c-a042-bb1e2a5fdfd5"
	getProduct, err := p.Repository.GetProduct(context.Background(), filter)
	p.Suite.NoError(err)
	p.Suite.NotNil(getProduct)
	p.Suite.Equal(getProduct.Id, filter["id"])
	p.Suite.Equal(getProduct.Description, "test")
	p.Suite.Equal(getProduct.Name, "test")
	p.Suite.Equal(getProduct.Category, "test")
	p.Suite.Equal(getProduct.MadeIn, "test")
	p.Suite.Equal(getProduct.Color, "test")
	p.Suite.Equal(getProduct.Count, int64(1))
	p.Suite.Equal(getProduct.Discount, float32(1))
	p.Suite.Equal(getProduct.AgeMin, int64(1))
	p.Suite.Equal(getProduct.AgeMax, int64(1))
	p.Suite.Equal(getProduct.TemperatureMin, int64(1))
	p.Suite.Equal(getProduct.TemperatureMax, int64(1))
	p.Suite.Equal(getProduct.ForGender, "test")
	p.Suite.Equal(getProduct.Size, int64(1))
	p.Suite.NotNil(getProduct.CreatedAt)
	p.Suite.NotNil(getProduct.UpdatedAt)
}

func (p *ProductRepositorySuiteTest) TestGetProducts() {
	listRequest := &entity.ListRequest{
		Page:  1,
		Limit: 10,
	}
	getProducts, err := p.Repository.GetProducts(context.Background(), listRequest)
	p.Suite.NoError(err)
	p.Suite.NotNil(getProducts)
	p.Suite.LessOrEqual(len(getProducts), 10)
}

func (p *ProductRepositorySuiteTest) TestUpdateProduct() {
	productReq := &entity.Product{
		Id:             "19e41c4a-5a82-47fa-90be-ea3f2860b59a",
		Name:           "update test",
		Description:    "update test",
		Category:       "update test",
		MadeIn:         "update test",
		Color:          "update test",
		Cost:           1,
		Count:          1,
		Discount:       1,
		AgeMin:         1,
		AgeMax:         1,
		TemperatureMin: 1,
		TemperatureMax: 1,
		ForGender:      "uptest",
		Size:           1,
		UpdatedAt:      time.Now(),
	}

	err := p.Repository.UpdateProduct(context.Background(), productReq)
	p.Suite.NoError(err)
}

func (p *ProductRepositorySuiteTest) TestDeleteProduct(){
	idReq := "19e41c4a-5a82-47fa-90be-ea3f2860b59a"
	err := p.Repository.DeleteProduct(context.Background(), idReq)
	p.Suite.NoError(err)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ProductRepositorySuiteTest))
	suite.Run(t, new(ProductRepositorySuiteTest))
}
