package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/GunarsK-rpg/public-api/internal/cache"
	"github.com/GunarsK-rpg/public-api/internal/repository"
)

func setupClassifierRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/classifiers", func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Set("username", "testuser")
		c.Next()
	}, handler.GetAllClassifiers)
	r.GET("/classifiers/source-books", func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Set("username", "testuser")
		c.Next()
	}, handler.GetSourceBooks)
	return r
}

func setupClassifierRouterNoAuth(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/classifiers/source-books", handler.GetSourceBooks)
	return r
}

func TestGetAllClassifiers_CacheHit(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	c := cache.New(client)

	callCount := 0
	mock := &mockRepo{
		getAllClassifiersFunc: func(_ context.Context, _ repository.AuthContext) (json.RawMessage, error) {
			callCount++
			return json.RawMessage(`{"fresh":"data"}`), nil
		},
	}
	handler := New(mock, c)

	// Pre-populate cache
	ctx := context.Background()
	cached := json.RawMessage(`{"cached":"data"}`)
	if err := c.Set(ctx, classifiersCacheKey, cached, 0); err != nil {
		t.Fatal(err)
	}

	router := setupClassifierRouter(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != string(cached) {
		t.Fatalf("expected cached data %s, got %s", cached, w.Body.String())
	}
	if callCount != 0 {
		t.Fatalf("expected 0 DB calls on cache hit, got %d", callCount)
	}
	if w.Header().Get("Cache-Control") != "private, max-age=3600" {
		t.Fatalf("expected Cache-Control header, got %q", w.Header().Get("Cache-Control"))
	}
}

func TestGetAllClassifiers_CacheMiss(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	c := cache.New(client)

	dbData := json.RawMessage(`{"from":"database"}`)
	callCount := 0
	mock := &mockRepo{
		getAllClassifiersFunc: func(_ context.Context, _ repository.AuthContext) (json.RawMessage, error) {
			callCount++
			return dbData, nil
		},
	}
	handler := New(mock, c)

	router := setupClassifierRouter(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != string(dbData) {
		t.Fatalf("expected DB data %s, got %s", dbData, w.Body.String())
	}
	if callCount != 1 {
		t.Fatalf("expected 1 DB call on cache miss, got %d", callCount)
	}

	// Verify data was stored in cache
	result, err := c.Get(context.Background(), classifiersCacheKey)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != string(dbData) {
		t.Fatalf("expected cache to contain %s, got %s", dbData, result)
	}

	// Verify TTL was set
	ttl := mr.TTL(classifiersCacheKey)
	if ttl < 55*time.Minute || ttl > 61*time.Minute {
		t.Fatalf("expected TTL ~1h, got %v", ttl)
	}
}

func TestGetAllClassifiers_SecondCallHitsCache(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	c := cache.New(client)

	dbData := json.RawMessage(`{"from":"database"}`)
	callCount := 0
	mock := &mockRepo{
		getAllClassifiersFunc: func(_ context.Context, _ repository.AuthContext) (json.RawMessage, error) {
			callCount++
			return dbData, nil
		},
	}
	handler := New(mock, c)

	router := setupClassifierRouter(handler)

	// First call: cache miss, hits DB
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/classifiers", nil)
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("first call: expected 200, got %d", w1.Code)
	}
	if w1.Header().Get("Cache-Control") != "private, max-age=3600" {
		t.Fatalf("first call: expected Cache-Control header, got %q", w1.Header().Get("Cache-Control"))
	}

	// Second call: cache hit, skips DB
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/classifiers", nil)
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("second call: expected 200, got %d", w2.Code)
	}
	if w2.Header().Get("Cache-Control") != "private, max-age=3600" {
		t.Fatalf("second call: expected Cache-Control header, got %q", w2.Header().Get("Cache-Control"))
	}
	if callCount != 1 {
		t.Fatalf("expected 1 DB call total (second should hit cache), got %d", callCount)
	}
	if w2.Body.String() != string(dbData) {
		t.Fatalf("expected same data on second call, got %s", w2.Body.String())
	}
}

func TestGetAllClassifiers_RedisFallbackToDB(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	c := cache.New(client)

	// Kill Redis to simulate failure
	mr.Close()

	dbData := json.RawMessage(`{"from":"database"}`)
	callCount := 0
	mock := &mockRepo{
		getAllClassifiersFunc: func(_ context.Context, _ repository.AuthContext) (json.RawMessage, error) {
			callCount++
			return dbData, nil
		},
	}
	handler := New(mock, c)

	router := setupClassifierRouter(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on Redis failure, got %d", w.Code)
	}
	if w.Body.String() != string(dbData) {
		t.Fatalf("expected DB data on Redis failure, got %s", w.Body.String())
	}
	if callCount != 1 {
		t.Fatalf("expected 1 DB call on Redis failure, got %d", callCount)
	}
}

func TestGetSourceBooks_Success(t *testing.T) {
	dbData := json.RawMessage(`[{"id":1,"name":"Stormlight RPG"}]`)
	mock := &mockRepo{
		getSourceBooksFunc: func(_ context.Context, _ repository.AuthContext) (json.RawMessage, error) {
			return dbData, nil
		},
	}
	handler := New(mock, nil)
	router := setupClassifierRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers/source-books", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != string(dbData) {
		t.Fatalf("expected %s, got %s", dbData, w.Body.String())
	}
}

func TestGetSourceBooks_NoAuth(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := setupClassifierRouterNoAuth(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers/source-books", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestGetSourceBooks_RepoError(t *testing.T) {
	mock := &mockRepo{
		getSourceBooksFunc: func(_ context.Context, _ repository.AuthContext) (json.RawMessage, error) {
			return nil, errors.New("db error")
		},
	}
	handler := New(mock, nil)
	router := setupClassifierRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers/source-books", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
