package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	commonHandlers "github.com/GunarsK-portfolio/portfolio-common/handlers"
	"github.com/GunarsK-portfolio/portfolio-common/logger"

	"github.com/GunarsK-rpg/public-api/internal/constants"
	"github.com/GunarsK-rpg/public-api/internal/repository"
)

// classifierScope holds the per-request context needed by every homebrew
// classifier write: identity, scope to inject, and cache key to invalidate.
// Exactly one of sourceBookID / heroID is non-nil.
type classifierScope struct {
	auth         repository.AuthContext
	sourceBookID *int64
	heroID       *int64
	invalidate   func(context.Context, *gin.Context)
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

// parseClassifierType reads :type and validates it against the allow-list.
func parseClassifierType(c *gin.Context) (urlType string, ok bool) {
	urlType = c.Param("type")
	if _, valid := constants.ClassifierTypeSuffix(urlType); !valid {
		commonHandlers.RespondError(c, http.StatusBadRequest, fmt.Sprintf("unknown classifier type: %q", urlType))
		return "", false
	}
	return urlType, true
}

// resolveSourceBookIDByCode looks up the integer id of a source book by its
// UUID code, responding 404 / 500 on miss / DB error.
func (h *Handler) resolveSourceBookIDByCode(c *gin.Context, auth repository.AuthContext, code string) (int64, bool) {
	raw, err := h.repo.GetSourceBookByCode(c.Request.Context(), auth, code)
	if err != nil {
		HandlePgxError(c, err)
		return 0, false
	}
	if raw == nil || string(raw) == "null" {
		commonHandlers.RespondError(c, http.StatusNotFound, "source book not found")
		return 0, false
	}
	id, ok := parseSourceBookID(raw)
	if !ok {
		commonHandlers.LogAndRespondError(c, http.StatusInternalServerError,
			fmt.Errorf("failed to parse source book payload"), "failed to parse source book")
		return 0, false
	}
	return id, true
}

// requireAuth fetches the auth context or responds 401.
func (h *Handler) requireAuth(c *gin.Context) (repository.AuthContext, bool) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return repository.AuthContext{}, false
	}
	return auth, true
}

// bookPreamble runs auth + type validation + book lookup + cache binding.
func (h *Handler) bookPreamble(c *gin.Context) (urlType string, sc *classifierScope, ok bool) {
	auth, ok := h.requireAuth(c)
	if !ok {
		return "", nil, false
	}
	urlType, ok = parseClassifierType(c)
	if !ok {
		return "", nil, false
	}
	bookID, ok := h.resolveSourceBookIDByCode(c, auth, c.Param("code"))
	if !ok {
		return "", nil, false
	}
	return urlType, &classifierScope{
		auth:         auth,
		sourceBookID: &bookID,
		invalidate:   func(ctx context.Context, c *gin.Context) { h.invalidateBookCache(ctx, c, bookID) },
	}, true
}

// heroPreamble runs auth + type validation + hero id parse + access check + cache binding.
func (h *Handler) heroPreamble(c *gin.Context) (urlType string, sc *classifierScope, ok bool) {
	auth, ok := h.requireAuth(c)
	if !ok {
		return "", nil, false
	}
	urlType, ok = parseClassifierType(c)
	if !ok {
		return "", nil, false
	}
	heroID, err := getPathParamInt64(c, "id")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return "", nil, false
	}
	if err := h.repo.ValidateHeroAccess(c.Request.Context(), auth, heroID); err != nil {
		HandlePgxError(c, err)
		return "", nil, false
	}
	return urlType, &classifierScope{
		auth:       auth,
		heroID:     &heroID,
		invalidate: func(ctx context.Context, c *gin.Context) { h.invalidateHeroCache(ctx, c, heroID) },
	}, true
}

