package postgresql

import (
	"context"
	"log"
	"product-service/internal/entity"
	"product-service/internal/infrastructure/repository"
	"testing"
	"time"

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

var (
	product = &entity.Product{
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
		Size:           1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	order = &entity.Order{
		Id:        uuid.NewString(),
		UserID:    uuid.NewString(),
		Status:    "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
)

func (p *ProductRepositorySuiteTest) TestCreateProduct() {

	createProduct, err := p.Repository.CreateProduct(context.Background(), product)

	p.Suite.NoError(err)
	p.Suite.NotNil(createProduct)
	p.Suite.Equal(product.Id, createProduct.Id)
	p.Suite.Equal(product.Description, createProduct.Description)
	p.Suite.Equal(product.Name, createProduct.Name)
	p.Suite.Equal(product.Category, createProduct.Category)
	p.Suite.Equal(product.MadeIn, createProduct.MadeIn)
	p.Suite.Equal(product.Color, createProduct.Color)
	p.Suite.Equal(product.Cost, createProduct.Cost)
	p.Suite.Equal(product.Count, createProduct.Count)
	p.Suite.Equal(product.Discount, createProduct.Discount)
	p.Suite.Equal(product.AgeMin, createProduct.AgeMin)
	p.Suite.Equal(product.AgeMax, createProduct.AgeMax)
	p.Suite.Equal(product.TemperatureMin, createProduct.TemperatureMin)
	p.Suite.Equal(product.TemperatureMax, createProduct.TemperatureMax)
	p.Suite.Equal(product.ForGender, createProduct.ForGender)
	p.Suite.Equal(product.Size, createProduct.Size)
	p.Suite.NotNil(createProduct.CreatedAt)
	p.Suite.NotNil(createProduct.UpdatedAt)
}

func (p *ProductRepositorySuiteTest) TestGetProduct() {

	createProduct, err := p.Repository.CreateProduct(context.Background(), product)
	p.Suite.NoError(err)

	filter := make(map[string]string)
	filter["id"] = createProduct.Id
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

	createProduct, err := p.Repository.CreateProduct(context.Background(), product)
	p.Suite.NoError(err)

	productReq := &entity.Product{
		Id:             createProduct.Id,
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

	err = p.Repository.UpdateProduct(context.Background(), productReq)
	p.Suite.NoError(err)
}

func (p *ProductRepositorySuiteTest) TestDeleteProduct() {

	createProduct, err := p.Repository.CreateProduct(context.Background(), product)
	p.Suite.NoError(err)

	idReq := createProduct.Id
	err = p.Repository.DeleteProduct(context.Background(), idReq)
	p.Suite.NoError(err)
}

func (p *ProductRepositorySuiteTest) TestCreateOrder() {

	productForID, err := p.Repository.CreateProduct(context.Background(), product)
	p.Suite.NoError(err)

	order.ProductID = productForID.Id

	createOrder, err := p.Repository.CreateOrder(context.Background(), order)
	p.Suite.NoError(err)
	p.Suite.NotNil(createOrder)
	p.Suite.Equal(order.Id, createOrder.Id)
	p.Suite.Equal(order.ProductID, createOrder.ProductID)
	p.Suite.Equal(order.UserID, createOrder.UserID)
	p.Suite.Equal(order.Status, createOrder.Status)
	p.Suite.NotNil(createOrder.CreatedAt)
	p.Suite.NotNil(createOrder.UpdatedAt)
}

func (p *ProductRepositorySuiteTest) TestGetOrderByID() {

	productForID, err := p.Repository.CreateProduct(context.Background(), product)
	p.Suite.NoError(err)
	order.ProductID = productForID.Id
	createOrder, err := p.Repository.CreateOrder(context.Background(), order)
	p.Suite.NoError(err)

	param := make(map[string]string)
	param["id"] = createOrder.Id
	getOrder, err := p.Repository.GetOrderByID(context.Background(), param)
	p.Suite.NoError(err)
	p.Suite.NotNil(getOrder)
	p.Suite.Equal(getOrder.Id, param["id"])
	p.Suite.NotNil(getOrder.ProductID)
	p.Suite.NotNil(getOrder.UserID)
	p.Suite.Equal(getOrder.Status, "test")
	p.Suite.IsType(time.Now(), getOrder.CreatedAt)
	p.Suite.IsType(time.Now(), getOrder.UpdatedAt)
}

func (p *ProductRepositorySuiteTest) TestGetAllOrders() {
	filter := &entity.ListRequest{
		Page:  1,
		Limit: 10,
	}
	getOrders, err := p.Repository.GetAllOrders(context.Background(), filter)
	p.Suite.NoError(err)
	p.Suite.LessOrEqual(len(getOrders), 10)
	p.Suite.IsType([]*entity.Order{}, getOrders)
}

func (p *ProductRepositorySuiteTest) TestCancelOrder() {

	productForID, err := p.Repository.CreateProduct(context.Background(), product)
	p.Suite.NoError(err)
	order.ProductID = productForID.Id
	createOrder, err := p.Repository.CreateOrder(context.Background(), order)
	p.Suite.NoError(err)
	
	id := createOrder.Id
	err = p.Repository.CancelOrder(context.Background(), id)
	p.Suite.NoError(err)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ProductRepositorySuiteTest))
}
