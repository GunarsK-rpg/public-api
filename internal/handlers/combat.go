package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	commonHandlers "github.com/GunarsK-portfolio/portfolio-common/handlers"
)

// NPCs

// GetNpcOptions returns lightweight NPC list for picker.
func (h *Handler) GetNpcOptions(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetNpcOptions)
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

// Combat NPC instances

// CreateCombatNpc adds an NPC instance to a combat.
func (h *Handler) CreateCombatNpc(c *gin.Context) {
	handlePost(c, h.repo.UpsertCombatNpc)
}

// UpdateCombatNpc updates a combat NPC instance.
func (h *Handler) UpdateCombatNpc(c *gin.Context) {
	handlePost(c, h.repo.UpsertCombatNpc)
}

// DeleteCombatNpc removes an NPC instance from a combat.
func (h *Handler) DeleteCombatNpc(c *gin.Context) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	instanceID, err := getPathParamInt64(c, "iid")
	if err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
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

	deleted, err := h.repo.DeleteCombatNpc(c.Request.Context(), auth, instanceID, combatID, campaignID)
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

// Combat NPC resource patches

// PatchCombatNpcHp updates combat NPC current HP.
func (h *Handler) PatchCombatNpcHp(c *gin.Context) {
	handlePost(c, h.repo.PatchCombatNpcHp)
}

// PatchCombatNpcFocus updates combat NPC current focus.
func (h *Handler) PatchCombatNpcFocus(c *gin.Context) {
	handlePost(c, h.repo.PatchCombatNpcFocus)
}

// PatchCombatNpcInvestiture updates combat NPC current investiture.
func (h *Handler) PatchCombatNpcInvestiture(c *gin.Context) {
	handlePost(c, h.repo.PatchCombatNpcInvestiture)
}
