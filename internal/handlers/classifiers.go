package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/GunarsK-rpg/public-api/internal/models/requests"
)

// Batch getter (all classifiers in one call)

func (h *Handler) GetAllClassifiers(c *gin.Context) {
	handleGet(c, h.repo.GetAllClassifiers)
}

// Simple getters (no parameters)

func (h *Handler) GetAttributeTypes(c *gin.Context) {
	handleGet(c, h.repo.GetAttributeTypes)
}

func (h *Handler) GetAttributes(c *gin.Context) {
	handleGet(c, h.repo.GetAttributes)
}

func (h *Handler) GetDerivedStats(c *gin.Context) {
	handleGet(c, h.repo.GetDerivedStats)
}

func (h *Handler) GetDerivedStatValues(c *gin.Context) {
	handleGet(c, h.repo.GetDerivedStatValues)
}

func (h *Handler) GetSkills(c *gin.Context) {
	handleGet(c, h.repo.GetSkills)
}

func (h *Handler) GetExpertiseTypes(c *gin.Context) {
	handleGet(c, h.repo.GetExpertiseTypes)
}

func (h *Handler) GetPaths(c *gin.Context) {
	handleGet(c, h.repo.GetPaths)
}

func (h *Handler) GetSurges(c *gin.Context) {
	handleGet(c, h.repo.GetSurges)
}

func (h *Handler) GetRadiantOrders(c *gin.Context) {
	handleGet(c, h.repo.GetRadiantOrders)
}

func (h *Handler) GetAncestries(c *gin.Context) {
	handleGet(c, h.repo.GetAncestries)
}

func (h *Handler) GetActivationTypes(c *gin.Context) {
	handleGet(c, h.repo.GetActivationTypes)
}

func (h *Handler) GetActionTypes(c *gin.Context) {
	handleGet(c, h.repo.GetActionTypes)
}

func (h *Handler) GetDamageTypes(c *gin.Context) {
	handleGet(c, h.repo.GetDamageTypes)
}

func (h *Handler) GetUnits(c *gin.Context) {
	handleGet(c, h.repo.GetUnits)
}

func (h *Handler) GetEquipmentTypes(c *gin.Context) {
	handleGet(c, h.repo.GetEquipmentTypes)
}

func (h *Handler) GetEquipmentAttributes(c *gin.Context) {
	handleGet(c, h.repo.GetEquipmentAttributes)
}

func (h *Handler) GetConditions(c *gin.Context) {
	handleGet(c, h.repo.GetConditions)
}

func (h *Handler) GetInjuries(c *gin.Context) {
	handleGet(c, h.repo.GetInjuries)
}

func (h *Handler) GetGoalStatus(c *gin.Context) {
	handleGet(c, h.repo.GetGoalStatus)
}

func (h *Handler) GetConnectionTypes(c *gin.Context) {
	handleGet(c, h.repo.GetConnectionTypes)
}

func (h *Handler) GetCompanionTypes(c *gin.Context) {
	handleGet(c, h.repo.GetCompanionTypes)
}

func (h *Handler) GetCultures(c *gin.Context) {
	handleGet(c, h.repo.GetCultures)
}

func (h *Handler) GetTiers(c *gin.Context) {
	handleGet(c, h.repo.GetTiers)
}

func (h *Handler) GetLevels(c *gin.Context) {
	handleGet(c, h.repo.GetLevels)
}

func (h *Handler) GetStartingKits(c *gin.Context) {
	handleGet(c, h.repo.GetStartingKits)
}

// Filtered getters (query params → JSONB filter)

func (h *Handler) GetExpertises(c *gin.Context) {
	handleGetFiltered[requests.GetExpertisesQuery](c, h.repo.GetExpertises)
}

func (h *Handler) GetSpecialties(c *gin.Context) {
	handleGetFiltered[requests.GetSpecialtiesQuery](c, h.repo.GetSpecialties)
}

func (h *Handler) GetSingerForms(c *gin.Context) {
	handleGetFiltered[requests.GetSingerFormsQuery](c, h.repo.GetSingerForms)
}

func (h *Handler) GetTalents(c *gin.Context) {
	handleGetFiltered[requests.GetTalentsQuery](c, h.repo.GetTalents)
}

func (h *Handler) GetActions(c *gin.Context) {
	handleGetFiltered[requests.GetActionsQuery](c, h.repo.GetActions)
}

func (h *Handler) GetActionLinks(c *gin.Context) {
	handleGetFiltered[requests.GetActionLinksQuery](c, h.repo.GetActionLinks)
}

func (h *Handler) GetEquipments(c *gin.Context) {
	handleGetFiltered[requests.GetEquipmentsQuery](c, h.repo.GetEquipments)
}
