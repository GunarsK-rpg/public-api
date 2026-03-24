package middleware

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/GunarsK-rpg/public-api/internal/cache"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestUserSync_NoAuth_Returns401(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	appCache := cache.New(client)

	router := gin.New()
	router.Use(UserSync(nil, appCache))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestUserSync_CacheHit_SkipsDB(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	appCache := cache.New(client)

	// Compute the exact cache key the middleware would generate for user_id=1, username="testuser", displayName=""
	claimsHash := fmt.Sprintf("%x", sha256.Sum256([]byte("testuser\x00")))[:12]
	cacheKey := fmt.Sprintf("%s%d:%s", userSyncKeyPrefix, int64(1), claimsHash)

	if err := appCache.SetFlag(context.Background(), cacheKey, userSyncTTL); err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Set("username", "testuser")
		c.Request.RemoteAddr = "1.2.3.4:1234"
		c.Next()
	})
	// Pass nil pool — if cache hit works, pool is never touched
	router.Use(UserSync(nil, appCache))
	router.GET("/test", func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if !handlerCalled {
		t.Error("expected handler to be called on cache hit")
	}
}
