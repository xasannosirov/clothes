package usecase

import (
	"context"
	"log"
	"time"

	"product-service/internal/entity"
	"product-service/internal/infrastructure/repository"
	"product-service/internal/pkg/otlp"
)

type Product interface {
	CreateProduct(ctx context.Context, req *entity.Product) (*entity.Product, error)
	GetProductByID(ctx context.Context, params map[string]string) (*entity.Product, error)
	GetAllProducts(ctx context.Context, req *entity.ListRequest) ([]*entity.Product, error)
	UpdateProduct(ctx context.Context, req *entity.Product) error
	DeleteProduct(ctx context.Context, id string) error

	CreateOrder(ctx context.Context, req *entity.Order) (*entity.Order, error)
	CancelOrder(ctx context.Context, id string) error
	GetOrderByID(ctx context.Context, params map[string]string) (*entity.Order, error)
	GetAllOrders(ctx context.Context, req *entity.ListRequest) ([]*entity.Order, error)

	GetDiscountProducts(ctx context.Context, req *entity.ListRequest) ([]*entity.Product, error)
	SearchProduct(ctx context.Context, req *entity.Filter) ([]*entity.Product, error)
	RecommentProducts(ctx context.Context, req *entity.Recom) ([]*entity.Product, error)

	LikeProduct(ctx context.Context, req *entity.LikeProduct) (bool, error)
	SaveProduct(ctx context.Context, req *entity.SaveProduct) (bool, error)

	CommentToProduct(ctx context.Context, req *entity.CommentToProduct) (bool, error)
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

type productService struct {
	BaseUseCase
	repo       repository.Product
	ctxTimeout time.Duration
}

func NewProductService(ctxTimeout time.Duration, repo repository.Product) Product {
	return &productService{
		BaseUseCase: BaseUseCase{},
		repo:        repo,
		ctxTimeout:  ctxTimeout,
	}
}

func (u *productService) CreateProduct(ctx context.Context, req *entity.Product) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, "product_grpc-usercase", "CreateProduct")
	defer span.End()

	u.beforeRequest(&req.Id, &req.CreatedAt, &req.UpdatedAt)

	return u.repo.CreateProduct(ctx, req)
}

func (u *productService) GetProductByID(ctx context.Context, params map[string]string) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetProduct(ctx, params)
}

func (u *productService) GetAllProducts(ctx context.Context, req *entity.ListRequest) ([]*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetProducts(ctx, req)
}

func (u *productService) UpdateProduct(ctx context.Context, req *entity.Product) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(nil, nil, &req.UpdatedAt)

	return u.repo.UpdateProduct(ctx, req)
}

func (u *productService) DeleteProduct(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.DeleteProduct(ctx, id)
}

func (u *productService) CancelOrder(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.CancelOrder(ctx, id)
}

func (u *productService) CreateOrder(ctx context.Context, req *entity.Order) (*entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(&req.Id, &req.CreatedAt, &req.UpdatedAt)

	return u.repo.CreateOrder(ctx, req)
}

func (u *productService) GetAllOrders(ctx context.Context, req *entity.ListRequest) ([]*entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetAllOrders(ctx, req)
}

func (u *productService) GetOrderByID(ctx context.Context, params map[string]string) (*entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetOrderByID(ctx, params)
}

func (u *productService) GetDiscountProducts(ctx context.Context, req *entity.ListRequest) ([]*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetDiscountProducts(ctx, req)
}

func (u *productService) SearchProduct(ctx context.Context, req *entity.Filter) ([]*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.SearchProduct(ctx, req)
}

func (u *productService) RecommentProducts(ctx context.Context, req *entity.Recom) ([]*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.RecommentProducts(ctx, req)
}

func (u *productService) LikeProduct(ctx context.Context, req *entity.LikeProduct) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	respStatus, err := u.repo.IsUnique(ctx, "wishlist", req.User_id, req.Product_id)

	if err != nil {
		log.Println("error while is check is unique", err)
		return false, err
	} else if respStatus {
		err := u.repo.DeleteLikeProduct(ctx, req.User_id, req.Product_id)
		if err != nil {
			return false, err
		}
		return false, nil
	} else {
		u.beforeRequest(&req.Id, &req.Created_at, &req.Updated_at)

		resp, err := u.repo.LikeProduct(ctx, req)
		if err != nil {
			return false, err
		}
		return resp, nil
	}
}

func (u *productService) SaveProduct(ctx context.Context, req *entity.SaveProduct) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	respStatus, err := u.repo.IsUnique(ctx, "saves", req.User_id, req.Product_id)

	if err != nil {
		log.Println("error while is check is unique", err)
		return false, err
	} else if respStatus {
		err := u.repo.DeleteSaveProduct(ctx, req.User_id, req.Product_id)
		if err != nil {
			return false, err
		}
		return false, nil
	} else {
		u.beforeRequest(&req.Id, &req.Created_at, &req.Updated_at)

		resp, err := u.repo.SaveProduct(ctx, req)
		if err != nil {
			return false, err
		}
		return resp, nil
	}
}

func (u *productService) CommentToProduct(ctx context.Context, req *entity.CommentToProduct) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(&req.Id, &req.Created_at, &req.Updated_at)

	status, err := u.repo.CommentToProduct(ctx, req)
	if err != nil {
		return false, err
	}
	return status, nil
}

// GetProductComments implements Product.
func (u *productService) GetProductComments(ctx context.Context, req *entity.GetWithID) ([]*entity.CommentToProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetProductComments(ctx, req)
}

// GetProductLikes implements Product.
func (u *productService) GetProductLikes(ctx context.Context, req *entity.GetWithID) ([]*entity.LikeProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetProductLikes(ctx, req)
}

// GetProductOrders implements Product.
func (u *productService) GetProductOrders(ctx context.Context, req *entity.GetWithID) ([]*entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetProductOrders(ctx, req)
}

// GetProductStars implements Product.
func (u *productService) GetProductStars(ctx context.Context, req *entity.GetWithID) ([]*entity.StarProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetProductStars(ctx, req)
}

// GetSavedProductsByUserID implements Product.
func (u *productService) GetSavedProductsByUserID(ctx context.Context, req string) ([]*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetSavedProductsByUserID(ctx, req)
}

func (u *productService) GetWishlistByUserID(ctx context.Context, req string) ([]*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetWishlistByUserID(ctx, req)
}

// GetOrderedProductsByUserID implements Product.
func (u *productService) GetOrderedProductsByUserID(ctx context.Context, req string) ([]*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetOrderedProductsByUserID(ctx, req)
}

// GetAllComments implements Product.
func (u *productService) GetAllComments(ctx context.Context, req *entity.ListRequest) ([]*entity.CommentToProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetAllComments(ctx, req)
}

// GetAllStars implements Product.
func (u *productService) GetAllStars(ctx context.Context, req *entity.ListRequest) ([]*entity.StarProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetAllStars(ctx, req)
}

// StarProduct implements Product.
func (u *productService) StarProduct(ctx context.Context, req *entity.StarProduct) (*entity.StarProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.StarProduct(ctx, req)
}

// GetDisableProducts implements Product.
func (u *productService) GetDisableProducts(ctx context.Context, req *entity.ListRequest) ([]*entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetDisableProducts(ctx, req)
}
