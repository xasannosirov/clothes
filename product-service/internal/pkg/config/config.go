package config

import (
	"os"
)

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	RPCPort     string

	Context struct {
		Timeout string
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SslMode  string
	}

	UserService struct {
		Host string
		Port string
	}

	MediaService struct {
		Host string
		Port string
	}
}

func New() *Config {
	var config Config

	// general configuration
	config.APP = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.RPCPort = getEnv("RPC_PORT", ":3333")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "clothes-postgres") // postgres
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "root") // root
	config.DB.SslMode = getEnv("POSTGRES_SSLMODE", "disable")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "clothes_store")

	config.UserService.Host = getEnv("USER_SERVICE_RPC_HOST", "clothes-user-service") // user-service
	config.UserService.Port = getEnv("USER_SERVICE_RPC_PORT", ":1111")
	config.MediaService.Host = getEnv("MEDIA_SERVICE_RPC_HOST", "clothes-media-service") // media-service
	config.MediaService.Port = getEnv("MEDIA_SERVICE_RPC_PORT", ":2222")

	return &config
}

func getEnv(key string, defaultVaule string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultVaule
}
