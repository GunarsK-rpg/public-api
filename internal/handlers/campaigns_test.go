package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/GunarsK-rpg/public-api/internal/repository"
)

// =============================================================================
// GetCampaigns (handleGet delegate)
// =============================================================================

func TestGetCampaigns_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns", func(c *gin.Context) {
		withAuth(c)
		handler.GetCampaigns(c)
	})

	expected := json.RawMessage(`[{"id":1,"name":"Bridge Four"}]`)
	mock.getCampaignsFunc = func(_ context.Context, _ repository.AuthContext) (json.RawMessage, error) {
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/campaigns", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != string(expected) {
		t.Errorf("body = %s, want %s", w.Body.String(), expected)
	}
}

func TestGetCampaigns_NoAuth(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns", func(c *gin.Context) {
		handler.GetCampaigns(c)
	})

	w := performRequest(t, router, "GET", "/campaigns", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGetCampaigns_RepoError(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns", func(c *gin.Context) {
		withAuth(c)
		handler.GetCampaigns(c)
	})

	mock.getCampaignsFunc = func(_ context.Context, _ repository.AuthContext) (json.RawMessage, error) {
		return nil, errors.New("db error")
	}

	w := performRequest(t, router, "GET", "/campaigns", nil)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// =============================================================================
// GetCampaign (handleGetByID delegate)
// =============================================================================

func TestGetCampaign_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id", func(c *gin.Context) {
		withAuth(c)
		handler.GetCampaign(c)
	})

	expected := json.RawMessage(`{"id":3,"name":"Bridge Four"}`)
	mock.getCampaignFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
		if id != 3 {
			t.Errorf("id = %d, want 3", id)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/campaigns/3", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetCampaign_NotFound(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id", func(c *gin.Context) {
		withAuth(c)
		handler.GetCampaign(c)
	})

	mock.getCampaignFunc = func(_ context.Context, _ repository.AuthContext, _ int64) (json.RawMessage, error) {
		return nil, pgx.ErrNoRows
	}

	w := performRequest(t, router, "GET", "/campaigns/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// =============================================================================
// GetCampaignByCode (handleGetByString delegate)
// =============================================================================

func TestGetCampaignByCode_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/join/:code", func(c *gin.Context) {
		withAuth(c)
		handler.GetCampaignByCode(c)
	})

	expected := json.RawMessage(`{"id":1,"name":"Bridge Four"}`)
	mock.getCampaignByCodeFunc = func(_ context.Context, _ repository.AuthContext, code string) (json.RawMessage, error) {
		if code != "ABC123" {
			t.Errorf("code = %q, want %q", code, "ABC123")
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/campaigns/join/ABC123", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetCampaignByCode_NotFound(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/join/:code", func(c *gin.Context) {
		withAuth(c)
		handler.GetCampaignByCode(c)
	})

	mock.getCampaignByCodeFunc = func(_ context.Context, _ repository.AuthContext, _ string) (json.RawMessage, error) {
		return json.RawMessage("null"), nil
	}

	w := performRequest(t, router, "GET", "/campaigns/join/BADCODE", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// =============================================================================
// CreateCampaign / UpdateCampaign (handlePost delegates)
// =============================================================================

func TestCreateCampaign_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/campaigns", func(c *gin.Context) {
		withAuth(c)
		handler.CreateCampaign(c)
	})

	expected := json.RawMessage(`{"id":1}`)
	mock.upsertCampaignFunc = func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
		return expected, nil
	}

	w := performRequest(t, router, "POST", "/campaigns", []byte(`{"name":"Bridge Four"}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

// =============================================================================
// DeleteCampaign (handleDelete delegate)
// =============================================================================

func TestDeleteCampaign_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteCampaign(c)
	})

	mock.deleteCampaignFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
		if id != 3 {
			t.Errorf("id = %d, want 3", id)
		}
		return true, nil
	}

	w := performRequest(t, router, "DELETE", "/campaigns/3", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

// =============================================================================
// RemoveHeroFromCampaign (custom multi-param handler)
// =============================================================================

func TestRemoveHeroFromCampaign_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/heroes/:hid", func(c *gin.Context) {
		withAuth(c)
		handler.RemoveHeroFromCampaign(c)
	})

	mock.removeHeroFromCampaignFunc = func(_ context.Context, _ repository.AuthContext, heroID, campaignID int64) (bool, error) {
		if heroID != 7 {
			t.Errorf("heroID = %d, want 7", heroID)
		}
		if campaignID != 3 {
			t.Errorf("campaignID = %d, want 3", campaignID)
		}
		return true, nil
	}

	w := performRequest(t, router, "DELETE", "/campaigns/3/heroes/7", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestRemoveHeroFromCampaign_NoAuth(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/heroes/:hid", func(c *gin.Context) {
		handler.RemoveHeroFromCampaign(c)
	})

	w := performRequest(t, router, "DELETE", "/campaigns/3/heroes/7", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestRemoveHeroFromCampaign_InvalidHeroID(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/heroes/:hid", func(c *gin.Context) {
		withAuth(c)
		handler.RemoveHeroFromCampaign(c)
	})

	w := performRequest(t, router, "DELETE", "/campaigns/3/heroes/abc", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRemoveHeroFromCampaign_InvalidCampaignID(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/heroes/:hid", func(c *gin.Context) {
		withAuth(c)
		handler.RemoveHeroFromCampaign(c)
	})

	w := performRequest(t, router, "DELETE", "/campaigns/abc/heroes/7", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRemoveHeroFromCampaign_NotFound(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/heroes/:hid", func(c *gin.Context) {
		withAuth(c)
		handler.RemoveHeroFromCampaign(c)
	})

	mock.removeHeroFromCampaignFunc = func(_ context.Context, _ repository.AuthContext, _, _ int64) (bool, error) {
		return false, nil
	}

	w := performRequest(t, router, "DELETE", "/campaigns/3/heroes/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestRemoveHeroFromCampaign_RepoError(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/heroes/:hid", func(c *gin.Context) {
		withAuth(c)
		handler.RemoveHeroFromCampaign(c)
	})

	mock.removeHeroFromCampaignFunc = func(_ context.Context, _ repository.AuthContext, _, _ int64) (bool, error) {
		return false, errors.New("db error")
	}

	w := performRequest(t, router, "DELETE", "/campaigns/3/heroes/7", nil)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
