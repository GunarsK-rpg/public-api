package handlers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/GunarsK-rpg/public-api/internal/repository"
)

// =============================================================================
// SetHeroAvatar
// =============================================================================

func TestSetHeroAvatar_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/heroes/:id/avatar", func(c *gin.Context) {
		withAuth(c)
		handler.SetHeroAvatar(c)
	})

	mock.upsertHeroAvatarFunc = func(_ context.Context, _ repository.AuthContext, heroID int64, avatarKey string) error {
		if heroID != 42 {
			t.Errorf("heroID = %d, want 42", heroID)
		}
		if avatarKey != "abc123" {
			t.Errorf("avatarKey = %q, want %q", avatarKey, "abc123")
		}
		return nil
	}

	w := performRequest(t, router, "POST", "/heroes/42/avatar", []byte(`{"avatarKey":"abc123"}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != `{"avatarKey":"abc123"}` {
		t.Errorf("body = %s, want %s", w.Body.String(), `{"avatarKey":"abc123"}`)
	}
}

func TestSetHeroAvatar_NoAuth(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/heroes/:id/avatar", func(c *gin.Context) {
		handler.SetHeroAvatar(c)
	})

	w := performRequest(t, router, "POST", "/heroes/42/avatar", []byte(`{"avatarKey":"abc"}`))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestSetHeroAvatar_InvalidID(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/heroes/:id/avatar", func(c *gin.Context) {
		withAuth(c)
		handler.SetHeroAvatar(c)
	})

	w := performRequest(t, router, "POST", "/heroes/abc/avatar", []byte(`{"avatarKey":"abc"}`))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestSetHeroAvatar_MissingAvatarKey(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/heroes/:id/avatar", func(c *gin.Context) {
		withAuth(c)
		handler.SetHeroAvatar(c)
	})

	w := performRequest(t, router, "POST", "/heroes/42/avatar", []byte(`{}`))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestSetHeroAvatar_RepoError(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/heroes/:id/avatar", func(c *gin.Context) {
		withAuth(c)
		handler.SetHeroAvatar(c)
	})

	mock.upsertHeroAvatarFunc = func(_ context.Context, _ repository.AuthContext, _ int64, _ string) error {
		return errors.New("db error")
	}

	w := performRequest(t, router, "POST", "/heroes/42/avatar", []byte(`{"avatarKey":"abc123"}`))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// =============================================================================
// DeleteHeroAvatar
// =============================================================================

func TestDeleteHeroAvatar_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/heroes/:id/avatar", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteHeroAvatar(c)
	})

	mock.deleteHeroAvatarFunc = func(_ context.Context, _ repository.AuthContext, heroID int64) error {
		if heroID != 42 {
			t.Errorf("heroID = %d, want 42", heroID)
		}
		return nil
	}

	w := performRequest(t, router, "DELETE", "/heroes/42/avatar", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestDeleteHeroAvatar_NoAuth(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/heroes/:id/avatar", func(c *gin.Context) {
		handler.DeleteHeroAvatar(c)
	})

	w := performRequest(t, router, "DELETE", "/heroes/42/avatar", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestDeleteHeroAvatar_RepoError(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/heroes/:id/avatar", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteHeroAvatar(c)
	})

	mock.deleteHeroAvatarFunc = func(_ context.Context, _ repository.AuthContext, _ int64) error {
		return errors.New("db error")
	}

	w := performRequest(t, router, "DELETE", "/heroes/42/avatar", nil)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// =============================================================================
// SetNpcAvatar
// =============================================================================

func TestSetNpcAvatar_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/campaigns/:id/npcs/:nid/avatar", func(c *gin.Context) {
		withAuth(c)
		handler.SetNpcAvatar(c)
	})

	mock.upsertNpcAvatarFunc = func(_ context.Context, _ repository.AuthContext, npcID, campaignID int64, avatarKey string) error {
		if npcID != 10 {
			t.Errorf("npcID = %d, want 10", npcID)
		}
		if campaignID != 3 {
			t.Errorf("campaignID = %d, want 3", campaignID)
		}
		if avatarKey != "npc-avatar" {
			t.Errorf("avatarKey = %q, want %q", avatarKey, "npc-avatar")
		}
		return nil
	}

	w := performRequest(t, router, "POST", "/campaigns/3/npcs/10/avatar", []byte(`{"avatarKey":"npc-avatar"}`))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestSetNpcAvatar_InvalidCampaignID(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/campaigns/:id/npcs/:nid/avatar", func(c *gin.Context) {
		withAuth(c)
		handler.SetNpcAvatar(c)
	})

	w := performRequest(t, router, "POST", "/campaigns/abc/npcs/10/avatar", []byte(`{"avatarKey":"x"}`))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestSetNpcAvatar_InvalidNpcID(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.POST("/campaigns/:id/npcs/:nid/avatar", func(c *gin.Context) {
		withAuth(c)
		handler.SetNpcAvatar(c)
	})

	w := performRequest(t, router, "POST", "/campaigns/3/npcs/abc/avatar", []byte(`{"avatarKey":"x"}`))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// =============================================================================
// DeleteNpcAvatar
// =============================================================================

func TestDeleteNpcAvatar_Success(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/npcs/:nid/avatar", func(c *gin.Context) {
		withAuth(c)
		handler.DeleteNpcAvatar(c)
	})

	mock.deleteNpcAvatarFunc = func(_ context.Context, _ repository.AuthContext, npcID, campaignID int64) error {
		if npcID != 10 {
			t.Errorf("npcID = %d, want 10", npcID)
		}
		if campaignID != 3 {
			t.Errorf("campaignID = %d, want 3", campaignID)
		}
		return nil
	}

	w := performRequest(t, router, "DELETE", "/campaigns/3/npcs/10/avatar", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestDeleteNpcAvatar_NoAuth(t *testing.T) {
	mock := &mockRepo{}
	handler := New(mock, nil)
	router := gin.New()
	router.DELETE("/campaigns/:id/npcs/:nid/avatar", func(c *gin.Context) {
		handler.DeleteNpcAvatar(c)
	})

	w := performRequest(t, router, "DELETE", "/campaigns/3/npcs/10/avatar", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
