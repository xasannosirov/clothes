package services

import (
	"context"
	pb "product-service/genproto/product_service"
	"product-service/internal/entity"
)

func (d *productRPC) SaveToBasket(ctx context.Context, in *pb.BasketCreateReq) (*pb.MoveResponse, error) {
	response, err := d.productUsecase.SaveToBasket(ctx, &entity.BasketCreateReq{
		UserID:    in.UserId,
		ProductID: in.ProductId,
	})

	if err != nil {
		return nil, err
	}

	return &pb.MoveResponse{
		Status: response.Status,
	}, nil
}

func (d *productRPC) GetUserBaskets(ctx context.Context, in *pb.GetWithID) (*pb.ListBaskedProducts, error) {
	products, err := d.productUsecase.GetUserBaskets(ctx, &entity.GetWithID{
		ID: in.Id,
	})
	if err != nil {
		return nil, err
	}

	var basketProducts pb.ListBaskedProducts
	for _, product := range products.Products {

		basketProducts.Products = append(basketProducts.Products, &pb.Product{
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

	return &basketProducts, nil
}
