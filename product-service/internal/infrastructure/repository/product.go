package repository

import (
	"context"
	"product-service/internal/entity"
)

type Product interface {
	CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	GetProduct(ctx context.Context, params map[string]string) (*entity.Product, error)
	GetProducts(ctx context.Context, req *entity.ListRequest) ([]*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, ID string) error

	CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	CancelOrder(ctx context.Context, ID string) error
	GetOrderByID(ctx context.Context, params map[string]string) (*entity.Order, error)
	GetAllOrders(ctx context.Context, req *entity.ListRequest) ([]*entity.Order, error)

	GetDiscountProducts(ctx context.Context, req *entity.ListRequest) ([]*entity.Product, error)
	SearchProduct(ctx context.Context, req *entity.Filter) ([]*entity.Product, error)
	RecommentProducts(ctx context.Context, req *entity.Recom) ([]*entity.Product, error)
}
