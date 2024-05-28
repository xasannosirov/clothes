package services

import (
	"context"
	pb "product-service/genproto/product_service"
	"product-service/internal/entity"
)

func (d *productRPC) SaveToBasket(ctx context.Context, in *pb.BasketCreateReq) (*pb.GetWithID, error) {
	respProduct, err := d.productUsecase.SaveToBasket(ctx, &entity.BasketCreateReq{
		UserID:    in.UserId,
		ProductID: in.ProductId,
	})

	if err != nil {
		return nil, err
	}

	return &pb.GetWithID{Id: respProduct.UserID}, nil
}

func (d *productRPC) DeleteFromBasket(ctx context.Context, in *pb.DeleteBasket) (*pb.MoveResponse, error) {
	err := d.productUsecase.DeleteFromBasket(ctx, in.UserId, in.ProductId)

	if err != nil {
		return nil, err
	}

	return &pb.MoveResponse{Status: true}, nil
}



func (d *productRPC) GetBasket(ctx context.Context, in *pb.GetWithID) (*pb.Basket, error) {
	baskets, err := d.productUsecase.GetBasket(ctx, &entity.GetWithID{
		ID: in.Id,
	})
	if err != nil {
		return nil, err
	}

	

	return &pb.Basket{
		UserId: baskets.UserID,
		ProductId: baskets.ProductIDs,
		Count: baskets.Count,
	}, nil
}
