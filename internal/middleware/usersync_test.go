package middleware

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/GunarsK-rpg/public-api/internal/cache"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

type mockDBExecer struct {
	execFn func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func (m *mockDBExecer) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return m.execFn(ctx, sql, arguments...)
}

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

func TestUserSync_CacheMiss_SyncsDB(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	appCache := cache.New(client)

	execCalled := false
	mock := &mockDBExecer{
		execFn: func(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			execCalled = true
			expectedSQL := "SELECT auth.sync_user($1, $2, $3)"
			if sql != expectedSQL {
				t.Errorf("sql = %q, want %q", sql, expectedSQL)
			}
			if len(args) != 3 {
				t.Fatalf("args len = %d, want 3", len(args))
			}
			if args[0] != int64(1) {
				t.Errorf("args[0] = %v, want int64(1)", args[0])
			}
			if args[1] != "testuser" {
				t.Errorf("args[1] = %v, want %q", args[1], "testuser")
			}
			if args[2] != (*string)(nil) {
				t.Errorf("args[2] = %v, want nil *string (no display name)", args[2])
			}
			return pgconn.NewCommandTag("SELECT 1"), nil
		},
	}

	handlerCalled := false
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Set("username", "testuser")
		c.Request.RemoteAddr = "1.2.3.4:1234"
		c.Next()
	})
	router.Use(UserSync(mock, appCache))
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
	if !execCalled {
		t.Error("expected DB exec to be called on cache miss")
	}
	if !handlerCalled {
		t.Error("expected handler to be called after DB sync")
	}

	// Verify cache was set after successful sync
	claimsHash := fmt.Sprintf("%x", sha256.Sum256([]byte("testuser\x00")))[:12]
	cacheKey := fmt.Sprintf("%s%d:%s", userSyncKeyPrefix, int64(1), claimsHash)
	synced, err := appCache.HasFlag(context.Background(), cacheKey)
	if err != nil {
		t.Fatal(err)
	}
	if !synced {
		t.Error("expected cache flag to be set after DB sync")
	}
}

func TestUserSync_DBFailure_Returns500(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	appCache := cache.New(client)

	mock := &mockDBExecer{
		execFn: func(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
			return pgconn.CommandTag{}, errors.New("connection refused")
		},
	}

	handlerCalled := false
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Set("username", "testuser")
		c.Request.RemoteAddr = "1.2.3.4:1234"
		c.Next()
	})
	router.Use(UserSync(mock, appCache))
	router.GET("/test", func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
	if handlerCalled {
		t.Error("handler should not be called on DB failure")
	}

	// Verify cache flag was NOT set after DB failure
	claimsHash := fmt.Sprintf("%x", sha256.Sum256([]byte("testuser\x00")))[:12]
	cacheKey := fmt.Sprintf("%s%d:%s", userSyncKeyPrefix, int64(1), claimsHash)
	synced, err := appCache.HasFlag(context.Background(), cacheKey)
	if err != nil {
		t.Fatal(err)
	}
	if synced {
		t.Error("cache flag should not be set after DB failure")
	}
}
