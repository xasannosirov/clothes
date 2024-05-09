package repository

import (
	"context"
	"product-service/internal/entity"
)

type Product interface {
	CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	GetProduct(ctx context.Context, params map[string]string) (*entity.Product, error)
	GetProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListProduct, error)
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, ID string) error

	CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	CancelOrder(ctx context.Context, ID string) error
	GetOrderByID(ctx context.Context, params map[string]string) (*entity.Order, error)
	GetAllOrders(ctx context.Context, req *entity.ListRequest) (*entity.ListOrders, error)

	GetDiscountProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListProduct, error)
	SearchProduct(ctx context.Context, req *entity.Filter) (*entity.ListProduct, error)
	RecommentProducts(ctx context.Context, req *entity.Recom) (*entity.ListProduct, error)

	IsUnique(ctx context.Context, tableName, userId, ProductId string) (bool, error)
	LikeProduct(ctx context.Context, req *entity.LikeProduct) (bool, error)
	DeleteLikeProduct(ctx context.Context, userId, productId string) error

	SaveProduct(ctx context.Context, req *entity.SaveProduct) (bool, error)
	DeleteSaveProduct(ctx context.Context, userId, productId string) error

	CommentToProduct(ctx context.Context, req *entity.CommentToProduct) (bool, error)
	GetAllComments(ctx context.Context, req *entity.ListRequest) (*entity.ListComments, error)

	GetProductOrders(ctx context.Context, req *entity.GetWithID) (*entity.ListOrders, error)
	GetProductComments(ctx context.Context, req *entity.GetWithID) (*entity.ListComments, error)
	GetProductLikes(ctx context.Context, req *entity.GetWithID) (*entity.ListLikes, error)
	GetProductStars(ctx context.Context, req *entity.GetWithID) (*entity.ListStars, error)

	GetSavedProductsByUserID(ctx context.Context, req string) (*entity.ListProduct, error)
	GetWishlistByUserID(ctx context.Context, req string) (*entity.ListProduct, error)
	GetOrderedProductsByUserID(ctx context.Context, req string) (*entity.ListProduct, error)

	StarProduct(ctx context.Context, req *entity.StarProduct) (*entity.StarProduct, error)
	GetAllStars(ctx context.Context, req *entity.ListRequest) (*entity.ListStars, error)
	GetDisableProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListOrders, error)

	CreateCategory(ctx context.Context, req *entity.Category) (*entity.Category, error)
	UpdateCategory(ctx context.Context, req *entity.Category) (*entity.Category, error)
	DeleteCategory(ctx context.Context, id string) error
	GetCategory(ctx context.Context, id string) (*entity.Category, error)
	ListCategory(ctx context.Context, req *entity.ListRequest) (*entity.LiestCategory, error)
}
