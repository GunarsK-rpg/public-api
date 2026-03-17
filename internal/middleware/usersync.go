package middleware

import (
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

		// Sync user to RPG database (empty display_name passed as NULL)
		var displayName *string
		if auth.DisplayName != "" {
			displayName = &auth.DisplayName
		}
		_, err = pool.Exec(
			c.Request.Context(),
			"SELECT auth.sync_user($1, $2, $3)",
			auth.UserID,
			auth.Username,
			displayName,
		)
		if err != nil {
			slog.Error("failed to sync user",
				"user_id", auth.UserID,
				"error", err,
			)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.Next()
	}
}
