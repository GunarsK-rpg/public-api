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
	"github.com/GunarsK-rpg/public-api/internal/repository"
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

// GetAllClassifiers returns classifiers scoped by optional campaignId, heroId or sourceBookId.
func (h *Handler) GetAllClassifiers(c *gin.Context) {
	var query requests.GetClassifiersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, "invalid query parameters")
		return
	}
	if err := query.Validate(); err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	auth, ok := h.requireAuth(c)
	if !ok {
		return
	}

	var sbIDs []int64
	if query.CampaignID != nil {
		ids, err := h.repo.GetCampaignSourceBookIDs(c.Request.Context(), auth, *query.CampaignID)
		if err != nil {
			HandlePgxError(c, err)
			return
		}
		sbIDs = ids
	}

	if query.SourceBookID != nil {
		depIDs, err := h.repo.GetSourceBookDependencyIDs(c.Request.Context(), auth, *query.SourceBookID)
		if err != nil {
			HandlePgxError(c, err)
			return
		}
		sbIDs = append([]int64{*query.SourceBookID}, depIDs...)
	}

	sourceBooks := make([]json.RawMessage, 0, len(sbIDs))
	for _, sbID := range sbIDs {
		sbData, err := h.fetchSourceBookClassifiers(c, auth, sbID)
		if err != nil {
			HandlePgxError(c, err)
			return
		}
		sourceBooks = append(sourceBooks, sbData)
	}

	var hero json.RawMessage
	if query.HeroID != nil {
		cacheKey := fmt.Sprintf("%s%d", classifiersCacheHeroPrefix, *query.HeroID)
		var err error
		hero, err = h.fetchClassifiers(c, auth, cacheKey, nil, query.HeroID)
		if err != nil {
			HandlePgxError(c, err)
			return
		}
	}

	global, err := h.fetchClassifiers(c, auth, classifiersCacheGlobalKey, nil, nil)
	if err != nil {
		HandlePgxError(c, err)
		return
	}

	sbJSON, _ := json.Marshal(sourceBooks)
	if hero == nil {
		hero = json.RawMessage(`{}`)
	}
	result := json.RawMessage(fmt.Sprintf(`{"global":%s,"sourceBooks":%s,"hero":%s}`, global, sbJSON, hero))

	c.Data(http.StatusOK, "application/json", result)
}

func (h *Handler) fetchSourceBookClassifiers(c *gin.Context, auth repository.AuthContext, sbID int64) (json.RawMessage, error) {
	cacheKey := fmt.Sprintf("%s%d", classifiersCacheSBPrefix, sbID)
	return h.fetchClassifiers(c, auth, cacheKey, &sbID, nil)
}

// fetchClassifiers gates the cache by filter scope, then returns from Redis or DB.
func (h *Handler) fetchClassifiers(
	c *gin.Context,
	auth repository.AuthContext,
	cacheKey string,
	sourceBookID *int64,
	heroID *int64,
) (json.RawMessage, error) {
	ctx := c.Request.Context()
	log := logger.GetLogger(c)

	if sourceBookID != nil {
		if err := h.repo.RequireSourceBookAccessible(ctx, auth, *sourceBookID); err != nil {
			return nil, err
		}
	}
	if heroID != nil {
		if err := h.repo.ValidateHeroAccess(ctx, auth, *heroID); err != nil {
			return nil, err
		}
	}

	if h.cache != nil {
		cached, err := h.cache.Get(ctx, cacheKey)
		if err != nil {
			log.Warn("redis cache get failed", "key", cacheKey, "error", err)
		}
		if cached != nil {
			return cached, nil
		}
	}

	result, err := h.repo.GetClassifiersFiltered(ctx, auth, buildClassifierFilter(sourceBookID, heroID))
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

func buildClassifierFilter(sourceBookID, heroID *int64) json.RawMessage {
	switch {
	case sourceBookID != nil:
		return json.RawMessage(fmt.Sprintf(`{"sourceBookId": %d}`, *sourceBookID))
	case heroID != nil:
		return json.RawMessage(fmt.Sprintf(`{"heroId": %d}`, *heroID))
	default:
		return json.RawMessage(`{"sourceBookId": null}`)
	}
}
