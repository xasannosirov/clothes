package grpc_client

import (
	"fmt"

	pbm "api-gateway/genproto/media_service"
	pbp "api-gateway/genproto/product_service"
	pbu "api-gateway/genproto/user_service"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"api-gateway/internal/pkg/config"
)

type ServiceClient interface {
	UserService() pbu.UserServiceClient
	MediaService() pbm.MediaServiceClient
	ProductService() pbp.ProductServiceClient
	Close()
}

type serviceClient struct {
	connections    []*grpc.ClientConn
	userService    pbu.UserServiceClient
	mediaService   pbm.MediaServiceClient
	productService pbp.ProductServiceClient
}

func New(cfg *config.Config) (ServiceClient, error) {
	connUserService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.UserService.Host, cfg.UserService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	connMediaService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.MediaService.Host, cfg.MediaService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	connProductService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.ProductService.Host, cfg.ProductService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return &serviceClient{
		userService:    pbu.NewUserServiceClient(connUserService),
		mediaService:   pbm.NewMediaServiceClient(connMediaService),
		productService: pbp.NewProductServiceClient(connProductService),
		connections: []*grpc.ClientConn{
			connUserService,
			connMediaService,
			connProductService,
		},
	}, nil
}

func (s *serviceClient) UserService() pbu.UserServiceClient {
	return s.userService
}
func (s *serviceClient) MediaService() pbm.MediaServiceClient {
	return s.mediaService
}
func (s *serviceClient) ProductService() pbp.ProductServiceClient {
	return s.productService
}

func (s *serviceClient) Close() {
	for _, conn := range s.connections {
		if err := conn.Close(); err != nil {
			// should be replaced by logger soon
			fmt.Printf("error while closing grpc connection: %v", err)
		}
	}
}
