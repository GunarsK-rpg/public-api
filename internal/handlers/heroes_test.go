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
	wantBody := `{"heroId":1,"pathCode":"windrunner"}`
	mock.upsertHeroPathFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
		if string(data) != wantBody {
			t.Errorf("data = %s, want %s", data, wantBody)
		}
		return expected, nil
	}

	w := performRequest(t, router, "POST", "/heroes/1/paths", []byte(wantBody))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != string(expected) {
		t.Errorf("body = %s, want %s", w.Body.String(), expected)
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

// =============================================================================
// Comprehensive hero wrapper handler matrix tests
// =============================================================================
//
// The hero handler methods above are thin wrappers around the shared
// handleGetByID, handlePost, and handleDelete helpers. These matrix tests
// exhaustively verify each wrapper across every edge case the helpers can
// produce: successful wiring (correct repo method + correct path param +
// correct payload pass-through), missing auth, invalid IDs, invalid JSON,
// null/nil results, not-found results, and repository errors.
//
// Individual Success tests above serve as readable examples; the matrices
// below guarantee that every hero sub-resource handler gets full edge-case
// coverage and catches copy-paste wiring bugs (wrong repo method, wrong
// path param name).

// -----------------------------------------------------------------------------
// GetByID pattern
// -----------------------------------------------------------------------------

type getByIDMockFn func(id int64) (json.RawMessage, error)

type getByIDHandlerCase struct {
	name    string
	method  func(*Handler) gin.HandlerFunc
	setMock func(*mockRepo, getByIDMockFn)
}

var heroGetByIDHandlerCases = []getByIDHandlerCase{
	{
		name:   "GetHero",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHero },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroSheet",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroSheet },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroSheetFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroAttributes",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroAttributes },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroAttributesFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroDefenses",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroDefenses },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroDefensesFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroDerivedStats",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroDerivedStats },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroDerivedStatsFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroSkills",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroSkills },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroSkillsFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroExpertises",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroExpertises },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroExpertisesFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroTalents",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroTalents },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroTalentsFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroPaths",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroPaths },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroPathsFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroEquipment",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroEquipment },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroEquipmentFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroConditions",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroConditions },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroConditionsFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroInjuries",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroInjuries },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroInjuriesFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroGoals",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroGoals },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroGoalsFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroConnections",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroConnections },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroConnectionsFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroNotes",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroNotes },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroNotesFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
	{
		name:   "GetHeroCultures",
		method: func(h *Handler) gin.HandlerFunc { return h.GetHeroCultures },
		setMock: func(m *mockRepo, fn getByIDMockFn) {
			m.getHeroCulturesFunc = func(_ context.Context, _ repository.AuthContext, id int64) (json.RawMessage, error) {
				return fn(id)
			}
		},
	},
}

func runHeroGetByIDMatrix(t *testing.T, fn func(*testing.T, getByIDHandlerCase)) {
	t.Helper()
	for _, tc := range heroGetByIDHandlerCases {
		t.Run(tc.name, func(t *testing.T) {
			fn(t, tc)
		})
	}
}

func TestHeroGetByIDHandlers_Success(t *testing.T) {
	runHeroGetByIDMatrix(t, func(t *testing.T, tc getByIDHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.GET("/heroes/:id", func(c *gin.Context) {
			withAuth(c)
			tc.method(handler)(c)
		})

		expected := json.RawMessage(`{"id":42}`)
		called := false
		tc.setMock(mock, func(id int64) (json.RawMessage, error) {
			called = true
			if id != 42 {
				t.Errorf("id = %d, want 42", id)
			}
			return expected, nil
		})

		w := performRequest(t, router, "GET", "/heroes/42", nil)
		if !called {
			t.Errorf("repo method not called; handler may be wired to wrong repo func")
		}
		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}
		if w.Body.String() != string(expected) {
			t.Errorf("body = %s, want %s", w.Body.String(), expected)
		}
	})
}

