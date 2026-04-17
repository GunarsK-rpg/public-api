package constants

// RPG resource constants for permission checks.
const (
	ResourceClassifiers = "classifiers"
	ResourceHeroes      = "heroes"
	ResourceCampaigns   = "campaigns"
)

// classifierMeta holds the DB function suffix (singular) and the
// soft-deletable table name for a URL plural form.
type classifierMeta struct {
	suffix string
	table  string
}

// classifiers maps URL plural form to its DB function suffix and table name.
// Used for routing homebrew CRUD requests to classifiers.upsert_<suffix> /
// delete_<suffix> and to classifiers.restore_classifier(table, id).
var classifiers = map[string]classifierMeta{
	"ancestries":        {"ancestry", "cl_ancestries"},
	"ancestry-subtypes": {"ancestry_subtype", "cl_ancestry_subtypes"},
	"cultures":          {"culture", "cl_cultures"},
	"path-types":        {"path_type", "cl_path_types"},
	"paths":             {"path", "cl_paths"},
	"specialties":       {"specialty", "cl_specialties"},
	"talents":           {"talent", "cl_talents"},
	"actions":           {"action", "cl_actions"},
	"action-types":      {"action_type", "cl_action_types"},
	"equipments":        {"equipment", "cl_equipments"},
	"equipment-types":   {"equipment_type", "cl_equipment_types"},
	"modifications":     {"modification", "cl_modifications"},
	"skills":            {"skill", "cl_skills"},
	"expertises":        {"expertise", "cl_expertises"},
	"starting-kits":     {"starting_kit", "cl_starting_kits"},
}

// ClassifierTypeSuffix returns the DB function suffix for a URL plural form.
// Returns ok=false if the type is not in the allow-list.
func ClassifierTypeSuffix(urlType string) (suffix string, ok bool) {
	m, ok := classifiers[urlType]
	if !ok {
		return "", false
	}
	return m.suffix, true
}

// ClassifierTableName returns the soft-deletable table name for a URL plural
// form. Returns ok=false if the type is not in the allow-list.
func ClassifierTableName(urlType string) (table string, ok bool) {
	m, ok := classifiers[urlType]
	if !ok {
		return "", false
	}
	return m.table, true
}
