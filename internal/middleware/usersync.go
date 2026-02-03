package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

// UserSync returns middleware that syncs the authenticated user to the RPG database.
// Must be applied after auth middleware (requires user_id and username in context).
func UserSync(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, err := handlers.GetAuthContext(c)
		if err != nil {
			slog.Error("failed to get auth context", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Sync user to RPG database
		_, err = pool.Exec(
			context.Background(),
			"SELECT auth.sync_user($1, $2)",
			auth.UserID,
			auth.Username,
		)
		if err != nil {
			slog.Error("failed to sync user",
				"user_id", auth.UserID,
				"username", auth.Username,
				"error", err,
			)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.Next()
	}
}
