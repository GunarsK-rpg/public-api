package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/GunarsK-rpg/public-api/internal/repository"
)

// =============================================================================
// injectScope
// =============================================================================

func TestInjectScope_BookID(t *testing.T) {
	id := int64(42)
	out, err := injectScope(json.RawMessage(`{"name":"Lashings","sourceBookId":999}`), &id, nil)
	if err != nil {
		t.Fatalf("injectScope() error = %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got["sourceBookId"].(float64) != 42 {
		t.Errorf("sourceBookId = %v, want 42", got["sourceBookId"])
	}
	if got["heroId"] != nil {
		t.Errorf("heroId = %v, want nil", got["heroId"])
	}
	if got["name"] != "Lashings" {
		t.Errorf("name = %v, want Lashings", got["name"])
	}
}

func TestInjectScope_HeroID(t *testing.T) {
	id := int64(7)
	out, err := injectScope(json.RawMessage(`{}`), nil, &id)
	if err != nil {
		t.Fatalf("injectScope() error = %v", err)
	}
	var got map[string]any
	_ = json.Unmarshal(out, &got)
	if got["heroId"].(float64) != 7 {
		t.Errorf("heroId = %v, want 7", got["heroId"])
	}
}

func TestInjectScope_BothNil(t *testing.T) {
	if _, err := injectScope(json.RawMessage(`{}`), nil, nil); err == nil {
		t.Error("expected error for both-nil scope")
	}
}

func TestInjectScope_BothSet(t *testing.T) {
	a, b := int64(1), int64(2)
	if _, err := injectScope(json.RawMessage(`{}`), &a, &b); err == nil {
		t.Error("expected error for both-set scope")
	}
}

func TestInjectScope_NonObject(t *testing.T) {
	id := int64(1)
	if _, err := injectScope(json.RawMessage(`[1,2,3]`), &id, nil); err == nil {
		t.Error("expected error for non-object payload")
	}
}

func TestInjectScope_EmptyBody(t *testing.T) {
	id := int64(5)
	out, err := injectScope(nil, &id, nil)
	if err != nil {
		t.Fatalf("injectScope() error = %v", err)
	}
	var got map[string]any
	_ = json.Unmarshal(out, &got)
	if got["sourceBookId"].(float64) != 5 {
		t.Errorf("sourceBookId = %v, want 5", got["sourceBookId"])
	}
}

// =============================================================================
// mergeCode + parseSourceBookID
// =============================================================================

func TestMergeCode_Overrides(t *testing.T) {
	out, err := mergeCode(json.RawMessage(`{"code":"old","name":"X"}`), "new-code")
	if err != nil {
		t.Fatalf("mergeCode() error = %v", err)
	}
	var got map[string]any
	_ = json.Unmarshal(out, &got)
	if got["code"] != "new-code" {
		t.Errorf("code = %v, want new-code", got["code"])
	}
}

func TestParseSourceBookID(t *testing.T) {
	id, ok := parseSourceBookID(json.RawMessage(`{"id":99,"name":"x"}`))
	if !ok || id != 99 {
		t.Errorf("got id=%d ok=%v, want 99 true", id, ok)
	}
	if _, ok := parseSourceBookID(json.RawMessage(`null`)); ok {
		t.Error("expected null payload to return ok=false")
	}
	if _, ok := parseSourceBookID(nil); ok {
		t.Error("expected nil payload to return ok=false")
	}
}

// =============================================================================
// Source-book handlers
// =============================================================================

func TestCreateSourceBook_Success(t *testing.T) {
	mock := &mockRepo{}
	mock.upsertSourceBookFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
		return json.RawMessage(`{"id":17,"name":"My Book"}`), nil
	}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/homebrew/source-books", func(c *gin.Context) {
		withAuth(c)
		handler.CreateSourceBook(c)
	})

	w := performRequest(t, router, "POST", "/homebrew/source-books", []byte(`{"name":"My Book","gameSystemCode":"homebrew"}`))
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d (body=%s)", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestUpdateSourceBook_OverridesCode(t *testing.T) {
	mock := &mockRepo{}
	var captured json.RawMessage
	mock.upsertSourceBookFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
		captured = data
		return json.RawMessage(`{"id":1}`), nil
	}
	handler := New(mock, nil)
	router := gin.New()
	router.PUT("/homebrew/source-books/:code", func(c *gin.Context) {
		withAuth(c)
		handler.UpdateSourceBook(c)
	})

	w := performRequest(t, router, "PUT", "/homebrew/source-books/url-code",
		[]byte(`{"name":"X","code":"client-code"}`))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	var got map[string]any
	_ = json.Unmarshal(captured, &got)
	if got["code"] != "url-code" {
		t.Errorf("captured code = %v, want url-code", got["code"])
	}
}