func TestHeroGetByIDHandlers_NoAuth(t *testing.T) {
	runHeroGetByIDMatrix(t, func(t *testing.T, tc getByIDHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.GET("/heroes/:id", func(c *gin.Context) {
			tc.method(handler)(c)
		})

		tc.setMock(mock, func(int64) (json.RawMessage, error) {
			t.Errorf("repo called despite missing auth")
			return nil, nil
		})

		w := performRequest(t, router, "GET", "/heroes/42", nil)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})
}

func TestHeroGetByIDHandlers_InvalidID(t *testing.T) {
	runHeroGetByIDMatrix(t, func(t *testing.T, tc getByIDHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.GET("/heroes/:id", func(c *gin.Context) {
			withAuth(c)
			tc.method(handler)(c)
		})

		tc.setMock(mock, func(int64) (json.RawMessage, error) {
			t.Errorf("repo called despite invalid id")
			return nil, nil
		})

		w := performRequest(t, router, "GET", "/heroes/abc", nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestHeroGetByIDHandlers_NotFound_NilResult(t *testing.T) {
	runHeroGetByIDMatrix(t, func(t *testing.T, tc getByIDHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.GET("/heroes/:id", func(c *gin.Context) {
			withAuth(c)
			tc.method(handler)(c)
		})

		tc.setMock(mock, func(int64) (json.RawMessage, error) {
			return nil, nil
		})

		w := performRequest(t, router, "GET", "/heroes/42", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestHeroGetByIDHandlers_NotFound_NullLiteral(t *testing.T) {
	runHeroGetByIDMatrix(t, func(t *testing.T, tc getByIDHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.GET("/heroes/:id", func(c *gin.Context) {
			withAuth(c)
			tc.method(handler)(c)
		})

		tc.setMock(mock, func(int64) (json.RawMessage, error) {
			return json.RawMessage(`null`), nil
		})

		w := performRequest(t, router, "GET", "/heroes/42", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestHeroGetByIDHandlers_RepoError(t *testing.T) {
	runHeroGetByIDMatrix(t, func(t *testing.T, tc getByIDHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.GET("/heroes/:id", func(c *gin.Context) {
			withAuth(c)
			tc.method(handler)(c)
		})

		tc.setMock(mock, func(int64) (json.RawMessage, error) {
			return nil, errors.New("db error")
		})

		w := performRequest(t, router, "GET", "/heroes/42", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})
}

// -----------------------------------------------------------------------------
// POST (handlePost) pattern
// -----------------------------------------------------------------------------

type postMockFn func(data json.RawMessage) (json.RawMessage, error)

type postHandlerCase struct {
	name    string
	method  func(*Handler) gin.HandlerFunc
	setMock func(*mockRepo, postMockFn)
}

var heroPostHandlerCases = []postHandlerCase{
	{
		name:   "CreateHero",
		method: func(h *Handler) gin.HandlerFunc { return h.CreateHero },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpdateHero",
		method: func(h *Handler) gin.HandlerFunc { return h.UpdateHero },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroAttribute",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroAttribute },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroAttributeFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroDefense",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroDefense },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroDefenseFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroDerivedStat",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroDerivedStat },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroDerivedStatFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroSkill",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroSkill },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroSkillFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroExpertise",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroExpertise },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroExpertiseFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroTalent",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroTalent },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroTalentFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroPath",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroPath },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroPathFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroEquipment",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroEquipment },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroEquipmentFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "AddEquipmentModification",
		method: func(h *Handler) gin.HandlerFunc { return h.AddEquipmentModification },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.addEquipmentModificationFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "AddFavoriteAction",
		method: func(h *Handler) gin.HandlerFunc { return h.AddFavoriteAction },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.addFavoriteActionFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroCondition",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroCondition },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroConditionFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroInjury",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroInjury },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroInjuryFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroGoal",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroGoal },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroGoalFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroConnection",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroConnection },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroConnectionFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroNote",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroNote },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroNoteFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "UpsertHeroCulture",
		method: func(h *Handler) gin.HandlerFunc { return h.UpsertHeroCulture },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.upsertHeroCultureFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "PatchHeroHealth",
		method: func(h *Handler) gin.HandlerFunc { return h.PatchHeroHealth },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.patchHeroHealthFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "PatchHeroFocus",
		method: func(h *Handler) gin.HandlerFunc { return h.PatchHeroFocus },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.patchHeroFocusFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "PatchHeroMagic",
		method: func(h *Handler) gin.HandlerFunc { return h.PatchHeroMagic },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.patchHeroMagicFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
	{
		name:   "PatchHeroCurrency",
		method: func(h *Handler) gin.HandlerFunc { return h.PatchHeroCurrency },
		setMock: func(m *mockRepo, fn postMockFn) {
			m.patchHeroCurrencyFunc = func(_ context.Context, _ repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
				return fn(data)
			}
		},
	},
}

func runHeroPostMatrix(t *testing.T, fn func(*testing.T, postHandlerCase)) {
	t.Helper()
	for _, tc := range heroPostHandlerCases {
		t.Run(tc.name, func(t *testing.T) {
			fn(t, tc)
		})
	}
}

func TestHeroPostHandlers_Success(t *testing.T) {
	runHeroPostMatrix(t, func(t *testing.T, tc postHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.POST("/r", func(c *gin.Context) {
			withAuth(c)
			tc.method(handler)(c)
		})

		body := []byte(`{"test":"payload"}`)
		expected := json.RawMessage(`{"id":42}`)
		called := false
		tc.setMock(mock, func(data json.RawMessage) (json.RawMessage, error) {
			called = true
			if string(data) != string(body) {
				t.Errorf("data = %s, want %s", data, body)
			}
			return expected, nil
		})

		w := performRequest(t, router, "POST", "/r", body)
		if !called {
			t.Errorf("repo method not called; handler may be wired to wrong repo func")
		}
		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}
		if w.Body.String() != string(expected) {
			t.Errorf("body = %s, want %s", w.Body.String(), expected)
		}
	})
}

func TestHeroPostHandlers_NoAuth(t *testing.T) {
	runHeroPostMatrix(t, func(t *testing.T, tc postHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.POST("/r", func(c *gin.Context) {
			tc.method(handler)(c)
		})

		tc.setMock(mock, func(json.RawMessage) (json.RawMessage, error) {
			t.Errorf("repo called despite missing auth")
			return nil, nil
		})

		w := performRequest(t, router, "POST", "/r", []byte(`{}`))
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})
}

func TestHeroPostHandlers_InvalidJSON(t *testing.T) {
	runHeroPostMatrix(t, func(t *testing.T, tc postHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.POST("/r", func(c *gin.Context) {
			withAuth(c)
			tc.method(handler)(c)
		})

		tc.setMock(mock, func(json.RawMessage) (json.RawMessage, error) {
			t.Errorf("repo called despite invalid JSON body")
			return nil, nil
		})

		w := performRequest(t, router, "POST", "/r", []byte(`{not valid json`))
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestHeroPostHandlers_EmptyBody(t *testing.T) {
	runHeroPostMatrix(t, func(t *testing.T, tc postHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.POST("/r", func(c *gin.Context) {
			withAuth(c)
			tc.method(handler)(c)
		})

		tc.setMock(mock, func(json.RawMessage) (json.RawMessage, error) {
			t.Errorf("repo called despite empty body")
			return nil, nil
		})

		w := performRequest(t, router, "POST", "/r", []byte(``))
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestHeroPostHandlers_RepoError(t *testing.T) {
	runHeroPostMatrix(t, func(t *testing.T, tc postHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.POST("/r", func(c *gin.Context) {
			withAuth(c)
			tc.method(handler)(c)
		})

		tc.setMock(mock, func(json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("db error")
		})

		w := performRequest(t, router, "POST", "/r", []byte(`{}`))
		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})
}

// -----------------------------------------------------------------------------
// DELETE (handleDelete) pattern
// -----------------------------------------------------------------------------

type deleteMockFn func(id int64) (bool, error)

type deleteHandlerCase struct {
	name       string
	method     func(*Handler) gin.HandlerFunc
	route      string // gin route, containing the param name the handler expects
	validURL   string // request URL with id=42 at the target param position
	invalidURL string // request URL with non-integer at the target param position
	setMock    func(*mockRepo, deleteMockFn)
}

var heroDeleteHandlerCases = []deleteHandlerCase{
	{
		name:       "DeleteHero",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHero },
		route:      "/heroes/:id",
		validURL:   "/heroes/42",
		invalidURL: "/heroes/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroAttribute",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroAttribute },
		route:      "/heroes/:id/attributes/:subId",
		validURL:   "/heroes/1/attributes/42",
		invalidURL: "/heroes/1/attributes/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroAttributeFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroDefense",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroDefense },
		route:      "/heroes/:id/defenses/:subId",
		validURL:   "/heroes/1/defenses/42",
		invalidURL: "/heroes/1/defenses/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroDefenseFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroDerivedStat",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroDerivedStat },
		route:      "/heroes/:id/derived-stats/:subId",
		validURL:   "/heroes/1/derived-stats/42",
		invalidURL: "/heroes/1/derived-stats/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroDerivedStatFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroSkill",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroSkill },
		route:      "/heroes/:id/skills/:subId",
		validURL:   "/heroes/1/skills/42",
		invalidURL: "/heroes/1/skills/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroSkillFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroExpertise",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroExpertise },
		route:      "/heroes/:id/expertises/:subId",
		validURL:   "/heroes/1/expertises/42",
		invalidURL: "/heroes/1/expertises/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroExpertiseFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroTalent",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroTalent },
		route:      "/heroes/:id/talents/:subId",
		validURL:   "/heroes/1/talents/42",
		invalidURL: "/heroes/1/talents/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroTalentFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroPath",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroPath },
		route:      "/heroes/:id/paths/:subId",
		validURL:   "/heroes/1/paths/42",
		invalidURL: "/heroes/1/paths/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroPathFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroEquipment",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroEquipment },
		route:      "/heroes/:id/equipment/:subId",
		validURL:   "/heroes/1/equipment/42",
		invalidURL: "/heroes/1/equipment/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroEquipmentFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "RemoveEquipmentModification",
		method:     func(h *Handler) gin.HandlerFunc { return h.RemoveEquipmentModification },
		route:      "/heroes/:id/equipment/:subId/modifications/:modId",
		validURL:   "/heroes/1/equipment/2/modifications/42",
		invalidURL: "/heroes/1/equipment/2/modifications/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.removeEquipmentModificationFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "RemoveFavoriteAction",
		method:     func(h *Handler) gin.HandlerFunc { return h.RemoveFavoriteAction },
		route:      "/heroes/:id/favorites/:subId",
		validURL:   "/heroes/1/favorites/42",
		invalidURL: "/heroes/1/favorites/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.removeFavoriteActionFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroCondition",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroCondition },
		route:      "/heroes/:id/conditions/:subId",
		validURL:   "/heroes/1/conditions/42",
		invalidURL: "/heroes/1/conditions/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroConditionFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroInjury",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroInjury },
		route:      "/heroes/:id/injuries/:subId",
		validURL:   "/heroes/1/injuries/42",
		invalidURL: "/heroes/1/injuries/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroInjuryFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroGoal",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroGoal },
		route:      "/heroes/:id/goals/:subId",
		validURL:   "/heroes/1/goals/42",
		invalidURL: "/heroes/1/goals/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroGoalFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroConnection",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroConnection },
		route:      "/heroes/:id/connections/:subId",
		validURL:   "/heroes/1/connections/42",
		invalidURL: "/heroes/1/connections/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroConnectionFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroNote",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroNote },
		route:      "/heroes/:id/notes/:subId",
		validURL:   "/heroes/1/notes/42",
		invalidURL: "/heroes/1/notes/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroNoteFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
	{
		name:       "DeleteHeroCulture",
		method:     func(h *Handler) gin.HandlerFunc { return h.DeleteHeroCulture },
		route:      "/heroes/:id/cultures/:subId",
		validURL:   "/heroes/1/cultures/42",
		invalidURL: "/heroes/1/cultures/abc",
		setMock: func(m *mockRepo, fn deleteMockFn) {
			m.deleteHeroCultureFunc = func(_ context.Context, _ repository.AuthContext, id int64) (bool, error) {
				return fn(id)
			}
		},
	},
}

