package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/GunarsK-rpg/public-api/internal/repository"
)

// =============================================================================
// GetNpcOptions / GetNpcLibrary (handleGetByID delegates)
// =============================================================================

func TestGetNpcOptions_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id/npcs", func(c *gin.Context) {
		withAuth(c)
		handler.GetNpcOptions(c)
	})

	expected := json.RawMessage(`[{"id":1,"name":"Szeth"}]`)
	mock.getNpcOptionsFunc = func(_ context.Context, _ repository.AuthContext, campaignID int64) (json.RawMessage, error) {
		if campaignID != 3 {
			t.Errorf("campaignID = %d, want 3", campaignID)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/campaigns/3/npcs", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetNpcLibrary_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id/npcs/library", func(c *gin.Context) {
		withAuth(c)
		handler.GetNpcLibrary(c)
	})

	expected := json.RawMessage(`[{"id":1}]`)
	mock.getNpcLibraryFunc = func(_ context.Context, _ repository.AuthContext, campaignID int64) (json.RawMessage, error) {
		if campaignID != 3 {
			t.Errorf("campaignID = %d, want 3", campaignID)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/campaigns/3/npcs/library", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

// =============================================================================
// GetNpc (custom multi-param handler: nid + id)
// =============================================================================

func TestGetNpc_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id/npcs/:nid", func(c *gin.Context) {
		withAuth(c)
		handler.GetNpc(c)
	})

	expected := json.RawMessage(`{"id":10,"name":"Szeth"}`)
	mock.getNpcFunc = func(_ context.Context, _ repository.AuthContext, npcID, campaignID int64) (json.RawMessage, error) {
		if npcID != 10 {
			t.Errorf("npcID = %d, want 10", npcID)
		}
		if campaignID != 3 {
			t.Errorf("campaignID = %d, want 3", campaignID)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/campaigns/3/npcs/10", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != string(expected) {
		t.Errorf("body = %s, want %s", w.Body.String(), expected)
	}
}

func TestGetNpc_NoAuth(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id/npcs/:nid", func(c *gin.Context) {
		handler.GetNpc(c)
	})

	w := performRequest(t, router, "GET", "/campaigns/3/npcs/10", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGetNpc_InvalidNpcID(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id/npcs/:nid", func(c *gin.Context) {
		withAuth(c)
		handler.GetNpc(c)
	})

	w := performRequest(t, router, "GET", "/campaigns/3/npcs/abc", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetNpc_InvalidCampaignID(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id/npcs/:nid", func(c *gin.Context) {
		withAuth(c)
		handler.GetNpc(c)
	})

	w := performRequest(t, router, "GET", "/campaigns/abc/npcs/10", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetNpc_NullResult(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id/npcs/:nid", func(c *gin.Context) {
		withAuth(c)
		handler.GetNpc(c)
	})

	mock.getNpcFunc = func(_ context.Context, _ repository.AuthContext, _, _ int64) (json.RawMessage, error) {
		return json.RawMessage("null"), nil
	}

	w := performRequest(t, router, "GET", "/campaigns/3/npcs/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestGetNpc_RepoError(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id/npcs/:nid", func(c *gin.Context) {
		withAuth(c)
		handler.GetNpc(c)
	})

	mock.getNpcFunc = func(_ context.Context, _ repository.AuthContext, _, _ int64) (json.RawMessage, error) {
		return nil, errors.New("db error")
	}

	w := performRequest(t, router, "GET", "/campaigns/3/npcs/10", nil)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// =============================================================================
// GetNpcByID (handleGetByID delegate)
// =============================================================================

