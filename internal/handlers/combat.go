package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	commonHandlers "github.com/GunarsK-portfolio/portfolio-common/handlers"
)

// NPCs (templates)

// GetNpcOptions returns lightweight NPC list for picker.
func (h *Handler) GetNpcOptions(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetNpcOptions)
}

// GetNpcLibrary returns all campaign NPCs including archived for library management.
func (h *Handler) GetNpcLibrary(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetNpcLibrary)
}

// GetNpc returns a full NPC stat block.
func (h *Handler) GetNpc(c *gin.Context) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	npcID, err := getPathParamInt64(c, "nid")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	campaignID, err := getPathParamInt64(c, "id")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.repo.GetNpc(c.Request.Context(), auth, npcID, campaignID)
	if err != nil {
		HandlePgxError(c, err)
		return
	}

	if result == nil || string(result) == "null" {
		commonHandlers.RespondError(c, http.StatusNotFound, "not found")
		return
	}

	c.Data(http.StatusOK, "application/json", result)
}

// GetNpcById returns a full NPC stat block by ID (no campaign scoping).
func (h *Handler) GetNpcById(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetNpcById)
}

// CreateNpc creates a custom NPC.
func (h *Handler) CreateNpc(c *gin.Context) {
	handlePost(c, h.repo.UpsertNpc)
}

// UpdateNpc updates a custom NPC.
func (h *Handler) UpdateNpc(c *gin.Context) {
	handlePost(c, h.repo.UpsertNpc)
}

// DeleteNpc soft-deletes a custom NPC.
func (h *Handler) DeleteNpc(c *gin.Context) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	npcID, err := getPathParamInt64(c, "nid")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	campaignID, err := getPathParamInt64(c, "id")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	deleted, err := h.repo.DeleteNpc(c.Request.Context(), auth, npcID, campaignID)
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

// Combats

// GetCombats returns all combats for a campaign.
func (h *Handler) GetCombats(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetCombats)
}

// GetCombat returns a single combat with NPC instances.
func (h *Handler) GetCombat(c *gin.Context) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	combatID, err := getPathParamInt64(c, "cid")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	campaignID, err := getPathParamInt64(c, "id")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.repo.GetCombat(c.Request.Context(), auth, combatID, campaignID)
	if err != nil {
		HandlePgxError(c, err)
		return
	}

	if result == nil || string(result) == "null" {
		commonHandlers.RespondError(c, http.StatusNotFound, "not found")
		return
	}

	c.Data(http.StatusOK, "application/json", result)
}

// CreateCombat creates a new combat encounter.
func (h *Handler) CreateCombat(c *gin.Context) {
	handlePost(c, h.repo.UpsertCombat)
}

// UpdateCombat updates a combat encounter.
func (h *Handler) UpdateCombat(c *gin.Context) {
	handlePost(c, h.repo.UpsertCombat)
}

// DeleteCombat deletes a combat encounter.
func (h *Handler) DeleteCombat(c *gin.Context) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	combatID, err := getPathParamInt64(c, "cid")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	campaignID, err := getPathParamInt64(c, "id")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	deleted, err := h.repo.DeleteCombat(c.Request.Context(), auth, combatID, campaignID)
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

// Combat round management

// EndCombatRound ends the current round — increments counter, resets turn speeds.
func (h *Handler) EndCombatRound(c *gin.Context) {
	handlePost(c, h.repo.EndCombatRound)
}

// NPC instances (combat + companion)

// GetNpcInstance returns a single NPC instance by ID.
func (h *Handler) GetNpcInstance(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetNpcInstance)
}

// CreateNpcInstance creates a combat or companion NPC instance.
func (h *Handler) CreateNpcInstance(c *gin.Context) {
	handlePost(c, h.repo.CreateNpcInstance)
}

// PatchNpcInstance updates an NPC instance (metadata or resource).
func (h *Handler) PatchNpcInstance(c *gin.Context) {
	handlePost(c, h.repo.PatchNpcInstance)
}

// DeleteNpcInstance removes an NPC instance.
func (h *Handler) DeleteNpcInstance(c *gin.Context) {
	handleDelete(c, "id", h.repo.DeleteNpcInstance)
}

// GetHeroCompanions returns all companion NPC instances for a hero.
func (h *Handler) GetHeroCompanions(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroNpcInstances)
}

// GetCompanionNpcOptions returns companion-eligible NPCs for the add companion picker.
func (h *Handler) GetCompanionNpcOptions(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetCompanionNpcOptions)
}
