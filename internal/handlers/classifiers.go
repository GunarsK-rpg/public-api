package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	commonHandlers "github.com/GunarsK-portfolio/portfolio-common/handlers"
	"github.com/GunarsK-portfolio/portfolio-common/logger"
)

const (
	classifiersCacheKey = "rpg:classifiers:all"
	classifiersCacheTTL = 1 * time.Hour
)

// GetSourceBooks returns source books visible to the current user.
func (h *Handler) GetSourceBooks(c *gin.Context) {
	handleGet(c, h.repo.GetSourceBooks)
}

// GetAllClassifiers returns all classifiers in a single batch call.
func (h *Handler) GetAllClassifiers(c *gin.Context) {
	ctx := c.Request.Context()

	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Try cache first
	if h.cache != nil {
		cached, err := h.cache.Get(ctx, classifiersCacheKey)
		if err != nil {
			logger.GetLogger(c).Warn("redis cache get failed", "key", classifiersCacheKey, "error", err)
		}
		if cached != nil {
			c.Header("Cache-Control", "private, max-age=3600")
			c.Data(http.StatusOK, "application/json", cached)
			return
		}
	}

	// Cache miss: fetch from DB
	result, err := h.repo.GetAllClassifiers(ctx, auth)
	if err != nil {
		HandlePgxError(c, err)
		return
	}

	// Store in cache (best-effort)
	if h.cache != nil {
		if err := h.cache.Set(ctx, classifiersCacheKey, result, classifiersCacheTTL); err != nil {
			logger.GetLogger(c).Warn("redis cache set failed", "key", classifiersCacheKey, "error", err)
		}
	}

	c.Header("Cache-Control", "private, max-age=3600")
	c.Data(http.StatusOK, "application/json", result)
}
