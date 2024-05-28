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



func (d *productRPC) ListBaskets(ctx context.Context, in *pb.GetWithID) (*pb.ListBasket, error) {
	baskets, err := d.productUsecase.ListBaskets(ctx, &entity.GetWithID{
		ID: in.Id,
	})
	if err != nil {
		return nil, err
	}

	pbBaskets := &pb.ListBasket{}
	for _, basket := range baskets.Baskets {
		pbBaskets.Baskets = append(pbBaskets.Baskets, &pb.Basket{
			UserId:    basket.UserID,
			ProductId: basket.ProductIDs,
			Count:     basket.Count,
		})
	}

	return &pb.ListBasket{
		Baskets:    pbBaskets.Baskets,
		TotalCount: baskets.TotalCount,
	}, nil
}
