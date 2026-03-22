package handlers

import (
	"net/http"

	commonHandlers "github.com/GunarsK-portfolio/portfolio-common/handlers"
	"github.com/gin-gonic/gin"
)

type setAvatarRequest struct {
	AvatarKey string `json:"avatarKey" binding:"required"`
}

// SetHeroAvatar sets the avatar key on a hero.
func (h *Handler) SetHeroAvatar(c *gin.Context) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	heroID, err := getPathParamInt64(c, "id")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	var req setAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, "avatarKey is required")
		return
	}

	if err := h.repo.UpsertHeroAvatar(c.Request.Context(), auth, heroID, req.AvatarKey); err != nil {
		HandlePgxError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatarKey": req.AvatarKey})
}

// DeleteHeroAvatar clears the avatar key on a hero.
func (h *Handler) DeleteHeroAvatar(c *gin.Context) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	heroID, err := getPathParamInt64(c, "id")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.DeleteHeroAvatar(c.Request.Context(), auth, heroID); err != nil {
		HandlePgxError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