func runHeroDeleteMatrix(t *testing.T, fn func(*testing.T, deleteHandlerCase)) {
	t.Helper()
	for _, tc := range heroDeleteHandlerCases {
		t.Run(tc.name, func(t *testing.T) {
			fn(t, tc)
		})
	}
}

func TestHeroDeleteHandlers_Success(t *testing.T) {
	runHeroDeleteMatrix(t, func(t *testing.T, tc deleteHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.DELETE(tc.route, func(c *gin.Context) {
			withAuth(c)
			tc.method(handler)(c)
		})

		called := false
		tc.setMock(mock, func(id int64) (bool, error) {
			called = true
			if id != 42 {
				t.Errorf("id = %d, want 42 (wrong path param extracted)", id)
			}
			return true, nil
		})

		w := performRequest(t, router, "DELETE", tc.validURL, nil)
		if !called {
			t.Errorf("repo method not called; handler may be wired to wrong repo func")
		}
		if w.Code != http.StatusNoContent {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
		}
	})
}

func TestHeroDeleteHandlers_NoAuth(t *testing.T) {
	runHeroDeleteMatrix(t, func(t *testing.T, tc deleteHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.DELETE(tc.route, func(c *gin.Context) {
			tc.method(handler)(c)
		})

		tc.setMock(mock, func(int64) (bool, error) {
			t.Errorf("repo called despite missing auth")
			return false, nil
		})

		w := performRequest(t, router, "DELETE", tc.validURL, nil)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})
}

