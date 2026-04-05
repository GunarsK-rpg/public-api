package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	commonHandlers "github.com/GunarsK-portfolio/portfolio-common/handlers"

	"github.com/GunarsK-rpg/public-api/internal/models/requests"
)

// GetHeroes returns a list of heroes filtered by campaign.
func (h *Handler) GetHeroes(c *gin.Context) {
	auth, err := GetAuthContext(c)
	if err != nil {
		commonHandlers.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var query requests.GetHeroesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.Error(err)
		commonHandlers.RespondError(c, http.StatusBadRequest, "invalid query parameters")
		return
	}
	if err := query.Validate(); err != nil {
		commonHandlers.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.repo.GetHeroes(c.Request.Context(), auth, query.CampaignID)
	if err != nil {
		HandlePgxError(c, err)
		return
	}

	c.Data(http.StatusOK, "application/json", result)
}

// GetHero returns a single hero by ID.
func (h *Handler) GetHero(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHero)
}

// GetHeroSheet returns complete hero sheet by ID.
func (h *Handler) GetHeroSheet(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroSheet)
}

// CreateHero creates a new hero.
func (h *Handler) CreateHero(c *gin.Context) {
	handlePost(c, h.repo.UpsertHero)
}

// UpdateHero updates an existing hero.
func (h *Handler) UpdateHero(c *gin.Context) {
	handlePost(c, h.repo.UpsertHero)
}

// DeleteHero deletes a hero by ID.
func (h *Handler) DeleteHero(c *gin.Context) {
	handleDelete(c, "id", h.repo.DeleteHero)
}

// Attributes

// GetHeroAttributes returns all attributes for a hero.
func (h *Handler) GetHeroAttributes(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroAttributes)
}

// UpsertHeroAttribute creates or updates a hero attribute.
func (h *Handler) UpsertHeroAttribute(c *gin.Context) {
	handlePost(c, h.repo.UpsertHeroAttribute)
}

// DeleteHeroAttribute deletes a hero attribute.
func (h *Handler) DeleteHeroAttribute(c *gin.Context) {
	handleDelete(c, "subId", h.repo.DeleteHeroAttribute)
}

// Defenses

// GetHeroDefenses returns all defenses for a hero.
func (h *Handler) GetHeroDefenses(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroDefenses)
}

// UpsertHeroDefense creates or updates a hero defense.
func (h *Handler) UpsertHeroDefense(c *gin.Context) {
	handlePost(c, h.repo.UpsertHeroDefense)
}

// DeleteHeroDefense deletes a hero defense.
func (h *Handler) DeleteHeroDefense(c *gin.Context) {
	handleDelete(c, "subId", h.repo.DeleteHeroDefense)
}

// Derived Stats

// GetHeroDerivedStats returns all derived stats for a hero.
func (h *Handler) GetHeroDerivedStats(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroDerivedStats)
}

// UpsertHeroDerivedStat creates or updates a hero derived stat.
func (h *Handler) UpsertHeroDerivedStat(c *gin.Context) {
	handlePost(c, h.repo.UpsertHeroDerivedStat)
}

// DeleteHeroDerivedStat deletes a hero derived stat.
func (h *Handler) DeleteHeroDerivedStat(c *gin.Context) {
	handleDelete(c, "subId", h.repo.DeleteHeroDerivedStat)
}

// Skills

// GetHeroSkills returns all skills for a hero.
func (h *Handler) GetHeroSkills(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroSkills)
}

// UpsertHeroSkill creates or updates a hero skill.
func (h *Handler) UpsertHeroSkill(c *gin.Context) {
	handlePost(c, h.repo.UpsertHeroSkill)
}

// DeleteHeroSkill deletes a hero skill.
func (h *Handler) DeleteHeroSkill(c *gin.Context) {
	handleDelete(c, "subId", h.repo.DeleteHeroSkill)
}

// Expertises

// GetHeroExpertises returns all expertises for a hero.
func (h *Handler) GetHeroExpertises(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroExpertises)
}

// UpsertHeroExpertise creates or updates a hero expertise.
func (h *Handler) UpsertHeroExpertise(c *gin.Context) {
	handlePost(c, h.repo.UpsertHeroExpertise)
}

// DeleteHeroExpertise deletes a hero expertise.
func (h *Handler) DeleteHeroExpertise(c *gin.Context) {
	handleDelete(c, "subId", h.repo.DeleteHeroExpertise)
}

// Talents

// GetHeroTalents returns all talents for a hero.
func (h *Handler) GetHeroTalents(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroTalents)
}

// UpsertHeroTalent creates or updates a hero talent.
func (h *Handler) UpsertHeroTalent(c *gin.Context) {
	handlePost(c, h.repo.UpsertHeroTalent)
}

// DeleteHeroTalent deletes a hero talent.
func (h *Handler) DeleteHeroTalent(c *gin.Context) {
	handleDelete(c, "subId", h.repo.DeleteHeroTalent)
}

// Equipment

// GetHeroEquipment returns all equipment for a hero.
func (h *Handler) GetHeroEquipment(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroEquipment)
}

// UpsertHeroEquipment creates or updates a hero equipment.
func (h *Handler) UpsertHeroEquipment(c *gin.Context) {
	handlePost(c, h.repo.UpsertHeroEquipment)
}

