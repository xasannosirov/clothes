package usecase

import (
	"context"
	"time"

	"product-service/internal/entity"
	"product-service/internal/infrastructure/repository"
)

type Product interface {
	CreateCategory(ctx context.Context, req *entity.Category) (*entity.Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	UpdateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error)
	GetCategory(ctx context.Context, req *entity.GetWithID) (*entity.Category, error)
	ListCategories(ctx context.Context, req *entity.ListRequest) (*entity.LiestCategory, error)
	SearchCategory(ctx context.Context, req *entity.SearchRequest) (*entity.ListProduct, error)
	UniqueCategory(ctx context.Context, req *entity.Params) (*entity.MoveResponse, error)

	CreateProduct(ctx context.Context, req *entity.Product) (*entity.Product, error)
	UpdateProduct(ctx context.Context, req *entity.Product) error
	DeleteProduct(ctx context.Context, id string) error
	GetProduct(ctx context.Context, params map[string]string) (*entity.Product, error)
	ListProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListProduct, error)
	SearchProduct(ctx context.Context, req *entity.SearchRequest) (*entity.ListProduct, error)
	GetDiscountProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListProduct, error)
	GetDisableProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListOrders, error)

	LikeProduct(ctx context.Context, req *entity.Like) (bool, error)
	UserWishlist(ctx context.Context, req *entity.SearchRequest) (*entity.ListProduct, error)
	IsUnique(ctx context.Context, tableName, userID, productID string) (bool, error)

	SaveToBasket(ctx context.Context, req *entity.BasketCreateReq) (*entity.MoveResponse, error)
	GetUserBaskets(ctx context.Context, req *entity.GetWithID) (*entity.ListProduct, error)

	// order
	CreateOrder(ctx context.Context, req *entity.Order) (*entity.Order, error)
	GetOrder(ctx context.Context, params map[string]string) (*entity.Order, error)
	DeleteOrder(ctx context.Context, params map[string]string) error
	UserOrderHistory(ctx context.Context, req *entity.SearchRequest) (*entity.ListProduct, error)

	// comment
	CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
	UpdateComment(ctx context.Context, category *entity.CommentUpdateRequest) (*entity.Comment, error)
	DeleteComment(ctx context.Context, req *entity.DeleteRequest) error
	GetComment(ctx context.Context, req *entity.GetRequest) (*entity.Comment, error)
	ListComment(ctx context.Context, req *entity.ListRequest) (*entity.CommentListResponse, error)
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

func (u *productService) CreateCategory(ctx context.Context, req *entity.Category) (*entity.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(&req.ID, nil, nil)

	return u.repo.CreateCategory(ctx, req)
}

func (u *productService) UpdateCategory(ctx context.Context, req *entity.Category) (*entity.Category, error) {
	return u.repo.UpdateCategory(ctx, req)
}

func (u *productService) DeleteCategory(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.DeleteCategory(ctx, id)
}

func (u *productService) GetCategory(ctx context.Context, req *entity.GetWithID) (*entity.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetCategory(ctx, req.ID)
}

func (u *productService) ListCategories(ctx context.Context, req *entity.ListRequest) (*entity.LiestCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.ListCategories(ctx, req)
}

func (u *productService) SearchCategory(ctx context.Context, req *entity.SearchRequest) (*entity.ListProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.SearchCategory(ctx, req)
}

func (u *productService) UniqueCategory(ctx context.Context, req *entity.Params) (*entity.MoveResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.UniqueCategory(ctx, req)
}

func (u *productService) CreateProduct(ctx context.Context, req *entity.Product) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(&req.Id, nil, nil)

	return u.repo.CreateProduct(ctx, req)
}

func (u *productService) UpdateProduct(ctx context.Context, req *entity.Product) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(nil, nil, nil)

	return u.repo.UpdateProduct(ctx, req)
}

func (u *productService) DeleteProduct(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.DeleteProduct(ctx, id)
}

func (u *productService) GetProduct(ctx context.Context, params map[string]string) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetProduct(ctx, params)
}

func (u *productService) ListProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.ListProducts(ctx, req)
}

func (u *productService) SearchProduct(ctx context.Context, req *entity.SearchRequest) (*entity.ListProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.SearchProduct(ctx, req)
}

func (u *productService) GetDisableProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListOrders, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetDisableProducts(ctx, req)
}

func (u *productService) GetDiscountProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetDiscountProducts(ctx, req)
}

func (u productService) SaveToBasket(ctx context.Context, req *entity.BasketCreateReq) (*entity.MoveResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.SaveToBasket(ctx, req)
}

func (u productService) GetUserBaskets(ctx context.Context, req *entity.GetWithID) (*entity.ListProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetUserBaskets(ctx, req)
}

func (u *productService) LikeProduct(ctx context.Context, req *entity.Like) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	respStatus, err := u.repo.IsUnique(ctx, "wishlist", req.UserID, req.ProductID)

	if err != nil {
		return false, err
	} else if respStatus {
		err := u.repo.DeleteLikeProduct(ctx, req.UserID, req.ProductID)
		if err != nil {
			return false, err
		}

		return false, nil
	} else {
		u.beforeRequest(&req.Id, nil, nil)
		resp, err := u.repo.LikeProduct(ctx, req)
		if err != nil {
			return false, err
		}

		return resp, nil
	}
}

func (u *productService) UserWishlist(ctx context.Context, req *entity.SearchRequest) (*entity.ListProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.UserWishlist(ctx, req)
}

func (u *productService) CreateOrder(ctx context.Context, req *entity.Order) (*entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(&req.Id, nil, nil)

	return u.repo.CreateOrder(ctx, req)
}

func (u *productService) GetOrder(ctx context.Context, req map[string]string) (*entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetOrder(ctx, req)
}

func (u *productService) DeleteOrder(ctx context.Context, params map[string]string) error {
	return nil
}

func (u *productService) UserOrderHistory(ctx context.Context, req *entity.SearchRequest) (*entity.ListProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.UserWishlist(ctx, req)
}

func (u *productService) IsUnique(ctx context.Context, tableName, userID, productID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.IsUnique(ctx, tableName, userID, productID)
}

func (p *productService) CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {

	p.beforeRequest(&comment.Id, &comment.CreatedAt, &comment.UpdatedAt)

	return p.repo.CreateComment(ctx, comment)
}

func (p *productService) UpdateComment(ctx context.Context, comment *entity.CommentUpdateRequest) (*entity.Comment, error) {

	p.beforeRequest(nil, nil, &comment.UpdatedAt)

	return p.repo.UpdateComment(ctx, comment)
}

func (p *productService) DeleteComment(ctx context.Context, req *entity.DeleteRequest) error {

	req.Deleted_at = time.Now()
	return p.repo.DeleteComment(ctx, req)
}

func (p *productService) GetComment(ctx context.Context, req *entity.GetRequest) (*entity.Comment, error) {

	return p.repo.GetComment(ctx, req)
}

func (p *productService) ListComment(ctx context.Context, req *entity.ListRequest) (*entity.CommentListResponse, error) {

	return p.repo.ListComment(ctx, req)
}
