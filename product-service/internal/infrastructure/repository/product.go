package repository

import (
	"context"
	"product-service/internal/entity"
)

type Product interface {
	CreateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	UpdateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error)
	GetCategory(ctx context.Context, categoryID string) (*entity.Category, error)
	ListCategories(ctx context.Context, listReq *entity.ListRequest) (*entity.LiestCategory, error)
	SearchCategory(ctx context.Context, searchFields *entity.SearchRequest) (*entity.ListProduct, error)
	UniqueCategory(ctx context.Context, params *entity.Params) (*entity.MoveResponse, error)

	CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, productID string) error
	GetProduct(ctx context.Context, params map[string]string) (*entity.Product, error)
	ListProducts(ctx context.Context, listReq *entity.ListRequest) (*entity.ListProduct, error)
	SearchProduct(ctx context.Context, searchFields *entity.SearchRequest) (*entity.ListProduct, error)
	GetDisableProducts(ctx context.Context, listReq *entity.ListRequest) (*entity.ListOrders, error)
	GetDiscountProducts(ctx context.Context, listReq *entity.ListRequest) (*entity.ListProduct, error)

	LikeProduct(ctx context.Context, like *entity.Like) (bool, error)
	DeleteLikeProduct(ctx context.Context, userID, productID string) error
	IsUnique(ctx context.Context, tableName, userID, productID string) (bool, error)
	UserWishlist(ctx context.Context, searchFields *entity.SearchRequest) (*entity.ListProduct, error)

	SaveToBasket(ctx context.Context, basket *entity.Basket) (*entity.Basket, error)
	DeleteFromBasket(ctx context.Context, basketID string) error
	UpdateBasket(ctx context.Context, basket *entity.Basket) (*entity.Basket, error)
	ListBaskets(ctx context.Context, getReq *entity.GetWithID) (*entity.ListBasket, error)

	CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	GetOrder(ctx context.Context, params map[string]string) (*entity.Order, error)
	DeleteOrder(ctx context.Context, params map[string]string) error
	UserOrderHistory(ctx context.Context, searchFields *entity.SearchRequest) (*entity.ListProduct, error)
}