// invalidateBookCache deletes the per-book classifier cache key. Errors logged + swallowed.
func (h *Handler) invalidateBookCache(ctx context.Context, c *gin.Context, bookID int64) {
	if h.cache == nil {
		return
	}
	key := fmt.Sprintf("%s%d", classifiersCacheSBPrefix, bookID)
	if err := h.cache.Delete(ctx, key); err != nil {
		logger.GetLogger(c).Warn("redis cache delete failed", "key", key, "error", err)
	}
}

// invalidateHeroCache deletes the per-hero classifier cache key.
func (h *Handler) invalidateHeroCache(ctx context.Context, c *gin.Context, heroID int64) {
	if h.cache == nil {
		return
	}
	key := fmt.Sprintf("%s%d", classifiersCacheHeroPrefix, heroID)
	if err := h.cache.Delete(ctx, key); err != nil {
		logger.GetLogger(c).Warn("redis cache delete failed", "key", key, "error", err)
	}
}

// readJSONBody reads and validates the request body as JSON. Empty body becomes "{}".
func readJSONBody(c *gin.Context) (json.RawMessage, bool) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, "failed to read request body")
		return nil, false
	}
	if len(body) == 0 {
		body = []byte("{}")
	}
	if !json.Valid(body) {
		commonHandlers.RespondError(c, http.StatusBadRequest, "invalid JSON")
		return nil, false
	}
	return body, true
}

// ----------------------------------------------------------------------------
// Generic classifier workers (book + hero share these)
// ----------------------------------------------------------------------------

func (h *Handler) doUpsertClassifier(c *gin.Context, urlType string, sc *classifierScope) {
	body, ok := readJSONBody(c)
	if !ok {
		return
	}
	merged, err := injectScope(body, sc.sourceBookID, sc.heroID)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.repo.UpsertClassifier(c.Request.Context(), sc.auth, urlType, merged)
	if err != nil {
		HandlePgxError(c, err)
		return
	}
	sc.invalidate(c.Request.Context(), c)
	c.Data(http.StatusOK, "application/json", result)
}

func (h *Handler) doDeleteClassifier(c *gin.Context, urlType string, sc *classifierScope) {
	cid, ok := h.requireCidInScope(c, urlType, sc)
	if !ok {
		return
	}
	deleted, err := h.repo.DeleteClassifier(c.Request.Context(), sc.auth, urlType, cid)
	if err != nil {
		HandlePgxError(c, err)
		return
	}
	if !deleted {
		commonHandlers.RespondError(c, http.StatusNotFound, "not found")
		return
	}
	sc.invalidate(c.Request.Context(), c)
	c.Status(http.StatusNoContent)
}

func (h *Handler) doRestoreClassifier(c *gin.Context, urlType string, sc *classifierScope) {
	cid, ok := h.requireCidInScope(c, urlType, sc)
	if !ok {
		return
	}
	result, err := h.repo.RestoreClassifier(c.Request.Context(), sc.auth, urlType, cid)
	if err != nil {
		HandlePgxError(c, err)
		return
	}
	if result == nil || string(result) == "null" {
		commonHandlers.RespondError(c, http.StatusNotFound, "not found or already active")
		return
	}
	sc.invalidate(c.Request.Context(), c)
	c.Data(http.StatusOK, "application/json", result)
}

// requireCidInScope parses :cid and verifies the classifier belongs to the
// request's book/hero scope. Keeps mutation + cache invalidation on the
// scope the caller actually named, even when the user owns both sides.
func (h *Handler) requireCidInScope(c *gin.Context, urlType string, sc *classifierScope) (int64, bool) {
	cid, err := getPathParamInt64(c, "cid")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return 0, false
	}
	match, err := h.repo.IsClassifierInScope(c.Request.Context(), sc.auth, urlType, cid, sc.sourceBookID, sc.heroID)
	if err != nil {
		HandlePgxError(c, err)
		return 0, false
	}
	if !match {
		commonHandlers.RespondError(c, http.StatusNotFound, "not found")
		return 0, false
	}
	return cid, true
}

// ----------------------------------------------------------------------------
// Source-book handlers
// ----------------------------------------------------------------------------

