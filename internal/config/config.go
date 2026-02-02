package config

import (
	common "github.com/GunarsK-portfolio/portfolio-common/config"
)

// Config holds all configuration for the service.
type Config struct {
	Service  common.ServiceConfig
	Database common.DatabaseConfig
	JWT      JWTConfig
}

// JWTConfig holds JWT authentication configuration.
type JWTConfig struct {
	Secret string
}

// Load loads all configuration from environment variables.
func Load() *Config {
	return &Config{
		Service:  common.NewServiceConfig(8090),
		Database: common.NewDatabaseConfig(),
		JWT: JWTConfig{
			Secret: common.GetEnvRequired("JWT_SECRET"),
		},
	}
}
