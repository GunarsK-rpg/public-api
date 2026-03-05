package cache

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"

	common "github.com/GunarsK-portfolio/portfolio-common/config"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client instance.
func NewRedisClient(cfg common.RedisConfig, environment string) *redis.Client {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, strconv.Itoa(cfg.Port)),
		Password: cfg.Password,
		DB:       0,
	}

	if environment == "production" {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: cfg.Host,
		}
	}

	client := redis.NewClient(options)

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	return client
}