func TestGetNpcByID_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/npcs/:id", func(c *gin.Context) {
		withAuth(c)
		handler.GetNpcByID(c)
	})

	expected := json.RawMessage(`{"id":10}`)
	mock.getNpcByIDFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
		if id != 10 {
			t.Errorf("id = %d, want 10", id)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/npcs/10", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

// =============================================================================
// CreateNpc / UpdateNpc (handlePost delegates)
// =============================================================================

func TestCreateNpc_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/campaigns/:id/npcs", func(c *gin.Context) {
		withAuth(c)
		handler.CreateNpc(c)
	})

	expected := json.RawMessage(`{"id":1}`)
	mock.upsertNpcFunc = func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
		return expected, nil
	}

	w := performRequest(t, router, "POST", "/campaigns/3/npcs", []byte(`{"campaignId":3,"name":"Szeth"}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

// =============================================================================
// DeleteNpc (custom multi-param handler: nid + id)
// =============================================================================

func TestDeleteNpc_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/npcs/:nid", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteNpc(c)
	})

	mock.deleteNpcFunc = func(_ context.Context, _ repository.AuthContext, npcID, campaignID int64) (bool, error) {
		if npcID != 10 {
			t.Errorf("npcID = %d, want 10", npcID)
		}
		if campaignID != 3 {
			t.Errorf("campaignID = %d, want 3", campaignID)
		}
		return true, nil
	}

	w := performRequest(t, router, "DELETE", "/campaigns/3/npcs/10", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestDeleteNpc_NoAuth(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/npcs/:nid", func(c *gin.Context) {
		handler.DeleteNpc(c)
	})

	w := performRequest(t, router, "DELETE", "/campaigns/3/npcs/10", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestDeleteNpc_InvalidNpcID(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/npcs/:nid", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteNpc(c)
	})

	w := performRequest(t, router, "DELETE", "/campaigns/3/npcs/abc", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestDeleteNpc_NotFound(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/npcs/:nid", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteNpc(c)
	})

	mock.deleteNpcFunc = func(_ context.Context, _ repository.AuthContext, _, _ int64) (bool, error) {
		return false, nil
	}

	w := performRequest(t, router, "DELETE", "/campaigns/3/npcs/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// =============================================================================
// GetCombat (custom multi-param handler: cid + id)
// =============================================================================

func TestGetCombat_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id/combats/:cid", func(c *gin.Context) {
		withAuth(c)
		handler.GetCombat(c)
	})

	expected := json.RawMessage(`{"id":5,"round":2}`)
	mock.getCombatFunc = func(_ context.Context, _ repository.AuthContext, combatID, campaignID int64) (json.RawMessage, error) {
		if combatID != 5 {
			t.Errorf("combatID = %d, want 5", combatID)
		}
		if campaignID != 3 {
			t.Errorf("campaignID = %d, want 3", campaignID)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/campaigns/3/combats/5", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetCombat_NoAuth(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id/combats/:cid", func(c *gin.Context) {
		handler.GetCombat(c)
	})

	w := performRequest(t, router, "GET", "/campaigns/3/combats/5", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGetCombat_InvalidCombatID(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id/combats/:cid", func(c *gin.Context) {
		withAuth(c)
		handler.GetCombat(c)
	})

	w := performRequest(t, router, "GET", "/campaigns/3/combats/abc", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetCombat_NullResult(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id/combats/:cid", func(c *gin.Context) {
		withAuth(c)
		handler.GetCombat(c)
	})

	mock.getCombatFunc = func(_ context.Context, _ repository.AuthContext, _, _ int64) (json.RawMessage, error) {
		return nil, nil
	}

	w := performRequest(t, router, "GET", "/campaigns/3/combats/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// =============================================================================
// DeleteCombat (custom multi-param handler: cid + id)
// =============================================================================

func TestDeleteCombat_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/combats/:cid", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteCombat(c)
	})

	mock.deleteCombatFunc = func(_ context.Context, _ repository.AuthContext, combatID, campaignID int64) (bool, error) {
		if combatID != 5 {
			t.Errorf("combatID = %d, want 5", combatID)
		}
		if campaignID != 3 {
			t.Errorf("campaignID = %d, want 3", campaignID)
		}
		return true, nil
	}

	w := performRequest(t, router, "DELETE", "/campaigns/3/combats/5", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestDeleteCombat_NotFound(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/combats/:cid", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteCombat(c)
	})

	mock.deleteCombatFunc = func(_ context.Context, _ repository.AuthContext, _, _ int64) (bool, error) {
		return false, nil
	}

	w := performRequest(t, router, "DELETE", "/campaigns/3/combats/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// =============================================================================
// GetCombats (handleGetByID delegate)
// =============================================================================

func TestGetCombats_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/campaigns/:id/combats", func(c *gin.Context) {
		withAuth(c)
		handler.GetCombats(c)
	})

	expected := json.RawMessage(`[{"id":1}]`)
	mock.getCombatsFunc = func(_ context.Context, _ repository.AuthContext, campaignID int64) (json.RawMessage, error) {
		if campaignID != 3 {
			t.Errorf("campaignID = %d, want 3", campaignID)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/campaigns/3/combats", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

// =============================================================================
// EndCombatRound (handlePost delegate)
// =============================================================================

func TestEndCombatRound_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/campaigns/:id/combats/:cid/end-round", func(c *gin.Context) {
		withAuth(c)
		handler.EndCombatRound(c)
	})

	expected := json.RawMessage(`{"round":3}`)
	mock.endCombatRoundFunc = func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
		return expected, nil
	}

	w := performRequest(t, router, "POST", "/campaigns/3/combats/5/end-round", []byte(`{"combatId":5}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

// =============================================================================
// NPC Instance handlers
// =============================================================================

func TestGetNpcInstance_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/npc-instances/:id", func(c *gin.Context) {
		withAuth(c)
		handler.GetNpcInstance(c)
	})

	expected := json.RawMessage(`{"id":1}`)
	mock.getNpcInstanceFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
		if id != 1 {
			t.Errorf("id = %d, want 1", id)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/npc-instances/1", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCreateNpcInstance_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/npc-instances", func(c *gin.Context) {
		withAuth(c)
		handler.CreateNpcInstance(c)
	})

	expected := json.RawMessage(`{"id":1}`)
	mock.createNpcInstanceFunc = func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
		return expected, nil
	}

	w := performRequest(t, router, "POST", "/npc-instances", []byte(`{"npcId":10,"combatId":5}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestPatchNpcInstance_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.PATCH("/npc-instances/:id", func(c *gin.Context) {
		withAuth(c)
		handler.PatchNpcInstance(c)
	})

	expected := json.RawMessage(`{"id":1}`)
	mock.patchNpcInstanceFunc = func(_ context.Context, _ repository.AuthContext, _ json.RawMessage) (json.RawMessage, error) {
		return expected, nil
	}

	w := performRequest(t, router, "PATCH", "/npc-instances/1", []byte(`{"id":1,"field":"hp","value":10}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestDeleteNpcInstance_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/npc-instances/:id", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteNpcInstance(c)
	})

	mock.deleteNpcInstanceFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
		if id != 1 {
			t.Errorf("id = %d, want 1", id)
		}
		return true, nil
	}

	w := performRequest(t, router, "DELETE", "/npc-instances/1", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

// =============================================================================
// Companion handlers
// =============================================================================

func TestGetHeroCompanions_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes/:id/companions", func(c *gin.Context) {
		withAuth(c)
		handler.GetHeroCompanions(c)
	})

	expected := json.RawMessage(`[{"id":1}]`)
	mock.getHeroNpcInstancesFunc = func(_ context.Context, _ repository.AuthContext, heroID int64) (json.RawMessage, error) {
		if heroID != 42 {
			t.Errorf("heroID = %d, want 42", heroID)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/heroes/42/companions", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetCompanionNpcOptions_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.GET("/heroes/:id/companion-npcs", func(c *gin.Context) {
		withAuth(c)
		handler.GetCompanionNpcOptions(c)
	})

	expected := json.RawMessage(`[{"id":1}]`)
	mock.getCompanionNpcOptionsFunc = func(_ context.Context, _ repository.AuthContext, heroID int64) (json.RawMessage, error) {
		if heroID != 42 {
			t.Errorf("heroID = %d, want 42", heroID)
		}
		return expected, nil
	}

	w := performRequest(t, router, "GET", "/heroes/42/companion-npcs", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}
