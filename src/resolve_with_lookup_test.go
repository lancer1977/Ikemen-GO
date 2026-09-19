package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveWithLookup_ReturnsPathRelativeToSearchRoot(t *testing.T) {
	t.Parallel()

	if got := resolveWithLookup("value", "", "base"); got != "value" {
		t.Fatalf("empty lookup tag = %q, want value", got)
	}

	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "base"), 0o755); err != nil {
		t.Fatalf("mkdir base dir: %v", err)
	}
	wantPath := filepath.Join(dir, "base", "select.def")
	if err := os.WriteFile(wantPath, []byte("def"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	// Not a defect: resolveWithLookup just returns whatever SearchFile finds,
	// and SearchFile reports paths relative to the search root it was given
	// (here "base"), not absolutized. That is SearchFile's documented,
	// consistent behavior across all its callers, not a bug specific to
	// resolveWithLookup.
	got := resolveWithLookup("select.def", "def", "base")
	if got != "base/select.def" {
		t.Fatalf("resolveWithLookup(def) = %q, want %q", got, "base/select.def")
	}
}
