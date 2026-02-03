package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	commonHandlers "github.com/GunarsK-portfolio/portfolio-common/handlers"

	"github.com/GunarsK-rpg/public-api/internal/repository"
)

// ErrMissingAuthContext indicates auth middleware did not run or set context values.
var ErrMissingAuthContext = errors.New("missing auth context: user_id or username not set")

// GetAuthContext extracts authentication context from Gin context.
// Returns error if auth middleware has not set the required values.
func GetAuthContext(c *gin.Context) (repository.AuthContext, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return repository.AuthContext{}, ErrMissingAuthContext
	}

	username, exists := c.Get("username")
	if !exists {
		return repository.AuthContext{}, ErrMissingAuthContext
	}

	uid, ok := userID.(int64)
	if !ok {
		return repository.AuthContext{}, ErrMissingAuthContext
	}

	uname, ok := username.(string)
	if !ok {
		return repository.AuthContext{}, ErrMissingAuthContext
	}

	return repository.AuthContext{
		UserID:    uid,
		Username:  uname,
		ClientIP:  c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}, nil
}

// RepoFunc is a function that calls a repository method and returns JSONB.
type RepoFunc func(ctx context.Context, auth repository.AuthContext) (json.RawMessage, error)

// handleJSONResponse handles the common pattern: auth → repo call → error handling → JSON response.
func handleJSONResponse(c *gin.Context, fn RepoFunc) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := fn(c.Request.Context(), auth)
	if err != nil {
		HandlePgxError(c, err)
		return
	}

	c.Data(http.StatusOK, "application/json", result)
}

// RepoFilterFunc is a function that calls a repository method with a JSONB filter and returns JSONB.
type RepoFilterFunc func(ctx context.Context, auth repository.AuthContext, filter json.RawMessage) (json.RawMessage, error)

// handleFilteredResponse binds query params, marshals to JSON, and calls a filtered repo method.
func handleFilteredResponse[Q any](c *gin.Context, fn RepoFilterFunc) {
	var query Q
	if err := c.ShouldBindQuery(&query); err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	filter, err := json.Marshal(query)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusInternalServerError, "failed to build filter")
		return
	}

	handleJSONResponse(c, func(ctx context.Context, auth repository.AuthContext) (json.RawMessage, error) {
		return fn(ctx, auth, filter)
	})
}

// Validatable is implemented by query structs that require custom validation.
type Validatable interface {
	Validate() error
}

// handleValidatedFilteredResponse is like handleFilteredResponse but calls Validate() on the query.
func handleValidatedFilteredResponse[Q Validatable](c *gin.Context, fn RepoFilterFunc) {
	var query Q
	if err := c.ShouldBindQuery(&query); err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	if err := query.Validate(); err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	filter, err := json.Marshal(query)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusInternalServerError, "failed to build filter")
		return
	}

	handleJSONResponse(c, func(ctx context.Context, auth repository.AuthContext) (json.RawMessage, error) {
		return fn(ctx, auth, filter)
	})
}
