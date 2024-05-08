package services

import (
	"context"
	grpc "product-service/internal/delivery"
	pb "product-service/genproto/product_service"
	"product-service/internal/entity"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (d *productRPC) LikeProduct(ctx context.Context, req *pb.Like) (*pb.MoveResponse, error) {
	status, err := d.productUsecase.LikeProduct(ctx, &entity.LikeProduct{
		ProductID: req.ProductId,
		UserID:    req.UserId,
	})

	if err != nil {
		d.logger.Error("productUseCase.CreateProduct", zap.Error(err))
		return nil, err
	}

	return &pb.MoveResponse{
		Status: status,
	}, nil
}

func (d *productRPC) SaveProduct(ctx context.Context, req *pb.Save) (*pb.MoveResponse, error) {
	status, err := d.productUsecase.SaveProduct(ctx, &entity.SaveProduct{
		Id:        req.Id,
		ProductID: req.ProductId,
		UserID:    req.UserId,
	})

	if err != nil {
		d.logger.Error("productUseCase.SaveProduct", zap.Error(err))
		return nil, err
	}

	return &pb.MoveResponse{
		Status: status,
	}, nil
}

func (d *productRPC) StarProduct(ctx context.Context, in *pb.Star) (*pb.MoveResponse, error) {
	id := uuid.New().String()
	_, err := d.productUsecase.StarProduct(ctx, &entity.StarProduct{
		Id:        id,
		ProductID: in.ProductId,
		UserID:    in.UserId,
		Stars:     in.Star,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		d.logger.Error("productUsecase.StarProduct", zap.Error(err))
		return &pb.MoveResponse{}, grpc.Error(ctx, err)
	}

	in.Id = id
	return &pb.MoveResponse{Status: true}, nil
}

func (d *productRPC) CommentToProduct(ctx context.Context, req *pb.Comment) (*pb.MoveResponse, error) {
	status, err := d.productUsecase.CommentToProduct(ctx, &entity.CommentToProduct{
		UserID:    req.UserId,
		ProductID: req.ProductId,
		Comment:   req.Comment,
	})
	if err != nil {
		return nil, err
	}
	return &pb.MoveResponse{
		Status: status,
	}, nil
}