func TestDeleteSourceBook_Success(t *testing.T) {
	mock := &mockRepo{
		getSourceBookByCodeFunc: func(_ context.Context, _ repository.AuthContext, code string) (json.RawMessage, error) {
			return json.RawMessage(`{"id":3,"code":"abc"}`), nil
		},
		deleteSourceBookByCodeFunc: func(_ context.Context, _ repository.AuthContext, _ string) (bool, error) {
			return true, nil
		},
	}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/homebrew/source-books/:code", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteSourceBook(c)
	})

	w := performRequest(t, router, "DELETE", "/homebrew/source-books/abc", nil)
	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d (body=%s)", w.Code, http.StatusNoContent, w.Body.String())
	}
}

func TestDeleteSourceBook_NotFound(t *testing.T) {
	mock := &mockRepo{
		getSourceBookByCodeFunc: func(_ context.Context, _ repository.AuthContext, _ string) (json.RawMessage, error) {
			return json.RawMessage(`null`), nil
		},
	}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/homebrew/source-books/:code", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteSourceBook(c)
	})

	w := performRequest(t, router, "DELETE", "/homebrew/source-books/missing", nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestRestoreSourceBook_Success(t *testing.T) {
	mock := &mockRepo{
		getSourceBookByCodeFunc: func(_ context.Context, _ repository.AuthContext, _ string) (json.RawMessage, error) {
			return json.RawMessage(`{"id":3}`), nil
		},
		restoreSourceBookByCodeFunc: func(_ context.Context, _ repository.AuthContext, _ string) (bool, error) {
			return true, nil
		},
	}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/homebrew/source-books/:code/restore", func(c *gin.Context) {
		withAuth(c)
		handler.RestoreSourceBook(c)
	})

	w := performRequest(t, router, "POST", "/homebrew/source-books/abc/restore", nil)
	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

// =============================================================================
// Book-scoped classifier handlers
// =============================================================================

func TestUpsertBookClassifier_Success_OverridesScope(t *testing.T) {
	mock := &mockRepo{
		getSourceBookByCodeFunc: func(_ context.Context, _ repository.AuthContext, _ string) (json.RawMessage, error) {
			return json.RawMessage(`{"id":42}`), nil
		},
	}
	var capturedType string
	var capturedData json.RawMessage
	mock.upsertClassifierFunc = func(_ context.Context, _ repository.AuthContext, ct string, data json.RawMessage) (json.RawMessage, error) {
		capturedType = ct
		capturedData = data
		return json.RawMessage(`{"id":7}`), nil
	}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/homebrew/source-books/:code/:type", func(c *gin.Context) {
		withAuth(c)
		handler.UpsertBookClassifier(c)
	})

	body := []byte(`{"name":"X","sourceBookId":999,"heroId":123}`)
	w := performRequest(t, router, "POST", "/homebrew/source-books/abc/talents", body)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	if capturedType != "talents" {
		t.Errorf("type = %q, want talents", capturedType)
	}
	var got map[string]any
	_ = json.Unmarshal(capturedData, &got)
	if got["sourceBookId"].(float64) != 42 {
		t.Errorf("sourceBookId = %v, want 42", got["sourceBookId"])
	}
	if got["heroId"] != nil {
		t.Errorf("heroId = %v, want nil", got["heroId"])
	}
}

