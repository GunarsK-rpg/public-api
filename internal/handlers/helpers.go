package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	commonHandlers "github.com/GunarsK-portfolio/portfolio-common/handlers"

	"github.com/GunarsK-rpg/public-api/internal/repository"
)

// ----------------------------------------------------------------------------
// Repository function types
// ----------------------------------------------------------------------------

// RepoFunc calls a repository method and returns JSONB.
type RepoFunc func(ctx context.Context, auth repository.AuthContext) (json.RawMessage, error)

// RepoIDFunc calls a repository method with an ID parameter.
type RepoIDFunc func(ctx context.Context, auth repository.AuthContext, id int64) (json.RawMessage, error)

// RepoUpsertFunc calls a repository upsert method with JSON data.
type RepoUpsertFunc func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)

// RepoStringFunc calls a repository method with a string parameter.
type RepoStringFunc func(ctx context.Context, auth repository.AuthContext, code string) (json.RawMessage, error)

// RepoDeleteFunc calls a repository delete method.
type RepoDeleteFunc func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)

// RepoDeleteStringFunc calls a repository delete method with a string parameter.
type RepoDeleteStringFunc func(ctx context.Context, auth repository.AuthContext, code string) (bool, error)

// RepoTwoIDFunc calls a repository method with two ID parameters.
type RepoTwoIDFunc func(ctx context.Context, auth repository.AuthContext, id1 int64, id2 int64) (json.RawMessage, error)

// RepoDeleteTwoIDFunc calls a repository delete method with two ID parameters.
type RepoDeleteTwoIDFunc func(ctx context.Context, auth repository.AuthContext, id1 int64, id2 int64) (bool, error)

// ----------------------------------------------------------------------------
// Auth context
// ----------------------------------------------------------------------------

// ErrMissingAuthContext indicates auth middleware did not run or set context values.
var ErrMissingAuthContext = errors.New("missing auth context: user_id or username not set")

// GetAuthContext extracts authentication context from Gin context.
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

	displayName, _ := c.Get("display_name")
	dname, _ := displayName.(string)

	return repository.AuthContext{
		UserID:      uid,
		Username:    uname,
		DisplayName: dname,
		ClientIP:    c.ClientIP(),
		UserAgent:   c.GetHeader("User-Agent"),
	}, nil
}

// ----------------------------------------------------------------------------
// Request handlers
// ----------------------------------------------------------------------------

// handleGet: auth → repo call → JSON response
func handleGet(c *gin.Context, fn RepoFunc) {
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

// handleGetByID: auth → path param → repo call → null check → JSON response
func handleGetByID(c *gin.Context, paramName string, fn RepoIDFunc) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := getPathParamInt64(c, paramName)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := fn(c.Request.Context(), auth, id)
	if err != nil {
		HandlePgxError(c, err)
		return
	}

	if result == nil || string(result) == "null" {
		commonHandlers.RespondError(c, http.StatusNotFound, "not found")
		return
	}

	c.Data(http.StatusOK, "application/json", result)
}

// handleGetByString: auth → string path param → repo call → null check → JSON response
func handleGetByString(c *gin.Context, paramName string, fn RepoStringFunc) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	value := c.Param(paramName)
	if value == "" {
		commonHandlers.RespondError(c, http.StatusBadRequest, fmt.Sprintf("missing path parameter: %s", paramName))
		return
	}

	result, err := fn(c.Request.Context(), auth, value)
	if err != nil {
		HandlePgxError(c, err)
		return
	}

	if result == nil || string(result) == "null" {
		commonHandlers.RespondError(c, http.StatusNotFound, "not found")
		return
	}

	c.Data(http.StatusOK, "application/json", result)
}

// handlePost: auth → read body → repo call → JSON response
func handlePost(c *gin.Context, fn RepoUpsertFunc) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, "failed to read request body")
		return
	}

	if !json.Valid(body) {
		commonHandlers.RespondError(c, http.StatusBadRequest, "invalid JSON")
		return
	}

	result, err := fn(c.Request.Context(), auth, body)
	if err != nil {
		HandlePgxError(c, err)
		return
	}

	c.Data(http.StatusOK, "application/json", result)
}

// handleDeleteCommon: auth → repo delete closure → 204 or 404. Callers parse
// their param and pass a closure that captures it.
func handleDeleteCommon(c *gin.Context, del func(ctx context.Context, auth repository.AuthContext) (bool, error)) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	deleted, err := del(c.Request.Context(), auth)
	if err != nil {
		HandlePgxError(c, err)
		return
	}

	if !deleted {
		commonHandlers.RespondError(c, http.StatusNotFound, "not found")
		return
	}

	c.Status(http.StatusNoContent)
}

// handleDeleteByString: auth → string path param → repo call → 204 or 404
func handleDeleteByString(c *gin.Context, paramName string, fn RepoDeleteStringFunc) {
	value := c.Param(paramName)
	if value == "" {
		commonHandlers.RespondError(c, http.StatusBadRequest, fmt.Sprintf("missing path parameter: %s", paramName))
		return
	}
	handleDeleteCommon(c, func(ctx context.Context, auth repository.AuthContext) (bool, error) {
		return fn(ctx, auth, value)
	})
}

// handleDelete: auth → path param → repo call → 204 or 404
func handleDelete(c *gin.Context, paramName string, fn RepoDeleteFunc) {
	id, err := getPathParamInt64(c, paramName)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	handleDeleteCommon(c, func(ctx context.Context, auth repository.AuthContext) (bool, error) {
		return fn(ctx, auth, id)
	})
}

// handleGetByTwoIDs: auth → two path params → repo call → null check → JSON response
func handleGetByTwoIDs(c *gin.Context, param1, param2 string, fn RepoTwoIDFunc) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id1, err := getPathParamInt64(c, param1)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	id2, err := getPathParamInt64(c, param2)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := fn(c.Request.Context(), auth, id1, id2)
	if err != nil {
		HandlePgxError(c, err)
		return
	}

	if result == nil || string(result) == "null" {
		commonHandlers.RespondError(c, http.StatusNotFound, "not found")
		return
	}

	c.Data(http.StatusOK, "application/json", result)
}

// handleDeleteByTwoIDs: auth → two path params → repo call → 204 or 404
func handleDeleteByTwoIDs(c *gin.Context, param1, param2 string, fn RepoDeleteTwoIDFunc) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id1, err := getPathParamInt64(c, param1)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	id2, err := getPathParamInt64(c, param2)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	deleted, err := fn(c.Request.Context(), auth, id1, id2)
	if err != nil {
		HandlePgxError(c, err)
		return
	}

	if !deleted {
		commonHandlers.RespondError(c, http.StatusNotFound, "not found")
		return
	}

	c.Status(http.StatusNoContent)
}

// ----------------------------------------------------------------------------
// Internal utilities
// ----------------------------------------------------------------------------

func getPathParamInt64(c *gin.Context, param string) (int64, error) {
	str := c.Param(param)
	if str == "" {
		return 0, fmt.Errorf("missing path parameter: %s", param)
	}
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: must be an integer", param)
	}
	return id, nil
}
