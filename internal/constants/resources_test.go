package constants

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClassifierTypeSuffix(t *testing.T) {
	cases := map[string]string{
		"talents":           "talent",
		"path-types":        "path_type",
		"ancestry-subtypes": "ancestry_subtype",
		"equipment":         "equipment",
		"starting-kits":     "starting_kit",
	}
	for url, want := range cases {
		got, ok := ClassifierTypeSuffix(url)
		if !ok || got != want {
			t.Errorf("%s -> %q (ok=%v), want %q (ok=true)", url, got, ok, want)
		}
	}

	for _, bad := range []string{"", "widgets", "talent", "TALENTS", "../../etc/passwd", "talents; DROP"} {
		if _, ok := ClassifierTypeSuffix(bad); ok {
			t.Errorf("ClassifierTypeSuffix(%q) ok=true, want false", bad)
		}
	}
}

func TestClassifierTableName(t *testing.T) {
	cases := map[string]string{
		"talents":           "cl_talents",
		"path-types":        "cl_path_types",
		"ancestry-subtypes": "cl_ancestry_subtypes",
		"equipment":         "cl_equipments",
		"starting-kits":     "cl_starting_kits",
	}
	for url, want := range cases {
		got, ok := ClassifierTableName(url)
		if !ok || got != want {
			t.Errorf("%s -> %q (ok=%v), want %q (ok=true)", url, got, ok, want)
		}
	}

	if _, ok := ClassifierTableName("widgets"); ok {
		t.Error("ClassifierTableName(widgets) ok=true, want false")
	}
}

// TestClassifierAllowListMatchesMigrations asserts that every URL-plural entry
// in the classifiers map has a matching upsert_<suffix> and delete_<suffix>
// SQL file under database/migrations/R/classifiers/. Catches drift when
// someone adds a new classifier type on one side without the other.
//
// File-existence check (not a DB query) so the test can run without a live
// database. If the migration files move, update the relative path below.
func TestClassifierAllowListMatchesMigrations(t *testing.T) {
	migrationsDir := filepath.Join("..", "..", "..", "database", "migrations", "R", "classifiers")
	if _, err := os.Stat(migrationsDir); err != nil {
		t.Skipf("migrations dir not reachable from test working dir: %v", err)
	}

	for urlType, meta := range classifiers {
		for _, op := range []string{"upsert", "delete"} {
			path := filepath.Join(migrationsDir, "R__"+op+"_"+meta.suffix+".sql")
			if _, err := os.Stat(path); err != nil {
				t.Errorf("url-type %q expects %s_%s function but %s is missing", urlType, op, meta.suffix, path)
			}
		}
	}
}
