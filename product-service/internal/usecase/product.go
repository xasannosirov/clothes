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
	GetAllProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListProduct, error)
	UpdateProduct(ctx context.Context, req *entity.Product) error
	DeleteProduct(ctx context.Context, id string) error

	CreateOrder(ctx context.Context, req *entity.Order) (*entity.Order, error)
	CancelOrder(ctx context.Context, id string) error
	GetOrderByID(ctx context.Context, params map[string]string) (*entity.Order, error)
	GetAllOrders(ctx context.Context, req *entity.ListRequest) (*entity.ListOrders, error)

	GetDiscountProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListProduct, error)
	SearchProduct(ctx context.Context, req *entity.Filter) (*entity.ListProduct, error)
	RecommentProducts(ctx context.Context, req *entity.Recom) (*entity.ListProduct, error)

	LikeProduct(ctx context.Context, req *entity.LikeProduct) (bool, error)
	SaveProduct(ctx context.Context, req *entity.SaveProduct) (bool, error)
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
	GetCategory(ctx context.Context, req *entity.GetWithID) (*entity.Category, error)
	DeleteCategory(ctx context.Context, id string) error
	ListCategory(ctx context.Context, req *entity.ListRequest) (*entity.LiestCategory, error)
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

func (u *productService) GetAllProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListProduct, error) {
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

func (u *productService) GetAllOrders(ctx context.Context, req *entity.ListRequest) (*entity.ListOrders, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetAllOrders(ctx, req)
}

func (u *productService) GetOrderByID(ctx context.Context, params map[string]string) (*entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetOrderByID(ctx, params)
}

func (u *productService) GetDiscountProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetDiscountProducts(ctx, req)
}

func (u *productService) SearchProduct(ctx context.Context, req *entity.Filter) (*entity.ListProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.SearchProduct(ctx, req)
}

func (u *productService) RecommentProducts(ctx context.Context, req *entity.Recom) (*entity.ListProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.RecommentProducts(ctx, req)
}

func (u *productService) LikeProduct(ctx context.Context, req *entity.LikeProduct) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	respStatus, err := u.repo.IsUnique(ctx, "wishlist", req.UserID, req.ProductID)

	if err != nil {
		log.Println("error while is check is unique", err)
		return false, err
	} else if respStatus {
		err := u.repo.DeleteLikeProduct(ctx, req.UserID, req.ProductID)
		if err != nil {
			return false, err
		}
		return false, nil
	} else {
		u.beforeRequest(&req.Id, &req.CreatedAt, &req.UpdatedAt)

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

	respStatus, err := u.repo.IsUnique(ctx, "saves", req.UserID, req.ProductID)

	if err != nil {
		log.Println("error while is check is unique", err)
		return false, err
	} else if respStatus {
		err := u.repo.DeleteSaveProduct(ctx, req.UserID, req.ProductID)
		if err != nil {
			return false, err
		}
		return false, nil
	} else {
		u.beforeRequest(&req.Id, &req.CreatedAt, &req.UpdatedAt)

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

	u.beforeRequest(&req.Id, &req.CreatedAt, &req.UpdatedAt)

	status, err := u.repo.CommentToProduct(ctx, req)
	if err != nil {
		return false, err
	}
	return status, nil
}

func (u *productService) GetProductComments(ctx context.Context, req *entity.GetWithID) (*entity.ListComments, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetProductComments(ctx, req)
}

func (u *productService) GetProductLikes(ctx context.Context, req *entity.GetWithID) (*entity.ListLikes, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetProductLikes(ctx, req)
}

func (u *productService) GetProductOrders(ctx context.Context, req *entity.GetWithID) (*entity.ListOrders, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetProductOrders(ctx, req)
}

func (u *productService) GetProductStars(ctx context.Context, req *entity.GetWithID) (*entity.ListStars, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetProductStars(ctx, req)
}

func (u *productService) GetSavedProductsByUserID(ctx context.Context, req string) (*entity.ListProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetSavedProductsByUserID(ctx, req)
}

func (u *productService) GetWishlistByUserID(ctx context.Context, req string) (*entity.ListProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetWishlistByUserID(ctx, req)
}

func (u *productService) GetOrderedProductsByUserID(ctx context.Context, req string) (*entity.ListProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetOrderedProductsByUserID(ctx, req)
}

func (u *productService) GetAllComments(ctx context.Context, req *entity.ListRequest) (*entity.ListComments, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetAllComments(ctx, req)
}

func (u *productService) GetAllStars(ctx context.Context, req *entity.ListRequest) (*entity.ListStars, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetAllStars(ctx, req)
}

func (u *productService) StarProduct(ctx context.Context, req *entity.StarProduct) (*entity.StarProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.StarProduct(ctx, req)
}

func (u *productService) GetDisableProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListOrders, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetDisableProducts(ctx, req)
}

func (u *productService) CreateCategory(ctx context.Context, req *entity.Category) (*entity.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.CreateCategory(ctx, req)
}

func (u *productService) DeleteCategory(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.DeleteCategory(ctx, id)
}

func (u *productService) ListCategory(ctx context.Context, req *entity.ListRequest) (*entity.LiestCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.ListCategory(ctx, req)
}

func (u *productService) GetCategory(ctx context.Context, req *entity.GetWithID) (*entity.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetCategory(ctx, req.ID)
}

func (u *productService) UpdateCategory(ctx context.Context, req *entity.Category) (*entity.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.UpdateCategory(ctx, req)
}
