package services

import (
	pb "product-service/genproto/product_service"
	grpcserviceclient "product-service/internal/infrastructure/grpc_service_client"
	"product-service/internal/usecase"

	"go.uber.org/zap"
)

type productRPC struct {
	logger         *zap.Logger
	productUsecase usecase.Product
	services       grpcserviceclient.ServiceClients
}

func NewRPC(logger *zap.Logger, productUsecase usecase.Product, clients grpcserviceclient.ServiceClients) pb.ProductServiceServer {
	return &productRPC{
		logger:         logger,
		productUsecase: productUsecase,
		services:       clients,
	}
}
