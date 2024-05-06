package v1

import (
	grpcClients "api-gateway/internal/infrastructure/grpc_client"
	repo "api-gateway/internal/infrastructure/repository/redis"
	"api-gateway/internal/pkg/config"
	"api-gateway/internal/usecase/refresh_token"
	"time"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
)

type HandlerV1 struct {
	Config         config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpcClients.ServiceClient
	redisStorage   repo.Cache
	RefreshToken   refresh_token.JWTHandler
	Enforcer       *casbin.Enforcer
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Config         config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpcClients.ServiceClient
	Redis          repo.Cache
	RefreshToken   refresh_token.JWTHandler
	Enforcer       *casbin.Enforcer
}

// New ...
func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		Config:         c.Config,
		Logger:         c.Logger,
		Service:        c.Service,
		ContextTimeout: c.ContextTimeout,
		redisStorage:   c.Redis,
		Enforcer:       c.Enforcer,
		RefreshToken:   c.RefreshToken,
	}
}
