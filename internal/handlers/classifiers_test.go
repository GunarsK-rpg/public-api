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
	"github.com/jackc/pgx/v5"
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
		getClassifiersFilteredFunc: func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
			callCount++
			return json.RawMessage(`{"fresh":"data"}`), nil
		},
	}
	handler := New(mock, c)

	// Pre-populate cache with global data
	ctx := context.Background()
	cached := json.RawMessage(`{"cached":"data"}`)
	if err := c.Set(ctx, classifiersCacheGlobalKey, cached, 0); err != nil {
		t.Fatal(err)
	}

	router := setupClassifierRouter(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	// Response wraps cached global in scoped shape
	var result struct {
		Global json.RawMessage `json:"global"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if string(result.Global) != string(cached) {
		t.Fatalf("expected global %s, got %s", cached, result.Global)
	}
	if callCount != 0 {
		t.Fatalf("expected 0 DB calls on cache hit, got %d", callCount)
	}
}

func TestGetAllClassifiers_CacheMiss(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	c := cache.New(client)

	dbData := json.RawMessage(`{"from":"database"}`)
	callCount := 0
	mock := &mockRepo{
		getClassifiersFilteredFunc: func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
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
	if callCount != 1 {
		t.Fatalf("expected 1 DB call on cache miss, got %d", callCount)
	}

	// Verify global data was stored in cache
	result, err := c.Get(context.Background(), classifiersCacheGlobalKey)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != string(dbData) {
		t.Fatalf("expected cache to contain %s, got %s", dbData, result)
	}

	// Verify TTL was set
	ttl := mr.TTL(classifiersCacheGlobalKey)
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
		getClassifiersFilteredFunc: func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
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

	// Second call: cache hit, skips DB
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/classifiers", nil)
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("second call: expected 200, got %d", w2.Code)
	}
	if callCount != 1 {
		t.Fatalf("expected 1 DB call total (second should hit cache), got %d", callCount)
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
		getClassifiersFilteredFunc: func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
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
	// Verify global data is in the response
	var result struct {
		Global json.RawMessage `json:"global"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if string(result.Global) != string(dbData) {
		t.Fatalf("expected global %s, got %s", dbData, result.Global)
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

// =============================================================================
// GetAllClassifiers with scoped query params
// =============================================================================

func TestGetAllClassifiers_WithCampaignID(t *testing.T) {
	globalData := json.RawMessage(`{"skills":[],"attrs":[1]}`)
	sbData := json.RawMessage(`{"skills":[1],"attrs":[]}`)
	mock := &mockRepo{
		getCampaignSourceBookIDsFunc: func(_ context.Context, _ repository.AuthContext, id int64) ([]int64, error) {
			if id != 5 {
				t.Errorf("campaignID = %d, want 5", id)
			}
			return []int64{10}, nil
		},
		getClassifiersFilteredFunc: func(_ context.Context, _ repository.AuthContext, filter json.RawMessage) (json.RawMessage, error) {
			if string(filter) == `{"sourceBookId": null}` {
				return globalData, nil
			}
			return sbData, nil
		},
	}
	handler := New(mock, nil)
	router := setupClassifierRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers?campaignId=5", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result struct {
		Global      json.RawMessage   `json:"global"`
		SourceBooks []json.RawMessage `json:"sourceBooks"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if string(result.Global) != string(globalData) {
		t.Fatalf("global = %s, want %s", result.Global, globalData)
	}
	if len(result.SourceBooks) != 1 || string(result.SourceBooks[0]) != string(sbData) {
		t.Fatalf("sourceBooks = %v, want [%s]", result.SourceBooks, sbData)
	}
}

func TestGetAllClassifiers_WithHeroID(t *testing.T) {
	globalData := json.RawMessage(`{"skills":[1]}`)
	heroData := json.RawMessage(`{"skills":[2]}`)
	mock := &mockRepo{
		getClassifiersFilteredFunc: func(_ context.Context, _ repository.AuthContext, filter json.RawMessage) (json.RawMessage, error) {
			if string(filter) == `{"sourceBookId": null}` {
				return globalData, nil
			}
			if string(filter) == `{"heroId": 7}` {
				return heroData, nil
			}
			return json.RawMessage(`{}`), nil
		},
	}
	handler := New(mock, nil)
	router := setupClassifierRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers?heroId=7", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result struct {
		Global json.RawMessage `json:"global"`
		Hero   json.RawMessage `json:"hero"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if string(result.Hero) != string(heroData) {
		t.Fatalf("hero = %s, want %s", result.Hero, heroData)
	}
}

func TestGetAllClassifiers_NoSourceBooks(t *testing.T) {
	globalData := json.RawMessage(`{"skills":[]}`)
	mock := &mockRepo{
		getCampaignSourceBookIDsFunc: func(_ context.Context, _ repository.AuthContext, _ int64) ([]int64, error) {
			return nil, nil
		},
		getClassifiersFilteredFunc: func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
			return globalData, nil
		},
	}
	handler := New(mock, nil)
	router := setupClassifierRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers?campaignId=5", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result struct {
		SourceBooks []json.RawMessage `json:"sourceBooks"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(result.SourceBooks) != 0 {
		t.Fatalf("expected empty sourceBooks, got %d", len(result.SourceBooks))
	}
}

func TestGetAllClassifiers_WithSourceBookID_AccessibleSuccess(t *testing.T) {
	globalData := json.RawMessage(`{"skills":[]}`)
	sbData := json.RawMessage(`{"ancestries":[{"id":1,"name":"Elf"}]}`)
	accessCalls := 0
	mock := &mockRepo{
		requireSourceBookAccessibleFunc: func(_ context.Context, _ repository.AuthContext, id int64) error {
			accessCalls++
			if id != 42 {
				t.Errorf("sourceBookID = %d, want 42", id)
			}
			return nil
		},
		getClassifiersFilteredFunc: func(_ context.Context, _ repository.AuthContext, filter json.RawMessage) (json.RawMessage, error) {
			if string(filter) == `{"sourceBookId": null}` {
				return globalData, nil
			}
			if string(filter) == `{"sourceBookId": 42}` {
				return sbData, nil
			}
			return json.RawMessage(`{}`), nil
		},
	}
	handler := New(mock, nil)
	router := setupClassifierRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers?sourceBookId=42", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if accessCalls != 1 {
		t.Fatalf("expected 1 access check, got %d", accessCalls)
	}

	var result struct {
		Global      json.RawMessage   `json:"global"`
		SourceBooks []json.RawMessage `json:"sourceBooks"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(result.SourceBooks) != 1 || string(result.SourceBooks[0]) != string(sbData) {
		t.Fatalf("sourceBooks = %v, want [%s]", result.SourceBooks, sbData)
	}
}

func TestGetAllClassifiers_WithSourceBookID_NotAccessible(t *testing.T) {
	mock := &mockRepo{
		requireSourceBookAccessibleFunc: func(_ context.Context, _ repository.AuthContext, _ int64) error {
			return pgx.ErrNoRows
		},
		getClassifiersFilteredFunc: func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`{}`), nil
		},
	}
	handler := New(mock, nil)
	router := setupClassifierRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers?sourceBookId=42", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 on inaccessible book, got %d", w.Code)
	}
}

func TestGetAllClassifiers_MutuallyExclusiveScopes(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := setupClassifierRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers?campaignId=1&sourceBookId=2", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for conflicting scopes, got %d", w.Code)
	}
}

func TestGetAllClassifiers_CampaignNotFound(t *testing.T) {
	mock := &mockRepo{
		getClassifiersFilteredFunc: func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`{}`), nil
		},
		getCampaignSourceBookIDsFunc: func(_ context.Context, _ repository.AuthContext, _ int64) ([]int64, error) {
			return nil, pgx.ErrNoRows
		},
	}
	handler := New(mock, nil)
	router := setupClassifierRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classifiers?campaignId=999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
