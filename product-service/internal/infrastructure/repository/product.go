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

	IsUnique(ctx context.Context, tableName, userId, ProductId string)(bool, error)
	LikeProduct(ctx context.Context, req *entity.LikeProduct)(bool, error)
	DeleteLikeProduct(ctx context.Context, userId, productId string)(error)

	SaveProduct(ctx context.Context, req *entity.SaveProduct)(bool, error)
	DeleteSaveProduct(ctx context.Context, userId, productId string)error

	CommentToProduct(ctx context.Context, req *entity.CommentToProduct)(bool, error)
	GetAllComments(ctx context.Context, req *entity.ListRequest) ([]*entity.CommentToProduct, error)

	GetProductOrders(ctx context.Context, req *entity.GetWithID) ([]*entity.Order, error)
	GetProductComments(ctx context.Context, req *entity.GetWithID) ([]*entity.CommentToProduct, error)
	GetProductLikes(ctx context.Context, req *entity.GetWithID) ([]*entity.LikeProduct, error)
	GetProductStars(ctx context.Context, req *entity.GetWithID) ([]*entity.StarProduct, error)

	GetSavedProductsByUserID(ctx context.Context, req string) ([]*entity.Product, error)
	GetWishlistByUserID(ctx context.Context, req string) ([]*entity.Product, error)
	GetOrderedProductsByUserID(ctx context.Context, req string) ([]*entity.Product, error)

	StarProduct(ctx context.Context, req *entity.StarProduct) (*entity.StarProduct, error)
	GetAllStars(ctx context.Context, req *entity.ListRequest) ([]*entity.StarProduct, error)

	GetDisableProducts(ctx context.Context, req *entity.ListRequest) ([]*entity.Order, error)
}
