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
// GetHeroes (custom handler with query param binding)
// =============================================================================

func TestGetHeroes_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes", func(c *gin.Context) {
		withAuth(c)
		handler.GetHeroes(c)
	})

	expected := json.RawMessage(`[{"id":1,"name":"Kaladin"}]`)
	mock.getHeroesFunc = func(_ context.Context, _ repository.AuthContext, campaignID *int64) (json.RawMessage, error) {
		if campaignID != nil {
			t.Errorf("campaignID = %v, want nil", *campaignID)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/heroes", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != string(expected) {
		t.Errorf("body = %s, want %s", w.Body.String(), expected)
	}
}

func TestGetHeroes_WithCampaignFilter(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes", func(c *gin.Context) {
		withAuth(c)
		handler.GetHeroes(c)
	})

	expected := json.RawMessage(`[{"id":1}]`)
	mock.getHeroesFunc = func(_ context.Context, _ repository.AuthContext, campaignID *int64) (json.RawMessage, error) {
		if campaignID == nil || *campaignID != 5 {
			t.Errorf("campaignID = %v, want *5", campaignID)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/heroes?campaign_id=5", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetHeroes_InvalidCampaignID(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes", func(c *gin.Context) {
		withAuth(c)
		handler.GetHeroes(c)
	})

	w := performRequest(t, router, "GET", "/heroes?campaign_id=abc", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetHeroes_NoAuth(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes", func(c *gin.Context) {
		handler.GetHeroes(c)
	})

	w := performRequest(t, router, "GET", "/heroes", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGetHeroes_RepoError(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes", func(c *gin.Context) {
		withAuth(c)
		handler.GetHeroes(c)
	})

	mock.getHeroesFunc = func(_ context.Context, _ repository.AuthContext, _ *int64) (json.RawMessage, error) {
		return nil, errors.New("db error")
	}

	w := performRequest(t, router, "GET", "/heroes", nil)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// =============================================================================
// GetHero (handleGetByID delegate)
// =============================================================================

func TestGetHero_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes/:id", func(c *gin.Context) {
		withAuth(c)
		handler.GetHero(c)
	})

	expected := json.RawMessage(`{"id":42,"name":"Kaladin"}`)
	mock.getHeroFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
		if id != 42 {
			t.Errorf("id = %d, want 42", id)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/heroes/42", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != string(expected) {
		t.Errorf("body = %s, want %s", w.Body.String(), expected)
	}
}

