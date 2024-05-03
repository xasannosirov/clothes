package v1

import (
	grpc_service_clients "api-gateway/internal/infrastructure/grpc_service_client"
	"api-gateway/internal/pkg/config"
	"time"

	"go.uber.org/zap"
)

type HandlerV1 struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpc_service_clients.ServiceClient
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpc_service_clients.ServiceClient
}

// New ...
func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		Config:         c.Config,
		Logger:         c.Logger,
		Service:        c.Service,
		ContextTimeout: c.ContextTimeout,
	}
}
