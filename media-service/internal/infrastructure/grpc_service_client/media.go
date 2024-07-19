package grpc_service_clients

import (
	"fmt"
	"media-service/internal/pkg/config"

	paymentproto "media-service/genproto/payment_service"
	productproto "media-service/genproto/product_service"
	userproto "media-service/genproto/user_service"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type ServiceClients interface {
	UserService() userproto.UserServiceClient
	ProductService() productproto.ProductServiceClient
	PaymentService() paymentproto.PaymentServiceClient
	Close()
}

type serviceClients struct {
	userService    userproto.UserServiceClient
	productService productproto.ProductServiceClient
	paymentService paymentproto.PaymentServiceClient
	services       []*grpc.ClientConn
}

func New(config *config.Config) (ServiceClients, error) {

	userServiceConnnection, err := grpc.Dial(
		fmt.Sprintf("%s%s", config.UserService.Host, config.UserService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	productServiceConnection, err := grpc.Dial(
		fmt.Sprintf("%s%s", config.ProductService.Host, config.ProductService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return &serviceClients{
		userService:    userproto.NewUserServiceClient(userServiceConnnection),
		productService: productproto.NewProductServiceClient(productServiceConnection),
		services:       []*grpc.ClientConn{},
	}, nil
}

func (s *serviceClients) Close() {
	// closing store-management service
	for _, conn := range s.services {
		conn.Close()
	}
}

func (s *serviceClients) UserService() userproto.UserServiceClient {
	return s.userService
}

func (s *serviceClients) ProductService() productproto.ProductServiceClient {
	return s.productService
}

func (s *serviceClients) PaymentService() paymentproto.PaymentServiceClient {
	return s.paymentService
}
