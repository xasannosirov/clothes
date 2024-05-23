package services

import (
	"context"
	pb "product-service/genproto/product_service"
	"product-service/internal/entity"
)

func (d *productRPC) LikeProduct(ctx context.Context, req *pb.Like) (*pb.MoveResponse, error) {
	status, err := d.productUsecase.LikeProduct(ctx, &entity.Like{
		ProductID: req.ProductId,
		UserID:    req.UserId,
	})

	if err != nil {
		return nil, err
	}

	return &pb.MoveResponse{
		Status: status,
	}, nil
}

func (p *productRPC) UserWishlist(ctx context.Context, in *pb.SearchRequest) (*pb.ListProduct, error) {
	products, err := p.productUsecase.UserOrderHistory(ctx, &entity.SearchRequest{
		Page:   in.Page,
		Limit:  in.Limit,
		Params: in.Params,
	})

	if err != nil {
		return nil, err
	}

	response := pb.ListProduct{}
	for _, product := range products.Products {
		response.Products = append(response.Products, &pb.Product{
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
	response.TotalCount = products.TotalCount

	return &response, nil
}
