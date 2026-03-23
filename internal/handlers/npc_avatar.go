package handlers

import (
	"net/http"

	commonHandlers "github.com/GunarsK-portfolio/portfolio-common/handlers"
	"github.com/gin-gonic/gin"
)

// SetNpcAvatar sets the avatar key on a custom NPC.
func (h *Handler) SetNpcAvatar(c *gin.Context) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	campaignID, err := getPathParamInt64(c, "id")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	npcID, err := getPathParamInt64(c, "nid")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	var req setAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.UpsertNpcAvatar(c.Request.Context(), auth, npcID, campaignID, req.AvatarKey); err != nil {
		HandlePgxError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatarKey": req.AvatarKey})
}

// DeleteNpcAvatar clears the avatar key on a custom NPC.
func (h *Handler) DeleteNpcAvatar(c *gin.Context) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	campaignID, err := getPathParamInt64(c, "id")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	npcID, err := getPathParamInt64(c, "nid")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.DeleteNpcAvatar(c.Request.Context(), auth, npcID, campaignID); err != nil {
		HandlePgxError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
