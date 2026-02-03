package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	commonHandlers "github.com/GunarsK-portfolio/portfolio-common/handlers"
)

// HandlePgxError maps pgx/PostgreSQL errors to appropriate HTTP responses.
func HandlePgxError(c *gin.Context, err error) {
	if errors.Is(err, pgx.ErrNoRows) {
		commonHandlers.RespondError(c, http.StatusNotFound, "not found")
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
