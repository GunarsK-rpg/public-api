package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/GunarsK-rpg/public-api/internal/repository"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// =============================================================================
// Test Helpers
// =============================================================================

// withAuth adds valid auth context to the Gin context.
func withAuth(c *gin.Context) {
	c.Set("user_id", int64(1))
	c.Set("username", "testuser")
}

// performRequest executes an HTTP request against a Gin router.
func performRequest(t *testing.T, r *gin.Engine, method, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if body != nil {
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(body))
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// =============================================================================
// GetAuthContext Tests
// =============================================================================

func TestGetAuthContext_Success(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", int64(42))
	c.Set("username", "kaladin")
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.RemoteAddr = "1.2.3.4:1234"
	c.Request.Header.Set("User-Agent", "test-agent")

	auth, err := GetAuthContext(c)
	if err != nil {
		t.Fatalf("GetAuthContext() error = %v", err)
	}
	if auth.UserID != 42 {
		t.Errorf("UserID = %d, want 42", auth.UserID)
	}
	if auth.Username != "kaladin" {
		t.Errorf("Username = %q, want %q", auth.Username, "kaladin")
	}
	if auth.ClientIP != "1.2.3.4" {
		t.Errorf("ClientIP = %q, want %q", auth.ClientIP, "1.2.3.4")
	}
	if auth.UserAgent != "test-agent" {
		t.Errorf("UserAgent = %q, want %q", auth.UserAgent, "test-agent")
	}
}

func TestGetAuthContext_MissingUserID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("username", "testuser")

	_, err := GetAuthContext(c)
	if !errors.Is(err, ErrMissingAuthContext) {
		t.Errorf("error = %v, want %v", err, ErrMissingAuthContext)
	}
}

func TestGetAuthContext_MissingUsername(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", int64(1))

	_, err := GetAuthContext(c)
	if !errors.Is(err, ErrMissingAuthContext) {
		t.Errorf("error = %v, want %v", err, ErrMissingAuthContext)
	}
}

func TestGetAuthContext_WrongUserIDType(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", "not-an-int64")
	c.Set("username", "testuser")

	_, err := GetAuthContext(c)
	if !errors.Is(err, ErrMissingAuthContext) {
		t.Errorf("error = %v, want %v", err, ErrMissingAuthContext)
	}
}

func TestGetAuthContext_WrongUsernameType(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", int64(1))
	c.Set("username", 12345)

	_, err := GetAuthContext(c)
	if !errors.Is(err, ErrMissingAuthContext) {
		t.Errorf("error = %v, want %v", err, ErrMissingAuthContext)
	}
}

// =============================================================================
// HandlePgxError Tests
// =============================================================================

// errorResponse matches the JSON shape from commonHandlers.RespondError/LogAndRespondError.
type errorResponse struct {
	Error string `json:"error"`
}

func TestHandlePgxError_NoRows(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	HandlePgxError(c, pgx.ErrNoRows)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
	var resp errorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp.Error != "not found" {
		t.Errorf("message = %q, want %q", resp.Error, "not found")
	}
}

func TestHandlePgxError_PgCodes(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		status  int
		message string
	}{
		{"unique_violation", "23505", http.StatusConflict, "resource already exists"},
		{"foreign_key_violation", "23503", http.StatusBadRequest, "referenced resource not found"},
		{"check_violation", "23514", http.StatusBadRequest, "validation constraint failed"},
		{"no_data_found", "P0002", http.StatusNotFound, "not found"},
		{"insufficient_privilege", "42501", http.StatusForbidden, "access denied"},
		{"invalid_parameter_value", "22023", http.StatusBadRequest, "invalid parameter value"},
		{"string_data_right_truncation", "22001", http.StatusBadRequest, "value too long"},
		{"raise_exception", "P0001", http.StatusBadRequest, "test error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/", nil)

			pgErr := &pgconn.PgError{Code: tt.code, Message: "test error"}
			HandlePgxError(c, pgErr)

			if w.Code != tt.status {
				t.Errorf("status = %d, want %d", w.Code, tt.status)
			}
			var resp errorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}
			if resp.Error != tt.message {
				t.Errorf("message = %q, want %q", resp.Error, tt.message)
			}
		})
	}
}

func TestHandlePgxError_UnknownPgCode(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	pgErr := &pgconn.PgError{Code: "99999", Message: "unknown"}
	HandlePgxError(c, pgErr)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
	var resp errorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp.Error != "internal server error" {
		t.Errorf("message = %q, want %q", resp.Error, "internal server error")
	}
}

func TestHandlePgxError_GenericError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	HandlePgxError(c, errors.New("something broke"))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
	var resp errorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp.Error != "internal server error" {
		t.Errorf("message = %q, want %q", resp.Error, "internal server error")
	}
}

// =============================================================================
// handleGet Tests
// =============================================================================

