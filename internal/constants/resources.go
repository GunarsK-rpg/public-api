package constants

// RPG resource constants for permission checks.
const (
	ResourceClassifiers = "classifiers"
	ResourceHeroes      = "heroes"
	ResourceCampaigns   = "campaigns"
)

// classifierTypes maps URL plural form to DB function suffix (singular).
// Used for routing homebrew CRUD requests to the matching upsert_<suffix> /
// delete_<suffix> SQL functions and for restore_classifier table-name lookup.
var classifierTypes = map[string]string{
	"ancestries":        "ancestry",
	"ancestry-subtypes": "ancestry_subtype",
	"cultures":          "culture",
	"path-types":        "path_type",
	"paths":             "path",
	"specialties":       "specialty",
	"talents":           "talent",
	"actions":           "action",
	"action-types":      "action_type",
	"equipment":         "equipment",
	"equipment-types":   "equipment_type",
	"modifications":     "modification",
	"skills":            "skill",
	"expertises":        "expertise",
	"starting-kits":     "starting_kit",
}

// classifierTables maps URL plural form to soft-deletable table name used
// by classifiers.restore_classifier(table_name, id).
var classifierTables = map[string]string{
	"ancestries":        "cl_ancestries",
	"ancestry-subtypes": "cl_ancestry_subtypes",
	"cultures":          "cl_cultures",
	"path-types":        "cl_path_types",
	"paths":             "cl_paths",
	"specialties":       "cl_specialties",
	"talents":           "cl_talents",
	"actions":           "cl_actions",
	"action-types":      "cl_action_types",
	"equipment":         "cl_equipments",
	"equipment-types":   "cl_equipment_types",
	"modifications":     "cl_modifications",
	"skills":            "cl_skills",
	"expertises":        "cl_expertises",
	"starting-kits":     "cl_starting_kits",
}

// ClassifierTypeSuffix returns the DB function suffix for an URL plural form.
// Returns ok=false if the type is not in the allow-list.
func ClassifierTypeSuffix(urlType string) (suffix string, ok bool) {
	suffix, ok = classifierTypes[urlType]
	return
}

// ClassifierTableName returns the soft-deletable table name for an URL plural
// form. Returns ok=false if the type is not in the allow-list.
func ClassifierTableName(urlType string) (table string, ok bool) {
	table, ok = classifierTables[urlType]
	return
}
