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
			commonHandlers.LogAndRespondError(c, http.StatusConflict, err, "resource already exists")
			return
		case "23503": // foreign_key_violation
			commonHandlers.LogAndRespondError(c, http.StatusBadRequest, err, "referenced resource not found")
			return
		case "23514": // check_violation
			commonHandlers.LogAndRespondError(c, http.StatusBadRequest, err, "validation constraint failed")
			return
		case "P0002": // no_data_found
			commonHandlers.LogAndRespondError(c, http.StatusNotFound, err, "not found")
			return
		case "42501": // insufficient_privilege
			commonHandlers.LogAndRespondError(c, http.StatusForbidden, err, "access denied")
			return
		case "22023": // invalid_parameter_value
			commonHandlers.LogAndRespondError(c, http.StatusBadRequest, err, "invalid parameter value")
			return
		case "22001": // string_data_right_truncation
			commonHandlers.LogAndRespondError(c, http.StatusBadRequest, err, "value too long")
			return
		case "P0001": // raise_exception (validation errors)
			commonHandlers.LogAndRespondError(c, http.StatusBadRequest, err, pgErr.Message)
			return
		case "55000": // object_not_in_prerequisite_state (e.g. soft-deleted row)
			commonHandlers.LogAndRespondError(c, http.StatusConflict, err, pgErr.Message)
			return
		}
	}

	commonHandlers.LogAndRespondError(c, http.StatusInternalServerError, err, "internal server error")
}
