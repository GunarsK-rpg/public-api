package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache provides get/set operations backed by Redis.
type Cache struct {
	client *redis.Client
}

// New creates a new Cache instance.
func New(client *redis.Client) *Cache {
	return &Cache{client: client}
}

// Get returns the cached value for key, or nil on miss.
func (c *Cache) Get(ctx context.Context, key string) (json.RawMessage, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Set stores data under key with the given TTL.
func (c *Cache) Set(ctx context.Context, key string, data json.RawMessage, ttl time.Duration) error {
	return c.client.Set(ctx, key, []byte(data), ttl).Err()
}
