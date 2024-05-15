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

func (d *productRPC) CreateProduct(ctx context.Context, in *pb.Product) (*pb.GetWithID, error) {
	ctx, span := otlp.Start(ctx, "product_grpc-delivery", "CreateProduct")
	span.SetAttributes(
		attribute.Key("guid").String(in.Id),
	)
	defer span.End()

	respProduct, err := d.productUsecase.CreateProduct(ctx, &entity.Product{
		Id:          in.Id,
		Name:        in.Name,
		Description: in.Description,
		Category:    in.Category,
		MadeIn:      in.MadeIn,
		Color:       in.Color,
		Count:       in.Count,
		Cost:        in.Cost,
		Discount:    in.Discount,
		AgeMin:      in.AgeMin,
		AgeMax:      in.AgeMax,
		ForGender:   in.ForGender,
		Size:        int64(in.Size()),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		d.logger.Error("productUseCase.CreateProduct", zap.Error(err))
		return &pb.GetWithID{}, grpc.Error(ctx, err)
	}
	return &pb.GetWithID{Id: respProduct.Id}, nil
}

func (d *productRPC) UpdateProduct(ctx context.Context, in *pb.Product) (*pb.Product, error) {
	err := d.productUsecase.UpdateProduct(ctx, &entity.Product{
		Id:          in.Id,
		Name:        in.Name,
		Description: in.Description,
		Category:    in.Category,
		MadeIn:      in.MadeIn,
		Color:       in.Color,
		Count:       in.Count,
		Cost:        in.Cost,
		Discount:    in.Discount,
		AgeMin:      in.AgeMin,
		AgeMax:      in.AgeMax,
		ForGender:   in.ForGender,
		Size:        int64(in.Size()),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		d.logger.Error("productUseCase.UpdateProduct", zap.Error(err))
		return &pb.Product{}, grpc.Error(ctx, err)
	}

	return in, nil
}

func (d *productRPC) GetProductByID(ctx context.Context, in *pb.GetWithID) (*pb.Product, error) {
	product, err := d.productUsecase.GetProductByID(ctx, map[string]string{"id": in.Id})
	if err != nil {
		d.logger.Error("productUseCase.GetProductByID", zap.Error(err))
		return &pb.Product{}, grpc.Error(ctx, err)
	}

	return &pb.Product{
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
	}, nil
}
func (d *productRPC) GetProductDelete(ctx context.Context, in *pb.GetWithID) (*pb.Product, error) {
	product, err := d.productUsecase.GetProductDelete(ctx, map[string]string{"id": in.Id})
	if err != nil {
		d.logger.Error("productUseCase.GetProductByID", zap.Error(err))
		return &pb.Product{}, grpc.Error(ctx, err)
	}
	return &pb.Product{
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
	}, nil
}

func (d *productRPC) DeleteProduct(ctx context.Context, in *pb.GetWithID) (*pb.DeleteResponse, error) {
	err := d.productUsecase.DeleteProduct(ctx, in.Id)
	if err != nil {
		d.logger.Error("productUseCase.DeleteProduct", zap.Error(err))
		return &pb.DeleteResponse{Status: false}, err
	}

	return &pb.DeleteResponse{Status: true}, nil
}

func (d *productRPC) GetAllProducts(ctx context.Context, in *pb.ListProductRequest) (*pb.ListProductResponse, error) {
	filter := &entity.ListProductRequest{
		Limit: in.Limit,
		Page:  in.Page,
		Name:  in.Name,
	}

	products, err := d.productUsecase.GetAllProducts(ctx, filter)
	if err != nil {
		d.logger.Error("productUseCase.List", zap.Error(err))
		return nil, err
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