func TestUpsertBookClassifier_UnknownType(t *testing.T) {
	mock := &mockRepo{
		getSourceBookByCodeFunc: func(_ context.Context, _ repository.AuthContext, _ string) (json.RawMessage, error) {
			return json.RawMessage(`{"id":1}`), nil
		},
	}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/homebrew/source-books/:code/:type", func(c *gin.Context) {
		withAuth(c)
		handler.UpsertBookClassifier(c)
	})

	w := performRequest(t, router, "POST", "/homebrew/source-books/abc/widgets", []byte(`{}`))
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUpsertBookClassifier_BookNotFound(t *testing.T) {
	mock := &mockRepo{
		getSourceBookByCodeFunc: func(_ context.Context, _ repository.AuthContext, _ string) (json.RawMessage, error) {
			return json.RawMessage(`null`), nil
		},
	}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/homebrew/source-books/:code/:type", func(c *gin.Context) {
		withAuth(c)
		handler.UpsertBookClassifier(c)
	})

	w := performRequest(t, router, "POST", "/homebrew/source-books/abc/talents", []byte(`{}`))
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestUpsertBookClassifier_RepoOwnershipDenied(t *testing.T) {
	mock := &mockRepo{
		getSourceBookByCodeFunc: func(_ context.Context, _ repository.AuthContext, _ string) (json.RawMessage, error) {
			return json.RawMessage(`{"id":1}`), nil
		},
		upsertClassifierFunc: func(_ context.Context, _ repository.AuthContext, _ string, _ json.RawMessage) (json.RawMessage, error) {
			return nil, &pgconn.PgError{Code: "42501", Message: "denied"}
		},
	}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/homebrew/source-books/:code/:type", func(c *gin.Context) {
		withAuth(c)
		handler.UpsertBookClassifier(c)
	})

	w := performRequest(t, router, "POST", "/homebrew/source-books/abc/talents", []byte(`{}`))
	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestDeleteBookClassifier_Success(t *testing.T) {
	mock := &mockRepo{
		getSourceBookByCodeFunc: func(_ context.Context, _ repository.AuthContext, _ string) (json.RawMessage, error) {
			return json.RawMessage(`{"id":1}`), nil
		},
		isClassifierInScopeFunc: func(_ context.Context, _ repository.AuthContext, _ string, _ int64, _, _ *int64) (bool, error) {
			return true, nil
		},
		deleteClassifierFunc: func(_ context.Context, _ repository.AuthContext, ct string, id int64) (bool, error) {
			if ct != "talents" || id != 5 {
				t.Errorf("got %s/%d, want talents/5", ct, id)
			}
			return true, nil
		},
	}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/homebrew/source-books/:code/:type/:cid", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteBookClassifier(c)
	})

	w := performRequest(t, router, "DELETE", "/homebrew/source-books/abc/talents/5", nil)
	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d (body=%s)", w.Code, http.StatusNoContent, w.Body.String())
	}
}

func TestRestoreBookClassifier_Success(t *testing.T) {
	mock := &mockRepo{
		getSourceBookByCodeFunc: func(_ context.Context, _ repository.AuthContext, _ string) (json.RawMessage, error) {
			return json.RawMessage(`{"id":1}`), nil
		},
		isClassifierInScopeFunc: func(_ context.Context, _ repository.AuthContext, _ string, _ int64, _, _ *int64) (bool, error) {
			return true, nil
		},
		restoreClassifierFunc: func(_ context.Context, _ repository.AuthContext, _ string, _ int64) (bool, error) {
			return true, nil
		},
	}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/homebrew/source-books/:code/:type/:cid/restore", func(c *gin.Context) {
		withAuth(c)
		handler.RestoreBookClassifier(c)
	})

	w := performRequest(t, router, "POST", "/homebrew/source-books/abc/talents/5/restore", nil)
	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

// =============================================================================
// Hero-scoped classifier handlers
// =============================================================================

