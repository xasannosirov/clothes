package grpc_service_clients

import (
	"fmt"
	mediaproto "user-service/genproto/media_service"
	paymentproto "user-service/genproto/payment_service"
	productproto "user-service/genproto/product_service"
	"user-service/internal/pkg/config"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type ServiceClients interface {
	ProductService() productproto.ProductServiceClient
	MediaService() mediaproto.MediaServiceClient
	PaymentService() paymentproto.PaymentServiceClient
	Close()
}

type serviceClients struct {
	productService productproto.ProductServiceClient
	mediaService   mediaproto.MediaServiceClient
	paymentService paymentproto.PaymentServiceClient
	services       []*grpc.ClientConn
}

func New(config *config.Config) (ServiceClients, error) {

	connProductService, err := grpc.Dial(
		fmt.Sprintf("%s%s", config.ProductService.Host, config.ProductService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	connMediaService, err := grpc.Dial(
		fmt.Sprintf("%s%s", config.MediaService.Host, config.MediaService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	connPaymentService, err := grpc.Dial(
		fmt.Sprintf("%s%s", config.PaymentService.Host, config.PaymentService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return &serviceClients{
		productService: productproto.NewProductServiceClient(connProductService),
		mediaService:   mediaproto.NewMediaServiceClient(connMediaService),
		paymentService: paymentproto.NewPaymentServiceClient(connPaymentService),
		services:       []*grpc.ClientConn{},
	}, nil
}

func (s *serviceClients) Close() {
	// closing investment service
	for _, conn := range s.services {
		conn.Close()
	}
}

func (s *serviceClients) ProductService() productproto.ProductServiceClient {
	return s.productService
}

func (s *serviceClients) MediaService() mediaproto.MediaServiceClient {
	return s.mediaService
}

func (s *serviceClients) PaymentService() paymentproto.PaymentServiceClient {
	return s.paymentService
}
