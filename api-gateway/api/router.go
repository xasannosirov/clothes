package api

import (
	"time"

	v1 "api-gateway/api/handlers/v1"
	"api-gateway/api/middleware"
	_ "api-gateway/api/docs"

	redisrepo "api-gateway/internal/infrastructure/repository/redis"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	grpcClients "api-gateway/internal/infrastructure/grpc_client"
	"api-gateway/internal/pkg/config"
	"api-gateway/internal/usecase/refresh_token"
)

type RouteOption struct {
	Config         config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpcClients.ServiceClient
	Cache          redisrepo.Cache
	Enforcer       *casbin.Enforcer
	RefreshToken   refresh_token.JWTHandler
}

// NewRoute
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRoute(option RouteOption) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	HandlerV1 := v1.New(&v1.HandlerV1Config{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		Service:        option.Service,
		Redis:          option.Cache,
		RefreshToken:   option.RefreshToken,
		Enforcer:       option.Enforcer,
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))
	// router.Use()
	router.Use(middleware.Tracing)
	router.Use(middleware.CheckCasbinPermission(option.Enforcer, option.Config))

	router.Static("/media", "./media")

	api := router.Group("/v1")

	// register verify login
	api.POST("/users/register", HandlerV1.Register)
	api.POST("/users/verify", HandlerV1.Verify)
	api.POST("/users/login", HandlerV1.LoginUser)
	api.POST("/users/token", HandlerV1.Token)
	api.POST("/users/forgetpassword", HandlerV1.ForgetPassword)
	api.POST("/users/verify/forgetpassword", HandlerV1.VerifyForgetPassword)

	//photo
	api.POST("/media/photo", HandlerV1.UploadMedia)
	api.GET("/media/get/:id", HandlerV1.GetMedia)
	api.DELETE("/media/delete/:id", HandlerV1.DeleteMedia)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