func TestUpsertHeroClassifier_Success(t *testing.T) {
	mock := &mockRepo{
		validateHeroAccessFunc: func(_ context.Context, _ repository.AuthContext, heroID int64) error {
			if heroID != 7 {
				t.Errorf("heroID = %d, want 7", heroID)
			}
			return nil
		},
	}
	var captured json.RawMessage
	mock.upsertClassifierFunc = func(_ context.Context, _ repository.AuthContext, _ string, data json.RawMessage) (json.RawMessage, error) {
		captured = data
		return json.RawMessage(`{"id":1}`), nil
	}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/homebrew/heroes/:id/:type", func(c *gin.Context) {
		withAuth(c)
		handler.UpsertHeroClassifier(c)
	})

	w := performRequest(t, router, "POST", "/homebrew/heroes/7/talents", []byte(`{"name":"X"}`))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	var got map[string]any
	_ = json.Unmarshal(captured, &got)
	if got["heroId"].(float64) != 7 {
		t.Errorf("heroId = %v, want 7", got["heroId"])
	}
}

func TestUpsertHeroClassifier_AccessDenied(t *testing.T) {
	mock := &mockRepo{
		validateHeroAccessFunc: func(_ context.Context, _ repository.AuthContext, _ int64) error {
			return &pgconn.PgError{Code: "42501", Message: "denied"}
		},
	}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/homebrew/heroes/:id/:type", func(c *gin.Context) {
		withAuth(c)
		handler.UpsertHeroClassifier(c)
	})

	w := performRequest(t, router, "POST", "/homebrew/heroes/7/talents", []byte(`{}`))
	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestDeleteHeroClassifier_Success(t *testing.T) {
	mock := &mockRepo{
		validateHeroAccessFunc: func(_ context.Context, _ repository.AuthContext, _ int64) error { return nil },
		isClassifierInScopeFunc: func(_ context.Context, _ repository.AuthContext, _ string, _ int64, _, _ *int64) (bool, error) {
			return true, nil
		},
		deleteClassifierFunc: func(_ context.Context, _ repository.AuthContext, _ string, _ int64) (bool, error) {
			return true, nil
		},
	}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/homebrew/heroes/:id/:type/:cid", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteHeroClassifier(c)
	})

	w := performRequest(t, router, "DELETE", "/homebrew/heroes/7/talents/5", nil)
	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestRestoreHeroClassifier_Success(t *testing.T) {
	mock := &mockRepo{
		validateHeroAccessFunc: func(_ context.Context, _ repository.AuthContext, _ int64) error { return nil },
		isClassifierInScopeFunc: func(_ context.Context, _ repository.AuthContext, _ string, _ int64, _, _ *int64) (bool, error) {
			return true, nil
		},
		restoreClassifierFunc: func(_ context.Context, _ repository.AuthContext, _ string, _ int64) (bool, error) {
			return true, nil
		},
	}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/homebrew/heroes/:id/:type/:cid/restore", func(c *gin.Context) {
		withAuth(c)
		handler.RestoreHeroClassifier(c)
	})

	w := performRequest(t, router, "POST", "/homebrew/heroes/7/talents/5/restore", nil)
	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

// =============================================================================
// handleDeleteByString
// =============================================================================

func TestHandleDeleteByString_Success(t *testing.T) {
	router := gin.New()
	router.DELETE("/items/:code", func(c *gin.Context) {
		withAuth(c)
		handleDeleteByString(c, "code", func(_ context.Context, _ repository.AuthContext, code string) (bool, error) {
			if code != "abc" {
				t.Errorf("code = %q, want abc", code)
			}
			return true, nil
		})
	})

	w := performRequest(t, router, "DELETE", "/items/abc", nil)
	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestHandleDeleteByString_NotFound(t *testing.T) {
	router := gin.New()
	router.DELETE("/items/:code", func(c *gin.Context) {
		withAuth(c)
		handleDeleteByString(c, "code", func(_ context.Context, _ repository.AuthContext, _ string) (bool, error) {
			return false, nil
		})
	})

	w := performRequest(t, router, "DELETE", "/items/missing", nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandleDeleteByString_RepoError(t *testing.T) {
	router := gin.New()
	router.DELETE("/items/:code", func(c *gin.Context) {
		withAuth(c)
		handleDeleteByString(c, "code", func(_ context.Context, _ repository.AuthContext, _ string) (bool, error) {
			return false, errors.New("db down")
		})
	})

	w := performRequest(t, router, "DELETE", "/items/x", nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
