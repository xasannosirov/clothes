package services

import (
	"context"
	pb "product-service/genproto/product_service"
	grpc "product-service/internal/delivery"
	"product-service/internal/entity"
	"product-service/internal/pkg/otlp"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func (d *productRPC) SaveToBasket(ctx context.Context, in *pb.Basket) (*pb.GetWithID, error) {
	ctx, span := otlp.Start(ctx, "product_grpc-delivery", "CreateBasket")
	span.SetAttributes(
		attribute.Key("guid").String(in.Id),
	)
	defer span.End()

	respProduct, err := d.productUsecase.CreateBasket(ctx, &entity.Basket{
		Id:         in.Id,
		UserId:     in.UserId,
		ProductId:  in.ProductId,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	})
	if err != nil {
		d.logger.Error("productUseCase.CreateBasket", zap.Error(err))
		return &pb.GetWithID{}, grpc.Error(ctx, err)
	}
	return &pb.GetWithID{Id: respProduct.Id}, nil
}

func (d *productRPC) GetBasketProduct(ctx context.Context, in *pb.RequestBasket) (*pb.Basket, error) {
	basket, err := d.productUsecase.GetBasket(ctx, in.Filter)
	if err != nil {
		d.logger.Error("productUseCase.GetProductByID", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	return &pb.Basket{
		Id:        basket.Id,
		UserId:    basket.UserId,
		ProductId: basket.ProductId,
		CreatedAt: basket.Created_at.String(),
		UpdatedAt: basket.Updated_at.String(),
	}, nil
}

func (d *productRPC) DeleteFromBasket(ctx context.Context, in *pb.RequestBasket) (*pb.DeleteResponse, error) {
	err := d.productUsecase.DeleteBasket(ctx, in.Filter["id"])
	if err != nil {
		d.logger.Error("productUseCase.DeleteProduct", zap.Error(err))
		return &pb.DeleteResponse{Status: false}, err
	}

	return &pb.DeleteResponse{Status: true}, nil
}

func (d *productRPC) GetBasketProducts(ctx context.Context, in *pb.ListBasketRequest) (*pb.ListBasketResponse, error) {
	filter := &entity.ListBasketReq{
		Limit: in.Limit,
		Page:  in.Page,
	}

	baskets, err := d.productUsecase.GetBaskets(ctx, filter)
	if err != nil {
		d.logger.Error("productUseCase.List", zap.Error(err))
		return nil, err
	}

	pbBaskets := &pb.ListBasketResponse{}
	for _, basket := range baskets.Basket {
		pbBaskets.Baskets = append(pbBaskets.Baskets, &pb.Basket{
			Id:        basket.Id,
			UserId:    basket.UserId,
			ProductId: basket.ProductId,
			CreatedAt: basket.Created_at.String(),
			UpdatedAt: basket.Updated_at.String(),
		})
	}

	return &pb.ListBasketResponse{
		Baskets:   pbBaskets.Baskets,
		TotalCount: int64(baskets.TotalCount),
	}, nil
}
