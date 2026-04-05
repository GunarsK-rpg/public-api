package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	commonHandlers "github.com/GunarsK-portfolio/portfolio-common/handlers"
	"github.com/GunarsK-portfolio/portfolio-common/logger"

	"github.com/GunarsK-rpg/public-api/internal/models/requests"
)

const (
	classifiersCacheGlobalKey  = "rpg:classifiers:global"
	classifiersCacheSBPrefix   = "rpg:classifiers:sb:"
	classifiersCacheHeroPrefix = "rpg:classifiers:hero:"
	classifiersCacheTTL        = 1 * time.Hour
)

// GetSourceBooks returns source books visible to the current user.
func (h *Handler) GetSourceBooks(c *gin.Context) {
	handleGet(c, h.repo.GetSourceBooks)
}

// GetAllClassifiers returns classifiers scoped by optional campaignId and/or heroId.
// Always returns {"global": {...}, "sourceBooks": [...], "hero": {...}}.
// Global classifiers are always included. Source books and hero layers added when params provided.
func (h *Handler) GetAllClassifiers(c *gin.Context) {
	var query requests.GetClassifiersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	// Always fetch global classifiers (unscoped + global scoped)
	global, err := h.fetchClassifiers(c, classifiersCacheGlobalKey, `{"sourceBookId": null}`)
	if err != nil {
		HandlePgxError(c, err)
		return
	}

	// Fetch campaign source books if campaignId provided
	sourceBooks := make([]json.RawMessage, 0)
	if query.CampaignID != nil {
		auth, err := GetAuthContext(c)
		if err != nil {
			commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		sbIDs, err := h.repo.GetCampaignSourceBookIDs(c.Request.Context(), auth, *query.CampaignID)
		if err != nil {
			HandlePgxError(c, err)
			return
		}

		for _, sbID := range sbIDs {
			cacheKey := fmt.Sprintf("%s%d", classifiersCacheSBPrefix, sbID)
			filter := fmt.Sprintf(`{"sourceBookId": %d}`, sbID)

			sbData, err := h.fetchClassifiers(c, cacheKey, filter)
			if err != nil {
				HandlePgxError(c, err)
				return
			}
			sourceBooks = append(sourceBooks, sbData)
		}
	}

	// Fetch hero classifiers if heroId provided
	var hero json.RawMessage
	if query.HeroID != nil {
		auth, err := GetAuthContext(c)
		if err != nil {
			commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		if err := h.repo.ValidateHeroAccess(c.Request.Context(), auth, *query.HeroID); err != nil {
			HandlePgxError(c, err)
			return
		}

		cacheKey := fmt.Sprintf("%s%d", classifiersCacheHeroPrefix, *query.HeroID)
		filter := fmt.Sprintf(`{"heroId": %d}`, *query.HeroID)

		hero, err = h.fetchClassifiers(c, cacheKey, filter)
		if err != nil {
			HandlePgxError(c, err)
			return
		}
	}

	// Build response
	sbJSON, _ := json.Marshal(sourceBooks)
	if hero == nil {
		hero = json.RawMessage(`{}`)
	}
	result := json.RawMessage(fmt.Sprintf(`{"global":%s,"sourceBooks":%s,"hero":%s}`, global, sbJSON, hero))

	c.Data(http.StatusOK, "application/json", result)
}

// fetchClassifiers returns classifiers from Redis cache or DB, caching on miss.
func (h *Handler) fetchClassifiers(c *gin.Context, cacheKey string, filter string) (json.RawMessage, error) {
	ctx := c.Request.Context()
	log := logger.GetLogger(c)

	if h.cache != nil {
		cached, err := h.cache.Get(ctx, cacheKey)
		if err != nil {
			log.Warn("redis cache get failed", "key", cacheKey, "error", err)
		}
		if cached != nil {
			return cached, nil
		}
	}

	auth, err := GetAuthContext(c)
	if err != nil {
		return nil, err
	}

	result, err := h.repo.GetClassifiersFiltered(ctx, auth, json.RawMessage(filter))
	if err != nil {
		return nil, err
	}

	if h.cache != nil {
		if err := h.cache.Set(ctx, cacheKey, result, classifiersCacheTTL); err != nil {
			log.Warn("redis cache set failed", "key", cacheKey, "error", err)
		}
	}

	return result, nil
}