func TestGetHero_NotFound(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes/:id", func(c *gin.Context) {
		withAuth(c)
		handler.GetHero(c)
	})

	mock.getHeroFunc = func(_ context.Context, _ repository.AuthContext, _ int64) (json.RawMessage, error) {
		return nil, pgx.ErrNoRows
	}

	w := performRequest(t, router, "GET", "/heroes/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestGetHero_NoAuth(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes/:id", func(c *gin.Context) {
		handler.GetHero(c)
	})

	w := performRequest(t, router, "GET", "/heroes/1", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGetHero_InvalidID(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes/:id", func(c *gin.Context) {
		withAuth(c)
		handler.GetHero(c)
	})

	w := performRequest(t, router, "GET", "/heroes/abc", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// =============================================================================
// GetHeroSheet
// =============================================================================

func TestGetHeroSheet_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes/:id/sheet", func(c *gin.Context) {
		withAuth(c)
		handler.GetHeroSheet(c)
	})

	expected := json.RawMessage(`{"id":1,"attributes":[]}`)
	mock.getHeroSheetFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
		if id != 1 {
			t.Errorf("id = %d, want 1", id)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/heroes/1/sheet", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

// =============================================================================
// CreateHero / UpdateHero (handlePost delegates)
// =============================================================================

func TestCreateHero_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/heroes", func(c *gin.Context) {
		withAuth(c)
		handler.CreateHero(c)
	})

	expected := json.RawMessage(`{"id":1}`)
	mock.upsertHeroFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
		if string(data) != `{"name":"Shallan"}` {
			t.Errorf("data = %s, want %s", data, `{"name":"Shallan"}`)
		}
		return expected, nil
	}

	w := performRequest(t, router, "POST", "/heroes", []byte(`{"name":"Shallan"}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCreateHero_NoAuth(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/heroes", func(c *gin.Context) {
		handler.CreateHero(c)
	})

	w := performRequest(t, router, "POST", "/heroes", []byte(`{}`))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestUpdateHero_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.PUT("/heroes/:id", func(c *gin.Context) {
		withAuth(c)
		handler.UpdateHero(c)
	})

	expected := json.RawMessage(`{"id":1}`)
	mock.upsertHeroFunc = func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
		return expected, nil
	}

	w := performRequest(t, router, "PUT", "/heroes/1", []byte(`{"id":1,"name":"Kaladin"}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

// =============================================================================
// DeleteHero (handleDelete delegate)
// =============================================================================

func TestDeleteHero_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/heroes/:id", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteHero(c)
	})

	mock.deleteHeroFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
		if id != 42 {
			t.Errorf("id = %d, want 42", id)
		}
		return true, nil
	}

	w := performRequest(t, router, "DELETE", "/heroes/42", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestDeleteHero_NotFound(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/heroes/:id", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteHero(c)
	})

	mock.deleteHeroFunc = func(_ context.Context, _ repository.AuthContext, _ int64) (bool, error) {
		return false, nil
	}

	w := performRequest(t, router, "DELETE", "/heroes/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// =============================================================================
// Sub-resource handlers (representative sample)
// =============================================================================

func TestGetHeroAttributes_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes/:id/attributes", func(c *gin.Context) {
		withAuth(c)
		handler.GetHeroAttributes(c)
	})

	expected := json.RawMessage(`[{"type":"STR","value":3}]`)
	mock.getHeroAttributesFunc = func(_ context.Context, _ repository.AuthContext, heroID int64) (json.RawMessage, error) {
		if heroID != 1 {
			t.Errorf("heroID = %d, want 1", heroID)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/heroes/1/attributes", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestUpsertHeroAttribute_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/heroes/:id/attributes", func(c *gin.Context) {
		withAuth(c)
		handler.UpsertHeroAttribute(c)
	})

	expected := json.RawMessage(`{"id":1}`)
	mock.upsertHeroAttributeFunc = func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
		return expected, nil
	}

	w := performRequest(t, router, "POST", "/heroes/1/attributes", []byte(`{"heroId":1,"code":"STR","value":4}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestDeleteHeroAttribute_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/heroes/:id/attributes/:subId", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteHeroAttribute(c)
	})

	mock.deleteHeroAttributeFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
		if id != 5 {
			t.Errorf("id = %d, want 5", id)
		}
		return true, nil
	}

	w := performRequest(t, router, "DELETE", "/heroes/1/attributes/5", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

// =============================================================================
// Resource patch handlers (representative sample)
// =============================================================================

func TestPatchHeroHealth_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.PATCH("/heroes/:id/health", func(c *gin.Context) {
		withAuth(c)
		handler.PatchHeroHealth(c)
	})

	expected := json.RawMessage(`{"currentHp":15}`)
	mock.patchHeroHealthFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
		return expected, nil
	}

	w := performRequest(t, router, "PATCH", "/heroes/1/health", []byte(`{"heroId":1,"currentHp":15}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestPatchHeroHealth_NoAuth(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.PATCH("/heroes/:id/health", func(c *gin.Context) {
		handler.PatchHeroHealth(c)
	})

	w := performRequest(t, router, "PATCH", "/heroes/1/health", []byte(`{}`))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

// =============================================================================
// Equipment modification handlers
// =============================================================================

func TestAddEquipmentModification_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/heroes/:id/equipment/:subId/modifications", func(c *gin.Context) {
		withAuth(c)
		handler.AddEquipmentModification(c)
	})

	expected := json.RawMessage(`{"id":1}`)
	mock.addEquipmentModificationFunc = func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
		return expected, nil
	}

	w := performRequest(t, router, "POST", "/heroes/1/equipment/2/modifications", []byte(`{"equipmentId":2,"code":"KEEN"}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRemoveEquipmentModification_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/heroes/:id/equipment/:subId/modifications/:modId", func(c *gin.Context) {
		withAuth(c)
		handler.RemoveEquipmentModification(c)
	})

	mock.removeEquipmentModificationFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
		if id != 7 {
			t.Errorf("id = %d, want 7", id)
		}
		return true, nil
	}

	w := performRequest(t, router, "DELETE", "/heroes/1/equipment/2/modifications/7", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

// =============================================================================
// Favorite action handlers
// =============================================================================

func TestAddFavoriteAction_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/heroes/:id/favorites", func(c *gin.Context) {
		withAuth(c)
		handler.AddFavoriteAction(c)
	})

	expected := json.RawMessage(`{"id":1}`)
	mock.addFavoriteActionFunc = func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
		return expected, nil
	}

	w := performRequest(t, router, "POST", "/heroes/1/favorites", []byte(`{"heroId":1,"actionCode":"STRIKE"}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRemoveFavoriteAction_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/heroes/:id/favorites/:subId", func(c *gin.Context) {
		withAuth(c)
		handler.RemoveFavoriteAction(c)
	})

	mock.removeFavoriteActionFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
		if id != 3 {
			t.Errorf("id = %d, want 3", id)
		}
		return true, nil
	}

	w := performRequest(t, router, "DELETE", "/heroes/1/favorites/3", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

// =============================================================================
// Hero path handlers
// =============================================================================

func TestGetHeroPaths_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes/:id/paths", func(c *gin.Context) {
		withAuth(c)
		handler.GetHeroPaths(c)
	})

	expected := json.RawMessage(`[{"id":1,"pathCode":"windrunner"}]`)
	mock.getHeroPathsFunc = func(_ context.Context, _ repository.AuthContext, heroID int64) (json.RawMessage, error) {
		if heroID != 1 {
			t.Errorf("heroID = %d, want 1", heroID)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/heroes/1/paths", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != string(expected) {
		t.Errorf("body = %s, want %s", w.Body.String(), expected)
	}
}

func TestUpsertHeroPath_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/heroes/:id/paths", func(c *gin.Context) {
		withAuth(c)
		handler.UpsertHeroPath(c)
	})

	expected := json.RawMessage(`{"id":1}`)
	mock.upsertHeroPathFunc = func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
		return expected, nil
	}

	w := performRequest(t, router, "POST", "/heroes/1/paths", []byte(`{"heroId":1,"pathCode":"windrunner"}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestDeleteHeroPath_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/heroes/:id/paths/:subId", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteHeroPath(c)
	})

	mock.deleteHeroPathFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
		if id != 5 {
			t.Errorf("id = %d, want 5", id)
		}
		return true, nil
	}

	w := performRequest(t, router, "DELETE", "/heroes/1/paths/5", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}
