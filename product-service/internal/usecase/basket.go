package usecase

import (
	"context"
	"product-service/internal/entity"
	"product-service/internal/pkg/otlp"
)

func (u productService) CreateBasket(ctx context.Context, req *entity.Basket) (*entity.Basket, error) {
	ctx, span := otlp.Start(ctx, "basket_grpc-usecase", "CreateProduct")
	defer span.End()

	u.beforeRequest(&req.Id, &req.Created_at, &req.Updated_at)

	return u.repo.CreateBasket(ctx, req)
}
func (u productService) GetBasket(ctx context.Context, req map[string]string) (*entity.Basket, error) {
	ctx, span := otlp.Start(ctx, "basket-grpc-usecase", "GetBasket")
	defer span.End()

	return u.repo.GetBasket(ctx, req)
}
func (u productService) GetBaskets(ctx context.Context, req *entity.ListBasketReq) (*entity.ListBasketRes, error) {
	ctx, span := otlp.Start(ctx, "basket-grpc-usecase", "GetBaskets")
	defer span.End()

	return u.repo.GetBaskets(ctx, req)
}
func (u productService) DeleteBasket(ctx context.Context, id string) error {
	ctx, span := otlp.Start(ctx, "basket-grpc-usecase", "DeleteBasket")
	defer span.End()
	return u.repo.DeleteBasket(ctx, id)
}
