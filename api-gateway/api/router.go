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
// @Description Online Clothes Store
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

	// registration
	apiV1.POST("/register", HandlerV1.Register)
	apiV1.POST("/login", HandlerV1.Login)
	apiV1.POST("/forgot/:email", HandlerV1.Forgot)
	apiV1.POST("/verify", HandlerV1.Verify)
	apiV1.PUT("/reset-password", HandlerV1.ResetPassword)
	apiV1.GET("/token/:refresh", HandlerV1.Token)

	// media
	apiV1.POST("/media/upload-photo", HandlerV1.UploadMedia)
	apiV1.GET("/media/:id", HandlerV1.GetMedia)
	apiV1.DELETE("/media/:id", HandlerV1.DeleteMedia)

	// users
	apiV1.POST("/user", HandlerV1.CreateUser)
	apiV1.PUT("/user", HandlerV1.UpdateUser)
	apiV1.DELETE("/user/:id", HandlerV1.DeleteUser)
	apiV1.GET("/user/:id", HandlerV1.GetUser)
	apiV1.GET("/users", HandlerV1.ListUsers)

	// products
	apiV1.POST("/product", HandlerV1.CreateProduct)
	apiV1.PUT("/product", HandlerV1.UpdateProduct)
	apiV1.DELETE("/product/:id", HandlerV1.DeleteProduct)
	apiV1.GET("/product/:id", HandlerV1.GetProduct)
	apiV1.GET("/products", HandlerV1.ListProducts)

	apiV1.POST("/order", HandlerV1.CreateOrder)
	apiV1.GET("/order/:id", HandlerV1.GetOrder)
	apiV1.DELETE("/order/:id", HandlerV1.CancelOrder)
	apiV1.GET("/orders", HandlerV1.ListOrders)

	apiV1.POST("/like-product", HandlerV1.LikeProduct)
	apiV1.POST("/save-product", HandlerV1.SaveProduct)
	apiV1.POST("/star-product", HandlerV1.StarToProduct)
	apiV1.POST("/comment-product", HandlerV1.CommentToProduct)

	apiV1.GET("/comments", HandlerV1.GetAllComments)
	apiV1.GET("/stars", HandlerV1.GetAllStars)

	apiV1.GET("/product/orders/:id", HandlerV1.GetProductOrders)
	apiV1.GET("/product/comments/:id", HandlerV1.GetProductComments)
	apiV1.GET("/product/likes/:id", HandlerV1.GetProductLikes)
	apiV1.GET("/product/stars/:id", HandlerV1.GetProductStars)

	apiV1.GET("/user/save/:id", HandlerV1.GetUserSavedProducts)
	apiV1.GET("/user/likes/:id", HandlerV1.GetUserLikesProducts)
	apiV1.GET("/user/orders/:id", HandlerV1.GetUserOrderedProducts)

	apiV1.GET("/search/:name", HandlerV1.SearchProduct)
	apiV1.GET("/recommendation", HandlerV1.RecommendProducts)
	apiV1.GET("/disable-orders", HandlerV1.GetDisableProducts)

	url := ginSwagger.URL("swagger/doc.json")
	apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
