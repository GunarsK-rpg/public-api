package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	commonHandlers "github.com/GunarsK-portfolio/portfolio-common/handlers"
)

// HandlePgxError maps pgx/PostgreSQL errors to appropriate HTTP responses.
func HandlePgxError(c *gin.Context, err error, notFoundMsg string) {
	if errors.Is(err, pgx.ErrNoRows) {
		commonHandlers.RespondError(c, http.StatusNotFound, notFoundMsg)
		return
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			commonHandlers.RespondError(c, http.StatusConflict, "Resource already exists")
			return
		case "23503": // foreign_key_violation
			commonHandlers.RespondError(c, http.StatusBadRequest, "Referenced resource not found")
			return
		case "23514": // check_violation
			commonHandlers.RespondError(c, http.StatusBadRequest, "Validation constraint failed")
			return
		case "P0002": // no_data_found
			commonHandlers.RespondError(c, http.StatusNotFound, pgErr.Message)
			return
		case "42501": // insufficient_privilege
			commonHandlers.RespondError(c, http.StatusForbidden, pgErr.Message)
			return
		case "P0001": // raise_exception (validation errors)
			commonHandlers.RespondError(c, http.StatusBadRequest, pgErr.Message)
			return
		}
	}

	commonHandlers.LogAndRespondError(c, http.StatusInternalServerError, err, "Internal server error")
}

// HandleNullResult checks if a JSONB result is NULL (not found) and responds 404.
// Returns true if the result was null (response already sent).
func HandleNullResult(c *gin.Context, result json.RawMessage, notFoundMsg string) bool {
	if result == nil || string(result) == "null" {
		commonHandlers.RespondError(c, http.StatusNotFound, notFoundMsg)
		return true
	}
	return false
}
