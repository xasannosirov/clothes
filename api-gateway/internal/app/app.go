package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	// "github.com/casbin/casbin/v2"
	"api-gateway/api"

	"github.com/casbin/casbin/v2/util"

	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"go.uber.org/zap"

	grpcService "api-gateway/internal/infrastructure/grpc_client"

	// "exam_5/api_gateway/internal/infrastructure/kafka"
	// "exam_5/api_gateway/internal/infrastructure/repository/postgresql"
	redisrepo "api-gateway/internal/infrastructure/repository/redis"
	"api-gateway/internal/pkg/config"
	"api-gateway/internal/pkg/logger"

	// "exam_5/api_gateway/internal/pkg/policy"

	"api-gateway/internal/pkg/otlp"
	// "exam_5/api_gateway/internal/pkg/policy"
	// "exam_5/api_gateway/internal/pkg/postgres"
	"api-gateway/internal/pkg/storage/redis"
	// "exam_5/api_gateway/internal/usecase/app_version"
	// "exam_5/api_gateway/internal/usecase/event"
	// "evrone_api_gateway/internal/usecase/refresh_token"
)

type App struct {
	Config       config.Config
	Logger       *zap.Logger
	server       *http.Server
	RedisDB      *redis.RedisDB
	ShutdownOTLP func() error
	Clients      grpcService.ServiceClient
	Enforcer     *casbin.Enforcer
}

func NewApp(cfg config.Config) (*App, error) {
	// logger init
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	// redis init
	redisdb, err := redis.New(&cfg)
	if err != nil {
		return nil, err
	}
	enforcer, err := casbin.NewEnforcer("auth.conf", "auth.csv")
	if err != nil {
		return nil, err
	}

	// otlp collector init
	shutdownOTLP, err := otlp.InitOTLPProvider(&cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:       cfg,
		Logger:       logger,
		ShutdownOTLP: shutdownOTLP,
		RedisDB:      redisdb,
		Enforcer:     enforcer,
	}, nil
}

func (a *App) Run() error {
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error while parsing context timeout: %v", err)
	}

	clients, err := grpcService.New(&a.Config)
	if err != nil {
		return err
	}
	a.Clients = clients

	// initialize cache
	cache := redisrepo.NewCache(a.RedisDB)

	// api init
	handler := api.NewRoute(api.RouteOption{
		Config:         a.Config,
		Logger:         a.Logger,
		ContextTimeout: contextTimeout,
		Service:        clients,
		Cache:          cache,
		Enforcer:       a.Enforcer,
	})
	err = a.Enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	roleManager := a.Enforcer.GetRoleManager().(*defaultrolemanager.RoleManagerImpl)

	roleManager.AddMatchingFunc("keyMatch", util.KeyMatch)
	roleManager.AddMatchingFunc("keyMatch3", util.KeyMatch3)

	// server init
	a.server, err = api.NewServer(&a.Config, handler)
	if err != nil {
		return fmt.Errorf("error while initializing server: %v", err)
	}

	return a.server.ListenAndServe()
}

func (a *App) Stop() {

	// close grpc connections
	a.Clients.Close()

	// shutdown server http
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.Logger.Error("shutdown server http ", zap.Error(err))
	}

	// shutdown otlp collector
	if err := a.ShutdownOTLP(); err != nil {
		a.Logger.Error("shutdown otlp collector", zap.Error(err))
	}

	// zap logger sync
	a.Logger.Sync()
}
