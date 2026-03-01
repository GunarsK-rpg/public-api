package handlers

import (
	"github.com/gin-gonic/gin"
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
