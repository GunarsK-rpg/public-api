package repository

import (
	"context"
	"encoding/json"
)

// HeroRepository defines methods for hero data access.
type HeroRepository interface {
	// Core CRUD
	GetHeroes(ctx context.Context, auth AuthContext, campaignID *int64) (json.RawMessage, error)
	GetHero(ctx context.Context, auth AuthContext, id int64) (json.RawMessage, error)
	GetHeroSheet(ctx context.Context, auth AuthContext, id int64) (json.RawMessage, error)
	UpsertHero(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	DeleteHero(ctx context.Context, auth AuthContext, id int64) (bool, error)

	// Sub-resource getters (13)
	GetHeroAttributes(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
	GetHeroDefenses(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
	GetHeroDerivedStats(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
	GetHeroSkills(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
	GetHeroExpertises(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
	GetHeroTalents(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
	GetHeroEquipment(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
	GetHeroConditions(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
	GetHeroInjuries(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
	GetHeroGoals(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
	GetHeroConnections(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
	GetHeroCompanions(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
	GetHeroCultures(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)

	// Sub-resource upserts (13)
	UpsertHeroAttribute(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	UpsertHeroDefense(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	UpsertHeroDerivedStat(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	UpsertHeroSkill(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	UpsertHeroExpertise(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	UpsertHeroTalent(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	UpsertHeroEquipment(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	UpsertHeroCondition(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	UpsertHeroInjury(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	UpsertHeroGoal(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	UpsertHeroConnection(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	UpsertHeroCompanion(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	UpsertHeroCulture(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)

	// Resource patches
	PatchHeroHealth(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	PatchHeroFocus(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	PatchHeroInvestiture(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	PatchHeroCurrency(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)

	// Equipment modification management
	AddEquipmentModification(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	RemoveEquipmentModification(ctx context.Context, auth AuthContext, id int64) (bool, error)

	// Sub-resource deletes (13)
	DeleteHeroAttribute(ctx context.Context, auth AuthContext, id int64) (bool, error)
	DeleteHeroDefense(ctx context.Context, auth AuthContext, id int64) (bool, error)
	DeleteHeroDerivedStat(ctx context.Context, auth AuthContext, id int64) (bool, error)
	DeleteHeroSkill(ctx context.Context, auth AuthContext, id int64) (bool, error)
	DeleteHeroExpertise(ctx context.Context, auth AuthContext, id int64) (bool, error)
	DeleteHeroTalent(ctx context.Context, auth AuthContext, id int64) (bool, error)
	DeleteHeroEquipment(ctx context.Context, auth AuthContext, id int64) (bool, error)
	DeleteHeroCondition(ctx context.Context, auth AuthContext, id int64) (bool, error)
	DeleteHeroInjury(ctx context.Context, auth AuthContext, id int64) (bool, error)
	DeleteHeroGoal(ctx context.Context, auth AuthContext, id int64) (bool, error)
	DeleteHeroConnection(ctx context.Context, auth AuthContext, id int64) (bool, error)
	DeleteHeroCompanion(ctx context.Context, auth AuthContext, id int64) (bool, error)
	DeleteHeroCulture(ctx context.Context, auth AuthContext, id int64) (bool, error)
}

// Core CRUD implementations

func (r *repository) GetHeroes(ctx context.Context, auth AuthContext, campaignID *int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_heroes($1)", campaignID)
}

func (r *repository) GetHero(ctx context.Context, auth AuthContext, id int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero($1)", id)
}

func (r *repository) GetHeroSheet(ctx context.Context, auth AuthContext, id int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_sheet($1)", id)
}

func (r *repository) UpsertHero(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero($1::jsonb)", data)
}

func (r *repository) DeleteHero(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero($1)", id)
}

// Sub-resource getters

func (r *repository) GetHeroAttributes(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_attributes($1)", heroID)
}

func (r *repository) GetHeroDefenses(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_defenses($1)", heroID)
}

func (r *repository) GetHeroDerivedStats(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_derived_stats($1)", heroID)
}

func (r *repository) GetHeroSkills(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_skills($1)", heroID)
}

func (r *repository) GetHeroExpertises(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_expertises($1)", heroID)
}

func (r *repository) GetHeroTalents(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_talents($1)", heroID)
}

func (r *repository) GetHeroEquipment(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_equipment($1)", heroID)
}

func (r *repository) GetHeroConditions(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_conditions($1)", heroID)
}

func (r *repository) GetHeroInjuries(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_injuries($1)", heroID)
}

func (r *repository) GetHeroGoals(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_goals($1)", heroID)
}

func (r *repository) GetHeroConnections(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_connections($1)", heroID)
}

func (r *repository) GetHeroCompanions(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_companions($1)", heroID)
}

func (r *repository) GetHeroCultures(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.get_hero_cultures($1)", heroID)
}

// Sub-resource upserts

func (r *repository) UpsertHeroAttribute(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_attribute($1::jsonb)", data)
}

func (r *repository) UpsertHeroDefense(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_defense($1::jsonb)", data)
}

func (r *repository) UpsertHeroDerivedStat(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_derived_stat($1::jsonb)", data)
}

func (r *repository) UpsertHeroSkill(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_skill($1::jsonb)", data)
}

func (r *repository) UpsertHeroExpertise(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_expertise($1::jsonb)", data)
}

func (r *repository) UpsertHeroTalent(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_talent($1::jsonb)", data)
}

func (r *repository) UpsertHeroEquipment(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_equipment($1::jsonb)", data)
}

func (r *repository) UpsertHeroCondition(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_condition($1::jsonb)", data)
}

func (r *repository) UpsertHeroInjury(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_injury($1::jsonb)", data)
}

func (r *repository) UpsertHeroGoal(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_goal($1::jsonb)", data)
}

func (r *repository) UpsertHeroConnection(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_connection($1::jsonb)", data)
}

func (r *repository) UpsertHeroCompanion(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_companion($1::jsonb)", data)
}

func (r *repository) UpsertHeroCulture(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_culture($1::jsonb)", data)
}

// Resource patches

func (r *repository) PatchHeroHealth(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.patch_hero_health($1::jsonb)", data)
}

func (r *repository) PatchHeroFocus(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.patch_hero_focus($1::jsonb)", data)
}

func (r *repository) PatchHeroInvestiture(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.patch_hero_investiture($1::jsonb)", data)
}

func (r *repository) PatchHeroCurrency(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.patch_hero_currency($1::jsonb)", data)
}

// Equipment modification management

func (r *repository) AddEquipmentModification(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT heroes.upsert_hero_equipment_modification($1::jsonb)", data)
}

func (r *repository) RemoveEquipmentModification(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_equipment_modification($1)", id)
}

// Sub-resource deletes

func (r *repository) DeleteHeroAttribute(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_attribute($1)", id)
}

func (r *repository) DeleteHeroDefense(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_defense($1)", id)
}

func (r *repository) DeleteHeroDerivedStat(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_derived_stat($1)", id)
}

func (r *repository) DeleteHeroSkill(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_skill($1)", id)
}

func (r *repository) DeleteHeroExpertise(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_expertise($1)", id)
}

func (r *repository) DeleteHeroTalent(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_talent($1)", id)
}

func (r *repository) DeleteHeroEquipment(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_equipment($1)", id)
}

func (r *repository) DeleteHeroCondition(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_condition($1)", id)
}

func (r *repository) DeleteHeroInjury(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_injury($1)", id)
}

func (r *repository) DeleteHeroGoal(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_goal($1)", id)
}

func (r *repository) DeleteHeroConnection(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_connection($1)", id)
}

func (r *repository) DeleteHeroCompanion(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_companion($1)", id)
}

func (r *repository) DeleteHeroCulture(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT heroes.delete_hero_culture($1)", id)
}
