package grpcserviceclient

import (
	"fmt"
	mediaproto "product-service/genproto/media_service"
	userproto "product-service/genproto/user_service"
	"product-service/internal/pkg/config"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type ServiceClients interface {
	UserService() userproto.UserServiceClient
	MediaService() mediaproto.MediaServiceClient
	Close()
}

type serviceClients struct {
	services     []*grpc.ClientConn
	userService  userproto.UserServiceClient
	mediaService mediaproto.MediaServiceClient
}

func New(config *config.Config) (ServiceClients, error) {
	userCONNECTION, err := grpc.Dial(
		fmt.Sprintf("%s%s", config.UserService.Host, config.UserService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	mediaCONNECTION, err := grpc.Dial(
		fmt.Sprintf("%s%s", config.MediaService.Host, config.MediaService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return &serviceClients{
		services: []*grpc.ClientConn{
			userCONNECTION,
			mediaCONNECTION,
		},
		userService:  userproto.NewUserServiceClient(userCONNECTION),
		mediaService: mediaproto.NewMediaServiceClient(mediaCONNECTION),
	}, nil
}

func (s *serviceClients) UserService() userproto.UserServiceClient {
	return s.userService
}

func (s *serviceClients) MediaService() mediaproto.MediaServiceClient {
	return s.mediaService
}

func (s *serviceClients) Close() {
	// closing investment service
	for _, conn := range s.services {
		conn.Close()
	}
}
