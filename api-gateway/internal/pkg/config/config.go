package config

import (
	"os"
	"strings"

	// "strings"
	"time"
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
	Server      struct {
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
		Endpoint              string
		AccessKey             string
		SecretKey             string
		Location              string
		MovieUploadBucketName string
	}
	Kafka struct {
		Address []string
		Topic   struct {
			UserCreateTopic string
		}
	}
	SMTP struct {
		Email         string
		EmailPassword string
		SMTPPort      string
		SMTPHost      string
	}
	UserService    webAddress
	MediaService   webAddress
	PaymentServie  webAddress
	ProductService webAddress
	OTLPCollector  webAddress
}

func NewConfig() (*Config, error) {
	var config Config

	// general configuration
	config.APP = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// server configuration
	config.Server.Host = getEnv("SERVER_HOST", "localhost")
	config.Server.Port = getEnv("SERVER_PORT", ":5555")
	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "localhost")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "examdb")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "4444")
	config.DB.SSLMode = getEnv("POSTGRES_SSLMODE", "disable")

	// redis configuration
	config.Redis.Host = getEnv("REDIS_HOST", "localhost")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.Name = getEnv("REDIS_DATABASE", "0")

	//user service
	config.UserService.Host = getEnv("USER_SERVICE_HOST", "localhost")
	config.UserService.Port = getEnv("USER_SERVICE_PORT", ":1111")

	//media service
	config.MediaService.Host = getEnv("MEDIA_SERVICE_HOST", "localhost")
	config.MediaService.Port = getEnv("MEDIA_SERVICE_PORT", ":2222")
 
	//product service
	config.ProductService.Host = getEnv("PRODUCT_SERVICE_HOST", "localhost")
	config.ProductService.Port = getEnv("PRODUCT_SERVICE_PORT", ":3333")

	//payment servicve
	config.PaymentServie.Host = getEnv("PAYMENT_SERVICE_HOST", "localhost")
	config.PaymentServie.Port = getEnv("PAYMENT_SERVICE_PORT", ":4444")

	// token configuration
	config.Token.Secret = getEnv("TOKEN_SECRET", "token_secret")

	// access ttl parse
	accessTTl, err := time.ParseDuration(getEnv("TOKEN_ACCESS_TTL", "1h"))
	if err != nil {
		return nil, err
	}
	// refresh ttl parse
	refreshTTL, err := time.ParseDuration(getEnv("TOKEN_REFRESH_TTL", "24h"))
	if err != nil {
		return nil, err
	}
	config.Token.AccessTTL = accessTTl
	config.Token.RefreshTTL = refreshTTL
	config.Token.SignInKey = getEnv("TOKEN_SIGNIN_KEY", "abdulazizXoshimov")

	// otlp collector configuration
	config.OTLPCollector.Host = getEnv("OTLP_COLLECTOR_HOST", "otel-collector")
	config.OTLPCollector.Port = getEnv("OTLP_COLLECTOR_PORT", ":4318")

	// kafka configuration
	config.Kafka.Address = strings.Split(getEnv("KAFKA_ADDRESS", "kafka:9091"), ",")
	config.Kafka.Topic.UserCreateTopic = getEnv("KAFKA_USER_CREATE_TOPIC", "api.create.user")

	config.SMTP.Email = getEnv("SMTP_EMAIL", "abdulazizxoshimov22@gmail.com")
	config.SMTP.EmailPassword = getEnv("SMTP_EMAIL_PASSWORD", "hxytgczqprxfsltu ")
	config.SMTP.SMTPPort = getEnv("SMTP_PORT", "587")
	config.SMTP.SMTPHost = getEnv("SMTP_HOST", "smtp.gmail.com")

	return &config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
