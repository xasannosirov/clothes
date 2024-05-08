package services

import (
	"context"
	pb "product-service/genproto/product_service"
	grpc "product-service/internal/delivery"

	"go.uber.org/zap"
)

func (d *productRPC) GetSavedProductsByUserID(ctx context.Context, in *pb.GetWithUserID) (*pb.ListProductResponse, error) {
	products, err := d.productUsecase.GetSavedProductsByUserID(ctx, in.UserId)
	if err != nil {
		d.logger.Error("productUsecase.GetSavedProductByUserID", zap.Error(err))
		return &pb.ListProductResponse{}, grpc.Error(ctx, err)
	}

	var pbProducts []*pb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, &pb.Product{
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
			Size_:       product.Size,
			CreatedAt:   product.CreatedAt.String(),
			UpdatedAt:   product.UpdatedAt.String(),
		})
	}
	return &pb.ListProductResponse{Products: pbProducts}, nil
}

func (d *productRPC) GetWishlistByUserID(ctx context.Context, in *pb.GetWithUserID) (*pb.ListProductResponse, error) {
	products, err := d.productUsecase.GetWishlistByUserID(ctx, in.UserId)
	if err != nil {
		d.logger.Error("productUsecase.GetWishlistByUserID", zap.Error(err))
		return &pb.ListProductResponse{}, grpc.Error(ctx, err)
	}

	var pbProducts []*pb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, &pb.Product{
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
			Size_:       product.Size,
			CreatedAt:   product.CreatedAt.String(),
			UpdatedAt:   product.UpdatedAt.String(),
		})
	}
	return &pb.ListProductResponse{Products: pbProducts}, nil
}

func (d *productRPC) GetOrderedProductsByUserID(ctx context.Context, in *pb.GetWithUserID) (*pb.ListProductResponse, error) {
	products, err := d.productUsecase.GetOrderedProductsByUserID(ctx, in.UserId)
	if err != nil {
		d.logger.Error("productUsecase.GetOrderedProductsByUserID", zap.Error(err))
		return &pb.ListProductResponse{}, grpc.Error(ctx, err)
	}

	var pbProducts []*pb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, &pb.Product{
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
			Size_:       product.Size,
			CreatedAt:   product.CreatedAt.String(),
			UpdatedAt:   product.UpdatedAt.String(),
		})
	}
	return &pb.ListProductResponse{Products: pbProducts}, nil
}
