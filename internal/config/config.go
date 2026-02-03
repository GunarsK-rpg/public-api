package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	common "github.com/GunarsK-portfolio/portfolio-common/config"
)

// Config holds all configuration for the service.
type Config struct {
	Service   common.ServiceConfig
	Database  common.DatabaseConfig
	JWTSecret string `validate:"required,min=32"`
}

// Load loads all configuration from environment variables.
func Load() *Config {
	cfg := &Config{
		Service:   common.NewServiceConfig(8090),
		Database:  common.NewDatabaseConfig(),
		JWTSecret: common.GetEnvRequired("JWT_SECRET"),
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic(fmt.Sprintf("Invalid configuration: %v", err))
	}

	return cfg
}
