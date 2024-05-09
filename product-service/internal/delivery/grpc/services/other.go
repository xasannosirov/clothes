package services

import (
	"context"
	pb "product-service/genproto/product_service"
	grpc "product-service/internal/delivery"
	"product-service/internal/entity"
	"time"

	"go.uber.org/zap"
)

func (d *productRPC) SearchProduct(ctx context.Context, in *pb.Filter) (*pb.ListProductResponse, error) {
	products, err := d.productUsecase.SearchProduct(ctx, &entity.Filter{
		Name: in.Name,
	})
	if err != nil {
		return nil, err
	}

	var serviceResponse pb.ListProductResponse
	for _, product := range products.Products {
		serviceResponse.Products = append(serviceResponse.Products, &pb.Product{
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
			CreatedAt:   product.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
		})
	}

	serviceResponse.TotalCount = products.TotalCount

	return &serviceResponse, nil
}

func (d *productRPC) Recommendation(ctx context.Context, in *pb.Recom) (*pb.ListProductResponse, error) {
	products, err := d.productUsecase.RecommentProducts(ctx, &entity.Recom{Age: uint8(in.Age), Gender: in.Gender})
	if err != nil {
		d.logger.Error("productUseCase.Recommendation", zap.Error(err))
		return &pb.ListProductResponse{}, grpc.Error(ctx, err)
	}

	pbProducts := &pb.ListProductResponse{}
	for _, product := range products.Products {
		pbProducts.Products = append(pbProducts.Products, &pb.Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Category:    product.Category,
			MadeIn:      product.MadeIn,
			Color:       product.Color,
			Count:       int64(product.Count),
			Cost:        float32(product.Cost),
			Discount:    float32(product.Discount),
			AgeMin:      product.AgeMin,
			AgeMax:      product.AgeMax,
			ForGender:   product.ForGender,
			Size_:       product.Size,
			CreatedAt:   product.CreatedAt.String(),
			UpdatedAt:   product.UpdatedAt.String(),
		})
	}

	return &pb.ListProductResponse{
		Products:   pbProducts.Products,
		TotalCount: products.TotalCount,
	}, nil
}

func (d *productRPC) GetDisableProducts(ctx context.Context, in *pb.ListRequest) (*pb.ListOrderResponse, error) {
	orders, err := d.productUsecase.GetDisableProducts(ctx, &entity.ListRequest{Page: in.Page, Limit: in.Limit})
	if err != nil {
		d.logger.Error("productUsecase.GetDisableProducts", zap.Error(err))
		return &pb.ListOrderResponse{}, grpc.Error(ctx, err)
	}

	pbOrders := &pb.ListOrderResponse{}
	for _, order := range orders.Orders {
		pbOrders.Orders = append(pbOrders.Orders, &pb.Order{
			Id:        order.Id,
			ProductId: order.ProductID,
			UserId:    order.UserID,
			Status:    order.Status,
			CreatedAt: order.CreatedAt.Format(time.RFC3339),
			UpdatedAt: order.UpdatedAt.Format(time.RFC3339),
		})
	}
	return &pb.ListOrderResponse{
		Orders:     pbOrders.Orders,
		TotalCount: orders.TotalCount,
	}, nil
}