// CreateSourceBook handles POST /api/v1/homebrew/source-books.
func (h *Handler) CreateSourceBook(c *gin.Context) { h.upsertSourceBook(c, "") }

// UpdateSourceBook handles PUT /api/v1/homebrew/source-books/:code.
func (h *Handler) UpdateSourceBook(c *gin.Context) { h.upsertSourceBook(c, c.Param("code")) }

func (h *Handler) upsertSourceBook(c *gin.Context, code string) {
	auth, ok := h.requireAuth(c)
	if !ok {
		return
	}
	body, ok := readJSONBody(c)
	if !ok {
		return
	}
	if code != "" {
		merged, err := mergeCode(body, code)
		if err != nil {
			commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		body = merged
	}

	result, err := h.repo.UpsertSourceBook(c.Request.Context(), auth, body)
	if err != nil {
		HandlePgxError(c, err)
		return
	}
	if id, parsed := parseSourceBookID(result); parsed {
		h.invalidateBookCache(c.Request.Context(), c, id)
	}
	c.Data(http.StatusOK, "application/json", result)
}

// GetSourceBookByCode handles GET /api/v1/homebrew/source-books/:code.
func (h *Handler) GetSourceBookByCode(c *gin.Context) {
	handleGetByString(c, "code", h.repo.GetSourceBookByCode)
}

// ListMyHomebrewSourceBooks handles GET /api/v1/homebrew/source-books.
// Returns the session user's own homebrew books, including inactive and soft-deleted rows.
func (h *Handler) ListMyHomebrewSourceBooks(c *gin.Context) {
	handleGet(c, h.repo.ListMyHomebrewSourceBooks)
}

// DeleteSourceBook handles DELETE /api/v1/homebrew/source-books/:code (soft).
// Returns 204 No Content on success (HTTP DELETE convention; frontend refetches).
func (h *Handler) DeleteSourceBook(c *gin.Context) {
	auth, ok := h.requireAuth(c)
	if !ok {
		return
	}
	code := c.Param("code")
	if code == "" {
		commonHandlers.RespondError(c, http.StatusBadRequest, "missing path parameter: code")
		return
	}
	bookID, ok := h.resolveSourceBookIDByCode(c, auth, code)
	if !ok {
		return
	}
	changed, err := h.repo.DeleteSourceBookByCode(c.Request.Context(), auth, code)
	if err != nil {
		HandlePgxError(c, err)
		return
	}
	if !changed {
		commonHandlers.RespondError(c, http.StatusNotFound, "not found")
		return
	}
	h.invalidateBookCache(c.Request.Context(), c, bookID)
	c.Status(http.StatusNoContent)
}

// RestoreSourceBook handles POST /api/v1/homebrew/source-books/:code/restore.
// Returns 200 with the restored book JSONB so the client can update its
// list/currentBook state without a second GET. Matches the upsert_* return
// convention.
func (h *Handler) RestoreSourceBook(c *gin.Context) {
	auth, ok := h.requireAuth(c)
	if !ok {
		return
	}
	code := c.Param("code")
	if code == "" {
		commonHandlers.RespondError(c, http.StatusBadRequest, "missing path parameter: code")
		return
	}
	bookID, ok := h.resolveSourceBookIDByCode(c, auth, code)
	if !ok {
		return
	}
	result, err := h.repo.RestoreSourceBookByCode(c.Request.Context(), auth, code)
	if err != nil {
		HandlePgxError(c, err)
		return
	}
	if result == nil || string(result) == "null" {
		commonHandlers.RespondError(c, http.StatusNotFound, "not found or already active")
		return
	}
	h.invalidateBookCache(c.Request.Context(), c, bookID)
	c.Data(http.StatusOK, "application/json", result)
}

// ----------------------------------------------------------------------------
// Book-scoped classifier handlers (thin wrappers around the workers)
// ----------------------------------------------------------------------------

// UpsertBookClassifier handles POST/PUT /api/v1/homebrew/source-books/:code/:type[/:cid].
func (h *Handler) UpsertBookClassifier(c *gin.Context) {
	urlType, sc, ok := h.bookPreamble(c)
	if !ok {
		return
	}
	h.doUpsertClassifier(c, urlType, sc)
}

// DeleteBookClassifier handles DELETE /api/v1/homebrew/source-books/:code/:type/:cid.
func (h *Handler) DeleteBookClassifier(c *gin.Context) {
	urlType, sc, ok := h.bookPreamble(c)
	if !ok {
		return
	}
	h.doDeleteClassifier(c, urlType, sc)
}

// RestoreBookClassifier handles POST /api/v1/homebrew/source-books/:code/:type/:cid/restore.
func (h *Handler) RestoreBookClassifier(c *gin.Context) {
	urlType, sc, ok := h.bookPreamble(c)
	if !ok {
		return
	}
	h.doRestoreClassifier(c, urlType, sc)
}

// ----------------------------------------------------------------------------
// Hero-scoped classifier handlers (Phase 4.1 reuse, wired now)
// ----------------------------------------------------------------------------

// UpsertHeroClassifier handles POST/PUT /api/v1/homebrew/heroes/:id/:type[/:cid].
func (h *Handler) UpsertHeroClassifier(c *gin.Context) {
	urlType, sc, ok := h.heroPreamble(c)
	if !ok {
		return
	}
	h.doUpsertClassifier(c, urlType, sc)
}

// DeleteHeroClassifier handles DELETE /api/v1/homebrew/heroes/:id/:type/:cid.
func (h *Handler) DeleteHeroClassifier(c *gin.Context) {
	urlType, sc, ok := h.heroPreamble(c)
	if !ok {
		return
	}
	h.doDeleteClassifier(c, urlType, sc)
}

// RestoreHeroClassifier handles POST /api/v1/homebrew/heroes/:id/:type/:cid/restore.
func (h *Handler) RestoreHeroClassifier(c *gin.Context) {
	urlType, sc, ok := h.heroPreamble(c)
	if !ok {
		return
	}
	h.doRestoreClassifier(c, urlType, sc)
}

// ----------------------------------------------------------------------------
// Internal utilities
// ----------------------------------------------------------------------------

// injectScope merges scope fields (sourceBookId or heroId, exactly one) into
// the top-level JSON object, overwriting any client-supplied value.
func injectScope(raw json.RawMessage, sourceBookID, heroID *int64) (json.RawMessage, error) {
	if (sourceBookID == nil) == (heroID == nil) {
		return nil, fmt.Errorf("injectScope: exactly one of sourceBookID or heroID must be set")
	}

	obj := map[string]json.RawMessage{}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &obj); err != nil {
			return nil, fmt.Errorf("payload must be a JSON object: %w", err)
		}
	}

	if sourceBookID != nil {
		v, _ := json.Marshal(*sourceBookID)
		obj["sourceBookId"] = v
		obj["heroId"] = json.RawMessage("null")
	} else {
		v, _ := json.Marshal(*heroID)
		obj["heroId"] = v
		obj["sourceBookId"] = json.RawMessage("null")
	}

	return json.Marshal(obj)
}

// mergeCode injects the `code` field into a top-level JSON object.
func mergeCode(raw json.RawMessage, code string) (json.RawMessage, error) {
	var obj map[string]json.RawMessage
	if len(raw) == 0 {
		obj = map[string]json.RawMessage{}
	} else if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, fmt.Errorf("payload must be a JSON object: %w", err)
	}
	if obj == nil {
		obj = map[string]json.RawMessage{}
	}
	v, _ := json.Marshal(code)
	obj["code"] = v
	return json.Marshal(obj)
}

// parseSourceBookID extracts the integer id from an upsert_source_book result.
// Returns ok=false when the payload is empty/null or cannot be unmarshaled.
func parseSourceBookID(raw json.RawMessage) (int64, bool) {
	if len(raw) == 0 || string(raw) == "null" {
		return 0, false
	}
	var book struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(raw, &book); err != nil {
		return 0, false
	}
	return book.ID, true
}
