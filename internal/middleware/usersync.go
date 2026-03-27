package middleware

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/GunarsK-portfolio/portfolio-common/logger"

	"github.com/GunarsK-rpg/public-api/internal/cache"
	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

const (
	userSyncTTL       = 15 * time.Minute
	userSyncKeyPrefix = "rpg:usersync:"
)

// DBExecer executes SQL statements. Implemented by *pgxpool.Pool.
type DBExecer interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// UserSync returns middleware that syncs the authenticated user to the RPG database.
// Uses Redis to skip redundant syncs within the TTL window.
// Must be applied after auth middleware (requires user_id and username in context).
func UserSync(pool DBExecer, appCache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, err := handlers.GetAuthContext(c)
		if err != nil {
			logger.GetLogger(c).Error("failed to get auth context", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		ctx := c.Request.Context()
		claimsHash := fmt.Sprintf("%x", sha256.Sum256([]byte(auth.Username+"\x00"+auth.DisplayName)))[:12]
		cacheKey := fmt.Sprintf("%s%d:%s", userSyncKeyPrefix, auth.UserID, claimsHash)

		synced, err := appCache.HasFlag(ctx, cacheKey)
		if err != nil {
			logger.GetLogger(c).Warn("usersync cache check failed, falling back to DB sync", "error", err)
		}
		if synced {
			c.Next()
			return
		}

		var displayName *string
		if auth.DisplayName != "" {
			displayName = &auth.DisplayName
		}
		_, err = pool.Exec(
			ctx,
			"SELECT auth.sync_user($1, $2, $3)",
			auth.UserID,
			auth.Username,
			displayName,
		)
		if err != nil {
			logger.GetLogger(c).Error("failed to sync user",
				"user_id", auth.UserID,
				"error", err,
			)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		if err := appCache.SetFlag(ctx, cacheKey, userSyncTTL); err != nil {
			logger.GetLogger(c).Warn("failed to set usersync cache flag", "user_id", auth.UserID, "error", err)
		}

		c.Next()
	}
}
