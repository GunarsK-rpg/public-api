package constants

import (
	"testing"
)

func TestClassifierTypeSuffix(t *testing.T) {
	cases := map[string]string{
		"talents":           "talent",
		"path-types":        "path_type",
		"ancestry-subtypes": "ancestry_subtype",
		"equipments":        "equipment",
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
		"equipments":        "cl_equipments",
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
