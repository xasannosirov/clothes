package config

import (
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	OtpSecret = "some_secret"
)

type webAddress struct {
	Host string
	Port string
}

type Config struct {
	APP         string
	Environment string
	LogLevel    string

	Server struct {
		Host         string
		Port         string
		ReadTimeout  string
		WriteTimeout string
		IdleTimeout  string
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SSLMode  string
	}

	Context struct {
		Timeout string
	}

	Redis struct {
		Host     string
		Port     string
		Password string
		Name     string
		Time     time.Time
	}

	Token struct {
		Secret     string
		AccessTTL  time.Duration
		RefreshTTL time.Duration
		SignInKey  string
	}

	Minio struct {
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
		Location        string
		BucketName      string
	}

	SMTP struct {
		Email         string
		EmailPassword string
		SMTPPort      string
		SMTPHost      string
	}

	Google struct {
		ClientId     string
		ClientSecret string
		RedirectURL  string
	}

	UserService    webAddress
	MediaService   webAddress
	ProductService webAddress
}

func NewConfig() (*Config, error) {
	var config Config

	// general configuration
	config.APP = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// server configuration
	config.Server.Host = getEnv("SERVER_HOST", "clothes-api") // api
	config.Server.Port = getEnv("SERVER_PORT", ":5555")
	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	// redis configuration
	config.Redis.Host = getEnv("REDIS_HOST", "clothes-redis-db") // redisdb
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.Name = getEnv("REDIS_DATABASE", "0")

	//user service
	config.UserService.Host = getEnv("USER_SERVICE_GRPC_HOST", "clothes-user-service") // user-service
	config.UserService.Port = getEnv("USER_SERVICE_GRPC_PORT", ":1111")

	//media service
	config.MediaService.Host = getEnv("MEDIA_SERVICE_GRPC_HOST", "clothes-media-service") //media-service
	config.MediaService.Port = getEnv("MEDIA_SERVICE_GRPC_PORT", ":2222")

	//product service
	config.ProductService.Host = getEnv("PRODUCT_SERVICE_GRPC_HOST", "clothes-product-service") // product-service
	config.ProductService.Port = getEnv("PRODUCT_SERVICE_GRPC_PORT", ":3333")

	// token configuration
	config.Token.Secret = getEnv("TOKEN_SECRET", "token_secret")

	// access ttl parse
	accessTTl, err := time.ParseDuration(getEnv("TOKEN_ACCESS_TTL", "180s"))
	if err != nil {
		return nil, err
	}
	// refresh ttl parse
	refreshTTL, err := time.ParseDuration(getEnv("TOKEN_REFRESH_TTL", "240h"))
	if err != nil {
		return nil, err
	}
	config.Token.AccessTTL = accessTTl
	config.Token.RefreshTTL = refreshTTL
	config.Token.SignInKey = getEnv("TOKEN_SIGNING_KEY", "debug")

	//smtp configuration
	config.SMTP.Email = getEnv("SMTP_EMAIL", "storegoclothes@gmail.com")
	config.SMTP.EmailPassword = getEnv("SMTP_EMAIL_PASSWORD", "dnrzomyvooftjcgc ")
	config.SMTP.SMTPPort = getEnv("SMTP_PORT", "587")
	config.SMTP.SMTPHost = getEnv("SMTP_HOST", "smtp.gmail.com")

	//minIO configuration
	config.Minio.AccessKeyID = getEnv("ACCESS_KEY_ID", "abdulaziz")
	config.Minio.SecretAccessKey = getEnv("SECRET_ACCESS_KEY", "abdulaziz")
	config.Minio.Endpoint = getEnv("ENDPOINT", "138.68.146.55:9000") // 13.201.56.179:9000
	config.Minio.BucketName = getEnv("BUCKET_NAME", "clothesstore")

	return &config, nil
}

func SetupConfig() *oauth2.Config {

	conf := &oauth2.Config{
		ClientID:     getEnv("CLIENT_ID", "627168078285-v33r2ijbkjvpsonc85nt1ns7rji08jti.apps.googleusercontent.com"),
		ClientSecret: getEnv("CLIENT_SECRET", "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"),
		RedirectURL:  getEnv("REDIRECT_URL", "http://localhost:5555/v1/google/callback"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
