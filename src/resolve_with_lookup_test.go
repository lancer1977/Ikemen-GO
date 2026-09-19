package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveWithLookup_ReturnsRelativePathNotAbsolute(t *testing.T) {
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

	// DEFECT: resolveWithLookup returns relative paths instead of absolute paths.
	// SearchFile finds the file and returns the path as it was used to search (relative),
	// but for configuration files, absolute paths are more useful and prevent issues
	// when the working directory changes. User impact: file paths in loaded configs are
	// relative to where the config was loaded from, not absolute, causing path resolution
	// failures if the working directory changes later.
	got := resolveWithLookup("select.def", "def", "base")
	if got != "base/select.def" {
		t.Fatalf("resolveWithLookup(def) = %q, want %q", got, "base/select.def")
	}
}
