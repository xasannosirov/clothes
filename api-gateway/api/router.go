package api

import (
	"time"

	_ "api-gateway/api/docs"
	v1 "api-gateway/api/handlers/v1"
	"api-gateway/api/middleware"

	redisrepo "api-gateway/internal/infrastructure/repository/redis"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	grpcClients "api-gateway/internal/infrastructure/grpc_client"
	"api-gateway/internal/pkg/config"
	"api-gateway/internal/pkg/token"
)

type RouteOption struct {
	Config         config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpcClients.ServiceClient
	Cache          redisrepo.Cache
	Enforcer       *casbin.Enforcer
	RefreshToken   token.JWTHandler
}

// NewRoute
// @securityDefinitions.apikey BearerAuth
// @in 			header
// @name 		Authorization
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

	router.Use(middleware.Tracing)
	router.Use(middleware.CheckCasbinPermission(option.Enforcer, option.Config))

	router.Static("/media", "./media")

	apiV1 := router.Group("/v1")

	// Auth
	apiV1.POST("/register", HandlerV1.Register)
	apiV1.POST("/login", HandlerV1.Login)
	apiV1.POST("/forgot/:email", HandlerV1.Forgot)
	apiV1.POST("/verify", HandlerV1.Verify)
	apiV1.PUT("/reset-password", HandlerV1.ResetPassword)
	apiV1.PUT("/update-password", HandlerV1.UpdatePassword)
	apiV1.GET("/token/:refresh", HandlerV1.NewToken)
	apiV1.GET("/google/login", HandlerV1.GoogleLogin)
	apiV1.GET("google/callback", HandlerV1.GoogleCallback)

	// User
	apiV1.POST("/user", HandlerV1.CreateUser)
	apiV1.PUT("/user", HandlerV1.UpdateUser)
	apiV1.DELETE("/user/:id", HandlerV1.DeleteUser)
	apiV1.GET("/user/:id", HandlerV1.GetUser)
	apiV1.GET("/users", HandlerV1.ListUsers)

	// Worker
	apiV1.POST("/worker", HandlerV1.CreateWorker)
	apiV1.PUT("/worker", HandlerV1.UpdateWorker)
	apiV1.DELETE("/worker/:id", HandlerV1.DeleteWorker)
	apiV1.GET("/worker/:id", HandlerV1.GetWorker)
	apiV1.GET("/workers", HandlerV1.ListWorker)

	// Category
	apiV1.POST("/category", HandlerV1.CreateCategory)
	apiV1.PUT("/category", HandlerV1.UpdateCategory)
	apiV1.DELETE("/category/:id", HandlerV1.DeleteCategory)
	apiV1.GET("/category/:id", HandlerV1.GetCategory)
	apiV1.GET("/categories", HandlerV1.ListCategory)
	apiV1.GET("/category/search", HandlerV1.SearchCategory)

	// Product
	apiV1.POST("/product", HandlerV1.CreateProduct)
	apiV1.PUT("/product", HandlerV1.UpdateProduct)
	apiV1.DELETE("/product/:id", HandlerV1.DeleteProduct)
	apiV1.GET("/product/:id", HandlerV1.GetProduct)
	apiV1.GET("/products", HandlerV1.ListProducts)
	apiV1.GET("/products/discount", HandlerV1.GetDicountProducts)

	// Media
	apiV1.POST("/media/upload-photo", HandlerV1.UploadMedia)
	apiV1.GET("/media/:id", HandlerV1.GetMedia)
	apiV1.DELETE("/media/:id", HandlerV1.DeleteMedia)

	// Wishlist
	apiV1.POST("/like/:id", HandlerV1.LikeProduct)
	apiV1.GET("/wishlist", HandlerV1.UserWishlist)

	// Basket, Stats, Payment, Order ...
	apiV1.POST("/basket", HandlerV1.SaveToBasket)
	apiV1.GET("/basket/:id", HandlerV1.GetBasketProduct)
	apiV1.DELETE("/basket/:id", HandlerV1.DeleteFromBasket)

	url := ginSwagger.URL("swagger/doc.json")
	apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
