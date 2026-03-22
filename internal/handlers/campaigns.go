package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	commonHandlers "github.com/GunarsK-portfolio/portfolio-common/handlers"
)

// GetCampaigns returns a list of campaigns.
func (h *Handler) GetCampaigns(c *gin.Context) {
	handleGet(c, h.repo.GetCampaigns)
}

// GetCampaign returns a single campaign by ID.
func (h *Handler) GetCampaign(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetCampaign)
}

// GetCampaignByCode returns a campaign by invite code.
func (h *Handler) GetCampaignByCode(c *gin.Context) {
	handleGetByString(c, "code", h.repo.GetCampaignByCode)
}

// CreateCampaign creates a new campaign.
func (h *Handler) CreateCampaign(c *gin.Context) {
	handlePost(c, h.repo.UpsertCampaign)
}

// UpdateCampaign updates an existing campaign.
func (h *Handler) UpdateCampaign(c *gin.Context) {
	handlePost(c, h.repo.UpsertCampaign)
}

// DeleteCampaign deletes a campaign by ID.
func (h *Handler) DeleteCampaign(c *gin.Context) {
	handleDelete(c, "id", h.repo.DeleteCampaign)
}

// RemoveHeroFromCampaign removes a hero from a campaign (owner only).
func (h *Handler) RemoveHeroFromCampaign(c *gin.Context) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	heroID, err := getPathParamInt64(c, "hid")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	campaignID, err := getPathParamInt64(c, "id")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	deleted, err := h.repo.RemoveHeroFromCampaign(c.Request.Context(), auth, heroID, campaignID)
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