func TestHandleGet_Success(t *testing.T) {
	router := gin.New()
	expected := json.RawMessage(`{"key":"value"}`)
	router.GET("/test", func(c *gin.Context) {
		withAuth(c)
		handleGet(c, func(_ context.Context, _ repository.AuthContext) (json.RawMessage, error) {
			return expected, nil
		})
	})

	w := performRequest(t, router, "GET", "/test", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != string(expected) {
		t.Errorf("body = %s, want %s", w.Body.String(), expected)
	}
}

func TestHandleGet_NoAuth(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handleGet(c, func(_ context.Context, _ repository.AuthContext) (json.RawMessage, error) {
			t.Fatal("repo should not be called without auth")
			return nil, nil
		})
	})

	w := performRequest(t, router, "GET", "/test", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestHandleGet_RepoError(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		withAuth(c)
		handleGet(c, func(_ context.Context, _ repository.AuthContext) (json.RawMessage, error) {
			return nil, errors.New("db down")
		})
	})

	w := performRequest(t, router, "GET", "/test", nil)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// =============================================================================
// handleGetByID Tests
// =============================================================================

func TestHandleGetByID_Success(t *testing.T) {
	router := gin.New()
	expected := json.RawMessage(`{"id":1}`)
	router.GET("/items/:id", func(c *gin.Context) {
		withAuth(c)
		handleGetByID(c, "id", func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
			if id != 42 {
				t.Errorf("id = %d, want 42", id)
			}
			return expected, nil
		})
	})

	w := performRequest(t, router, "GET", "/items/42", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != string(expected) {
		t.Errorf("body = %s, want %s", w.Body.String(), expected)
	}
}

