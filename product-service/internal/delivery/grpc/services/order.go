package services

import (
	"context"
	"log"
	pb "product-service/genproto/product_service"
	grpc "product-service/internal/delivery"
	"product-service/internal/entity"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (d *productRPC) CreateOrder(ctx context.Context, in *pb.Order) (*pb.GetWithID, error) {
	id := uuid.New().String()
	order, err := d.productUsecase.CreateOrder(ctx, &entity.Order{
		Id:        id,
		ProductID: in.ProductId,
		UserID:    in.UserId,
		Status:    in.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		d.logger.Error("productUsecase.CreateOrder", zap.Error(err))
		return &pb.GetWithID{}, grpc.Error(ctx, err)
	}

	in.Id = id
	return &pb.GetWithID{Id: order.Id}, nil
}

func (d *productRPC) CancelOrder(ctx context.Context, in *pb.GetWithID) (*pb.DeleteResponse, error) {
	err := d.productUsecase.CancelOrder(ctx, in.Id)
	if err != nil {
		d.logger.Error("productUsecase.CancelOrder", zap.Error(err))
		return &pb.DeleteResponse{Status: false}, grpc.Error(ctx, err)
	}

	return &pb.DeleteResponse{Status: true}, nil
}

func (d *productRPC) GetOrderByID(ctx context.Context, in *pb.GetWithID) (*pb.Order, error) {
	order, err := d.productUsecase.GetOrderByID(ctx, map[string]string{"id": in.Id})
	if err != nil {
		log.Println(err.Error())
		d.logger.Error("productUsecase.GetOrderByID", zap.Error(err))
		return &pb.Order{}, grpc.Error(ctx, err)
	}
	return &pb.Order{
		Id:        order.Id,
		ProductId: order.ProductID,
		UserId:    order.UserID,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
		UpdatedAt: order.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (d *productRPC) GetAllOrders(ctx context.Context, in *pb.ListRequest) (*pb.ListOrderResponse, error) {
	filter := &entity.ListRequest{
		Limit: in.Limit,
		Page:  in.Page,
	}

	orders, err := d.productUsecase.GetAllOrders(ctx, filter)
	if err != nil {
		d.logger.Error("productUsecase.GetAllOrders", zap.Error(err))
		return &pb.ListOrderResponse{}, err
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
