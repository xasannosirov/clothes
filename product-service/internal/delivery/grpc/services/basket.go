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

func (d *productRPC) GetBasket(ctx context.Context, in *pb.BasketGetReq) (*pb.Basket, error) {
	baskets, err := d.productUsecase.GetBasket(ctx, &entity.GetBAsketReq{
		UserId: in.UserId,
		Page:   in.Page,
		Limit:  in.Limit,
	})
	if err != nil {
		return nil, err
	}
	resBasket := pb.Basket{
		UserId:     baskets.UserID,
		TotalCount: baskets.TotalCount,
	}

	for _, productId := range baskets.ProductIDs {
		product, err := d.productUsecase.GetProduct(ctx, map[string]string{"id": productId})
		if err != nil {
			return nil, err
		}

		resBasket.Product = append(resBasket.Product, &pb.Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Category:    product.Category,
			MadeIn:      product.MadeIn,
			Color:       product.Color,
			Count:       product.Count,
			Cost:        product.Cost,
			Discount:    product.Discount,
			AgeMin:      product.AgeMin,
			AgeMax:      product.AgeMax,
			ForGender:   product.ForGender,
			ProductSize: product.Size,
		})

	}

	return &resBasket, nil
}
