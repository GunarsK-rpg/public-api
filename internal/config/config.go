package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"

	common "github.com/GunarsK-portfolio/portfolio-common/config"
)

// Config holds all configuration for the service.
type Config struct {
	Service     common.ServiceConfig
	Database    common.DatabaseConfig
	Redis       common.RedisConfig
	JWTSecret   string `validate:"required,min=32"`
	MaxBodySize int64  `validate:"min=1024"` // Minimum 1KB
}

// Load loads all configuration from environment variables.
func Load() *Config {
	cfg := &Config{
		Service:     common.NewServiceConfig(8090),
		Database:    common.NewDatabaseConfig(),
		Redis:       common.NewRedisConfig(),
		JWTSecret:   common.GetEnvRequired("JWT_SECRET"),
		MaxBodySize: getEnvInt64("MAX_BODY_SIZE", 64<<10), // Default 64KB
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic(fmt.Sprintf("Invalid configuration: %v", err))
	}

	return cfg
}

// getEnvInt64 returns the int64 value of an environment variable or the default.
func getEnvInt64(key string, defaultVal int64) int64 {
	if val := os.Getenv(key); val != "" {
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			slog.Warn("invalid environment variable value, using default",
				"key", key,
				"value", val,
				"default", defaultVal,
			)
			return defaultVal
		}
		return i
	}
	return defaultVal
}
