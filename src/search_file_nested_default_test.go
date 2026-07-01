package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSearchFile_UsesDefaultDirectoryUnderNestedRoot(t *testing.T) {
	dir := t.TempDir()
	root := filepath.Join(dir, "motif")
	fontDir := filepath.Join(root, "font")
	if err := os.MkdirAll(fontDir, 0o755); err != nil {
		t.Fatalf("mkdir nested default dir: %v", err)
	}
	wantPath := filepath.Join(fontDir, "select.fnt")
	if err := os.WriteFile(wantPath, []byte("font"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	got := SearchFile("select.fnt", []string{root}, "font")
	if got != filepath.ToSlash(wantPath) {
		t.Fatalf("SearchFile(nested default) = %q, want %q", got, filepath.ToSlash(wantPath))
	}
}