// DeleteHeroEquipment deletes a hero equipment.
func (h *Handler) DeleteHeroEquipment(c *gin.Context) {
	handleDelete(c, "subId", h.repo.DeleteHeroEquipment)
}

// Equipment Modifications

// AddEquipmentModification adds a modification to hero equipment.
func (h *Handler) AddEquipmentModification(c *gin.Context) {
	handlePost(c, h.repo.AddEquipmentModification)
}

// RemoveEquipmentModification removes a modification from hero equipment.
func (h *Handler) RemoveEquipmentModification(c *gin.Context) {
	handleDelete(c, "modId", h.repo.RemoveEquipmentModification)
}

// Favorite actions

// AddFavoriteAction adds an action to hero favorites.
func (h *Handler) AddFavoriteAction(c *gin.Context) {
	handlePost(c, h.repo.AddFavoriteAction)
}

// RemoveFavoriteAction removes an action from hero favorites.
func (h *Handler) RemoveFavoriteAction(c *gin.Context) {
	handleDelete(c, "subId", h.repo.RemoveFavoriteAction)
}

// Conditions

// GetHeroConditions returns all conditions for a hero.
func (h *Handler) GetHeroConditions(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroConditions)
}

// UpsertHeroCondition creates or updates a hero condition.
func (h *Handler) UpsertHeroCondition(c *gin.Context) {
	handlePost(c, h.repo.UpsertHeroCondition)
}

// DeleteHeroCondition deletes a hero condition.
func (h *Handler) DeleteHeroCondition(c *gin.Context) {
	handleDelete(c, "subId", h.repo.DeleteHeroCondition)
}

// Injuries

// GetHeroInjuries returns all injuries for a hero.
func (h *Handler) GetHeroInjuries(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroInjuries)
}

// UpsertHeroInjury creates or updates a hero injury.
func (h *Handler) UpsertHeroInjury(c *gin.Context) {
	handlePost(c, h.repo.UpsertHeroInjury)
}

// DeleteHeroInjury deletes a hero injury.
func (h *Handler) DeleteHeroInjury(c *gin.Context) {
	handleDelete(c, "subId", h.repo.DeleteHeroInjury)
}

// Goals

// GetHeroGoals returns all goals for a hero.
func (h *Handler) GetHeroGoals(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroGoals)
}

// UpsertHeroGoal creates or updates a hero goal.
func (h *Handler) UpsertHeroGoal(c *gin.Context) {
	handlePost(c, h.repo.UpsertHeroGoal)
}

// DeleteHeroGoal deletes a hero goal.
func (h *Handler) DeleteHeroGoal(c *gin.Context) {
	handleDelete(c, "subId", h.repo.DeleteHeroGoal)
}

// Connections

// GetHeroConnections returns all connections for a hero.
func (h *Handler) GetHeroConnections(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroConnections)
}

// UpsertHeroConnection creates or updates a hero connection.
func (h *Handler) UpsertHeroConnection(c *gin.Context) {
	handlePost(c, h.repo.UpsertHeroConnection)
}

// DeleteHeroConnection deletes a hero connection.
func (h *Handler) DeleteHeroConnection(c *gin.Context) {
	handleDelete(c, "subId", h.repo.DeleteHeroConnection)
}

// Notes

// GetHeroNotes returns all notes for a hero.
func (h *Handler) GetHeroNotes(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroNotes)
}

// UpsertHeroNote creates or updates a hero note.
func (h *Handler) UpsertHeroNote(c *gin.Context) {
	handlePost(c, h.repo.UpsertHeroNote)
}

// DeleteHeroNote deletes a hero note.
func (h *Handler) DeleteHeroNote(c *gin.Context) {
	handleDelete(c, "subId", h.repo.DeleteHeroNote)
}

// Cultures

// GetHeroCultures returns all cultures for a hero.
func (h *Handler) GetHeroCultures(c *gin.Context) {
	handleGetByID(c, "id", h.repo.GetHeroCultures)
}

// UpsertHeroCulture creates or updates a hero culture.
func (h *Handler) UpsertHeroCulture(c *gin.Context) {
	handlePost(c, h.repo.UpsertHeroCulture)
}

// DeleteHeroCulture deletes a hero culture.
func (h *Handler) DeleteHeroCulture(c *gin.Context) {
	handleDelete(c, "subId", h.repo.DeleteHeroCulture)
}

// Resource patches

// PatchHeroHealth updates hero current health.
func (h *Handler) PatchHeroHealth(c *gin.Context) {
	handlePost(c, h.repo.PatchHeroHealth)
}

// PatchHeroFocus updates hero current focus.
func (h *Handler) PatchHeroFocus(c *gin.Context) {
	handlePost(c, h.repo.PatchHeroFocus)
}

// PatchHeroInvestiture updates hero current investiture.
func (h *Handler) PatchHeroInvestiture(c *gin.Context) {
	handlePost(c, h.repo.PatchHeroInvestiture)
}

// PatchHeroCurrency updates hero currency.
func (h *Handler) PatchHeroCurrency(c *gin.Context) {
	handlePost(c, h.repo.PatchHeroCurrency)
}
