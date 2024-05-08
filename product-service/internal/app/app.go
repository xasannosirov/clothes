package app

import (
	"fmt"
	pb "product-service/genproto/product_service"
	"product-service/internal/delivery/grpc/server"
	"product-service/internal/delivery/grpc/services"
	grpc_service_clients "product-service/internal/infrastructure/grpc_service_client"
	"time"

	repo "product-service/internal/infrastructure/repository/postgresql"
	"product-service/internal/pkg/config"
	"product-service/internal/pkg/logger"
	"product-service/internal/pkg/otlp"
	"product-service/internal/pkg/postgres"
	"product-service/internal/usecase"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	ServiceClients grpc_service_clients.ServiceClients
	GrpcServer     *grpc.Server
	ShutdownOTLP   func() error
}

func NewApp(cfg *config.Config) (*App, error) {
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	shutdownOTLP, err := otlp.InitOTLPProvider(cfg)
	if err != nil {
		return nil, err
	}

	db, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	clients, err := grpc_service_clients.New(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:         cfg,
		Logger:         logger,
		DB:             db,
		GrpcServer:     grpcServer,
		ShutdownOTLP:   shutdownOTLP,
		ServiceClients: clients,
	}, nil
}

func (a *App) Run() error {
	var (
		contextTimeout time.Duration
	)

	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error during parse duration for context timeout : %w", err)
	}

	serviceClients, err := grpc_service_clients.New(a.Config)
	if err != nil {
		return fmt.Errorf("error during initialize service clients: %w", err)
	}
	a.ServiceClients = serviceClients

	productRepo := repo.NewProductsRepo(a.DB)

	productUseCase := usecase.NewProductService(contextTimeout, productRepo)

	pb.RegisterProductServiceServer(a.GrpcServer, services.NewRPC(a.Logger, productUseCase, a.ServiceClients))

	a.Logger.Info("gRPC Server Listening", zap.String("url", a.Config.RPCPort))
	if err := server.Run(a.Config, a.GrpcServer); err != nil {
		return fmt.Errorf("gRPC fatal to serve grpc server over %s %w", a.Config.RPCPort, err)
	}
	return nil
}

func (a *App) Stop() {

	// closing client service connections
	a.ServiceClients.Close()
	// stop gRPC server
	a.GrpcServer.Stop()

	// database connection
	a.DB.Close()

	// zap logger sync
	a.Logger.Sync()
}
