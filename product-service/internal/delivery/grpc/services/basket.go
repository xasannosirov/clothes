package services

import (
	"context"
	pb "product-service/genproto/product_service"
	"product-service/internal/entity"
)

func (d *productRPC) SaveToBasket(ctx context.Context, in *pb.Basket) (*pb.GetWithID, error) {
	respProduct, err := d.productUsecase.SaveToBasket(ctx, &entity.Basket{
		ID:        in.Id,
		UserID:    in.UserId,
		ProductID: in.ProductId,
		Count:     in.Count,
	})

	if err != nil {
		return nil, err
	}

	return &pb.GetWithID{Id: respProduct.ID}, nil
}

func (d *productRPC) DeleteFromBasket(ctx context.Context, in *pb.GetWithID) (*pb.MoveResponse, error) {
	err := d.productUsecase.DeleteFromBasket(ctx, in.Id)

	if err != nil {
		return nil, err
	}

	return &pb.MoveResponse{Status: true}, nil
}

func (p *productRPC) UpdateBasket(ctx context.Context, in *pb.Basket) (*pb.Basket, error) {
	basket, err := p.productUsecase.UpdateBasket(ctx, &entity.Basket{
		ID:        in.Id,
		ProductID: in.ProductId,
		UserID:    in.UserId,
		Count:     in.Count,
	})

	if err != nil {
		return nil, err
	}

	return &pb.Basket{
		Id:        basket.ID,
		ProductId: basket.ProductID,
		UserId:    basket.UserID,
		Count:     basket.Count,
	}, nil
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
			Id:        basket.ID,
			UserId:    basket.UserID,
			ProductId: basket.ProductID,
			Count:     basket.Count,
		})
	}

	return &pb.ListBasket{
		Baskets:    pbBaskets.Baskets,
		TotalCount: baskets.TotalCount,
	}, nil
}