func TestHeroDeleteHandlers_InvalidID(t *testing.T) {
	runHeroDeleteMatrix(t, func(t *testing.T, tc deleteHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.DELETE(tc.route, func(c *gin.Context) {
			withAuth(c)
			tc.method(handler)(c)
		})

		tc.setMock(mock, func(int64) (bool, error) {
			t.Errorf("repo called despite invalid id")
			return false, nil
		})

		w := performRequest(t, router, "DELETE", tc.invalidURL, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestHeroDeleteHandlers_NotFound(t *testing.T) {
	runHeroDeleteMatrix(t, func(t *testing.T, tc deleteHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.DELETE(tc.route, func(c *gin.Context) {
			withAuth(c)
			tc.method(handler)(c)
		})

		tc.setMock(mock, func(int64) (bool, error) {
			return false, nil
		})

		w := performRequest(t, router, "DELETE", tc.validURL, nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestHeroDeleteHandlers_RepoError(t *testing.T) {
	runHeroDeleteMatrix(t, func(t *testing.T, tc deleteHandlerCase) {
		mock := &mockRepo{}
		handler := New(mock, nil)
		router := gin.New()
		router.DELETE(tc.route, func(c *gin.Context) {
			withAuth(c)
			tc.method(handler)(c)
		})

		tc.setMock(mock, func(int64) (bool, error) {
			return false, errors.New("db error")
		})

		w := performRequest(t, router, "DELETE", tc.validURL, nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})
}
