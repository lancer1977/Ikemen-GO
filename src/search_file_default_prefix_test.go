package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSearchFile_SkipsDuplicatingDefaultDirectoryPrefix(t *testing.T) {
	dir := t.TempDir()
	defaultDir := filepath.Join(dir, "font")
	if err := os.MkdirAll(filepath.Join(dir, "font"), 0o755); err != nil {
		t.Fatalf("mkdir font dir: %v", err)
	}
	wantPath := filepath.Join(defaultDir, "select.fnt")
	if err := os.WriteFile(wantPath, []byte("font"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	got := SearchFile(filepath.ToSlash(filepath.Join("font", "select.fnt")), []string{dir}, "font")
	if got != filepath.ToSlash(wantPath) {
		t.Fatalf("SearchFile(default-prefix) = %q, want %q", got, filepath.ToSlash(wantPath))
	}
}