func TestHandleGetByID_NoAuth(t *testing.T) {
	router := gin.New()
	router.GET("/items/:id", func(c *gin.Context) {
		handleGetByID(c, "id", func(_ context.Context, _ repository.AuthContext, _ int64) (json.RawMessage, error) {
			t.Fatal("repo should not be called without auth")
			return nil, nil
		})
	})

	w := performRequest(t, router, "GET", "/items/1", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestHandleGetByID_InvalidID(t *testing.T) {
	router := gin.New()
	router.GET("/items/:id", func(c *gin.Context) {
		withAuth(c)
		handleGetByID(c, "id", func(_ context.Context, _ repository.AuthContext, _ int64) (json.RawMessage, error) {
			t.Fatal("repo should not be called with invalid ID")
			return nil, nil
		})
	})

	tests := []struct {
		name string
		id   string
	}{
		{"alphabetic", "abc"},
		{"float", "1.5"},
		{"special chars", "!@#"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performRequest(t, router, "GET", "/items/"+tt.id, nil)
			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestHandleGetByID_NullResult(t *testing.T) {
	router := gin.New()
	router.GET("/items/:id", func(c *gin.Context) {
		withAuth(c)
		handleGetByID(c, "id", func(_ context.Context, _ repository.AuthContext, _ int64) (json.RawMessage, error) {
			return json.RawMessage("null"), nil
		})
	})

	w := performRequest(t, router, "GET", "/items/1", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandleGetByID_NilResult(t *testing.T) {
	router := gin.New()
	router.GET("/items/:id", func(c *gin.Context) {
		withAuth(c)
		handleGetByID(c, "id", func(_ context.Context, _ repository.AuthContext, _ int64) (json.RawMessage, error) {
			return nil, nil
		})
	})

	w := performRequest(t, router, "GET", "/items/1", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandleGetByID_RepoError(t *testing.T) {
	router := gin.New()
	router.GET("/items/:id", func(c *gin.Context) {
		withAuth(c)
		handleGetByID(c, "id", func(_ context.Context, _ repository.AuthContext, _ int64) (json.RawMessage, error) {
			return nil, pgx.ErrNoRows
		})
	})

	w := performRequest(t, router, "GET", "/items/1", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// =============================================================================
// handleGetByString Tests
// =============================================================================

func TestHandleGetByString_Success(t *testing.T) {
	router := gin.New()
	expected := json.RawMessage(`{"code":"warrior"}`)
	router.GET("/items/:code", func(c *gin.Context) {
		withAuth(c)
		handleGetByString(c, "code", func(_ context.Context, _ repository.AuthContext, code string) (json.RawMessage, error) {
			if code != "warrior" {
				t.Errorf("code = %q, want %q", code, "warrior")
			}
			return expected, nil
		})
	})

	w := performRequest(t, router, "GET", "/items/warrior", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHandleGetByString_NoAuth(t *testing.T) {
	router := gin.New()
	router.GET("/items/:code", func(c *gin.Context) {
		handleGetByString(c, "code", func(_ context.Context, _ repository.AuthContext, _ string) (json.RawMessage, error) {
			t.Fatal("repo should not be called without auth")
			return nil, nil
		})
	})

	w := performRequest(t, router, "GET", "/items/test", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestHandleGetByString_NullResult(t *testing.T) {
	router := gin.New()
	router.GET("/items/:code", func(c *gin.Context) {
		withAuth(c)
		handleGetByString(c, "code", func(_ context.Context, _ repository.AuthContext, _ string) (json.RawMessage, error) {
			return json.RawMessage("null"), nil
		})
	})

	w := performRequest(t, router, "GET", "/items/nonexistent", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandleGetByString_RepoError(t *testing.T) {
	router := gin.New()
	router.GET("/items/:code", func(c *gin.Context) {
		withAuth(c)
		handleGetByString(c, "code", func(_ context.Context, _ repository.AuthContext, _ string) (json.RawMessage, error) {
			return nil, errors.New("db error")
		})
	})

	w := performRequest(t, router, "GET", "/items/test", nil)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// =============================================================================
// handlePost Tests
// =============================================================================

func TestHandlePost_Success(t *testing.T) {
	router := gin.New()
	expected := json.RawMessage(`{"id":1}`)
	router.POST("/items", func(c *gin.Context) {
		withAuth(c)
		handlePost(c, func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
			if string(data) != `{"name":"test"}` {
				t.Errorf("data = %s, want %s", data, `{"name":"test"}`)
			}
			return expected, nil
		})
	})

	w := performRequest(t, router, "POST", "/items", []byte(`{"name":"test"}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != string(expected) {
		t.Errorf("body = %s, want %s", w.Body.String(), expected)
	}
}

func TestHandlePost_NoAuth(t *testing.T) {
	router := gin.New()
	router.POST("/items", func(c *gin.Context) {
		handlePost(c, func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
			t.Fatal("repo should not be called without auth")
			return nil, nil
		})
	})

	w := performRequest(t, router, "POST", "/items", []byte(`{}`))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestHandlePost_InvalidJSON(t *testing.T) {
	router := gin.New()
	router.POST("/items", func(c *gin.Context) {
		withAuth(c)
		handlePost(c, func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
			t.Fatal("repo should not be called with invalid JSON")
			return nil, nil
		})
	})

	tests := []struct {
		name string
		body string
	}{
		{"broken json", `{"name": broken`},
		{"plain text", `hello world`},
		{"empty body", ``},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performRequest(t, router, "POST", "/items", []byte(tt.body))
			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestHandlePost_RepoError(t *testing.T) {
	router := gin.New()
	router.POST("/items", func(c *gin.Context) {
		withAuth(c)
		handlePost(c, func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
			return nil, &pgconn.PgError{Code: "23505", Message: "duplicate"}
		})
	})

	w := performRequest(t, router, "POST", "/items", []byte(`{"name":"test"}`))

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
	}
}

// =============================================================================
// handleDelete Tests
// =============================================================================

func TestHandleDelete_Success(t *testing.T) {
	router := gin.New()
	router.DELETE("/items/:id", func(c *gin.Context) {
		withAuth(c)
		handleDelete(c, "id", func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
			if id != 42 {
				t.Errorf("id = %d, want 42", id)
			}
			return true, nil
		})
	})

	w := performRequest(t, router, "DELETE", "/items/42", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestHandleDelete_NoAuth(t *testing.T) {
	router := gin.New()
	router.DELETE("/items/:id", func(c *gin.Context) {
		handleDelete(c, "id", func(_ context.Context, _ repository.AuthContext, _ int64) (bool, error) {
			t.Fatal("repo should not be called without auth")
			return false, nil
		})
	})

	w := performRequest(t, router, "DELETE", "/items/1", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestHandleDelete_InvalidID(t *testing.T) {
	router := gin.New()
	router.DELETE("/items/:id", func(c *gin.Context) {
		withAuth(c)
		handleDelete(c, "id", func(_ context.Context, _ repository.AuthContext, _ int64) (bool, error) {
			t.Fatal("repo should not be called with invalid ID")
			return false, nil
		})
	})

	w := performRequest(t, router, "DELETE", "/items/abc", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleDelete_NotFound(t *testing.T) {
	router := gin.New()
	router.DELETE("/items/:id", func(c *gin.Context) {
		withAuth(c)
		handleDelete(c, "id", func(_ context.Context, _ repository.AuthContext, _ int64) (bool, error) {
			return false, nil
		})
	})

	w := performRequest(t, router, "DELETE", "/items/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandleDelete_RepoError(t *testing.T) {
	router := gin.New()
	router.DELETE("/items/:id", func(c *gin.Context) {
		withAuth(c)
		handleDelete(c, "id", func(_ context.Context, _ repository.AuthContext, _ int64) (bool, error) {
			return false, errors.New("db error")
		})
	})

	w := performRequest(t, router, "DELETE", "/items/1", nil)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// =============================================================================
// getPathParamInt64 Tests
// =============================================================================

func TestGetPathParamInt64_Valid(t *testing.T) {
	router := gin.New()
	router.GET("/items/:id", func(c *gin.Context) {
		id, err := getPathParamInt64(c, "id")
		if err != nil {
			t.Fatalf("getPathParamInt64() error = %v", err)
		}
		if id != 42 {
			t.Errorf("id = %d, want 42", id)
		}
		c.Status(http.StatusOK)
	})

	w := performRequest(t, router, "GET", "/items/42", nil)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetPathParamInt64_NonInteger(t *testing.T) {
	router := gin.New()
	router.GET("/items/:id", func(c *gin.Context) {
		_, err := getPathParamInt64(c, "id")
		if err == nil {
			t.Fatal("getPathParamInt64() should return error for non-integer")
		}
		c.Status(http.StatusBadRequest)
	})

	w := performRequest(t, router, "GET", "/items/abc", nil)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
