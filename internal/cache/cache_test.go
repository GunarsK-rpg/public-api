package cache

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func setupTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	return client, mr
}

func TestCache_Get_Miss(t *testing.T) {
	client, _ := setupTestRedis(t)
	c := New(client)

	result, err := c.Get(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("expected no error on cache miss, got %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil on cache miss, got %s", result)
	}
}

func TestCache_Set_Then_Get(t *testing.T) {
	client, _ := setupTestRedis(t)
	c := New(client)
	ctx := context.Background()

	data := json.RawMessage(`{"attributes":[{"id":1,"name":"Strength"}]}`)

	err := c.Set(ctx, "rpg:classifiers:all", data, 1*time.Hour)
	if err != nil {
		t.Fatalf("expected no error on set, got %v", err)
	}

	result, err := c.Get(ctx, "rpg:classifiers:all")
	if err != nil {
		t.Fatalf("expected no error on get, got %v", err)
	}
	if string(result) != string(data) {
		t.Fatalf("expected %s, got %s", data, result)
	}
}

func TestCache_Set_Respects_TTL(t *testing.T) {
	client, mr := setupTestRedis(t)
	c := New(client)
	ctx := context.Background()

	data := json.RawMessage(`{"test":true}`)
	err := c.Set(ctx, "rpg:test", data, 1*time.Hour)
	if err != nil {
		t.Fatalf("expected no error on set, got %v", err)
	}

	// Fast-forward past TTL
	mr.FastForward(2 * time.Hour)

	result, err := c.Get(ctx, "rpg:test")
	if err != nil {
		t.Fatalf("expected no error after expiry, got %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil after TTL expiry, got %s", result)
	}
}
