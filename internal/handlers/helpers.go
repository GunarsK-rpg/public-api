package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"

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

	return repository.AuthContext{
		UserID:    userID.(int64),
		Username:  username.(string),
		ClientIP:  c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}, nil
}
