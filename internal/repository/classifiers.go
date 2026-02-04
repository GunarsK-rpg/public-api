package repository

import (
	"context"
	"encoding/json"
)

// ClassifierRepository defines methods for classifier data access.
type ClassifierRepository interface {
	// Batch getter (all classifiers in one call)
	GetAllClassifiers(ctx context.Context, auth AuthContext) (json.RawMessage, error)

	// Simple getters (no parameters)
	GetAttributeTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetAttributes(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetDerivedStats(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetDerivedStatValues(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetSkills(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetExpertiseTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetPaths(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetSurges(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetRadiantOrders(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetAncestries(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetActivationTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetActionTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetDamageTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetUnits(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetEquipmentTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetEquipmentAttributes(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetConditions(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetInjuries(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetGoalStatus(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetConnectionTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetCompanionTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetCultures(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetTiers(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetLevels(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetStartingKits(ctx context.Context, auth AuthContext) (json.RawMessage, error)

	// Filtered getters (take JSONB filter)
	GetExpertises(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error)
	GetSpecialties(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error)
	GetSingerForms(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error)
	GetTalents(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error)
	GetActions(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error)
	GetActionLinks(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error)
	GetEquipments(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error)
}

// Batch getter (all classifiers in one call)

func (r *repository) GetAllClassifiers(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_all_classifiers()")
}

// Simple getters (no parameters)

func (r *repository) GetAttributeTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_attribute_types()")
}

func (r *repository) GetAttributes(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_attributes()")
}

func (r *repository) GetDerivedStats(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_derived_stats()")
}

func (r *repository) GetDerivedStatValues(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_derived_stat_values()")
}

func (r *repository) GetSkills(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_skills()")
}

func (r *repository) GetExpertiseTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_expertise_types()")
}

func (r *repository) GetPaths(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_paths()")
}

func (r *repository) GetSurges(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_surges()")
}

func (r *repository) GetRadiantOrders(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_radiant_orders()")
}

func (r *repository) GetAncestries(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_ancestries()")
}

func (r *repository) GetActivationTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_activation_types()")
}

func (r *repository) GetActionTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_action_types()")
}

func (r *repository) GetDamageTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_damage_types()")
}

func (r *repository) GetUnits(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_units()")
}

func (r *repository) GetEquipmentTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_equipment_types()")
}

func (r *repository) GetEquipmentAttributes(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_equipment_attributes()")
}

func (r *repository) GetConditions(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_conditions()")
}

func (r *repository) GetInjuries(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_injuries()")
}

func (r *repository) GetGoalStatus(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_goal_status()")
}

func (r *repository) GetConnectionTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_connection_types()")
}

func (r *repository) GetCompanionTypes(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_companion_types()")
}

func (r *repository) GetCultures(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_cultures()")
}

func (r *repository) GetTiers(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_tiers()")
}

func (r *repository) GetLevels(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_levels()")
}

func (r *repository) GetStartingKits(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_starting_kits()")
}

// Filtered getters (take JSONB filter)

func (r *repository) GetExpertises(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_expertises($1)", filter)
}

func (r *repository) GetSpecialties(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_specialties($1)", filter)
}

func (r *repository) GetSingerForms(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_singer_forms($1)", filter)
}

func (r *repository) GetTalents(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_talents($1)", filter)
}

func (r *repository) GetActions(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_actions($1)", filter)
}

func (r *repository) GetActionLinks(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_action_links($1)", filter)
}

func (r *repository) GetEquipments(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_equipments($1)", filter)
}
