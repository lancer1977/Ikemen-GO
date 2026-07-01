package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSearchFile_TrimsQuotesCommentsAndUsesDefaultDirectoryAgain(t *testing.T) {
	dir := t.TempDir()
	defaultDir := filepath.Join(dir, "font")
	if err := os.MkdirAll(defaultDir, 0o755); err != nil {
		t.Fatalf("mkdir default dir: %v", err)
	}
	wantPath := filepath.Join(defaultDir, "select.fnt")
	if err := os.WriteFile(wantPath, []byte("font"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	got := SearchFile(`  "select.fnt" ; comment`, []string{dir}, "font")
	if got != filepath.ToSlash(wantPath) {
		t.Fatalf("SearchFile() = %q, want %q", got, filepath.ToSlash(wantPath))
	}
}

func TestSearchFile_ReturnsAbsoluteAndFallbackOriginalWhenMissing(t *testing.T) {
	dir := t.TempDir()
	absPath := filepath.Join(dir, "data.def")
	if err := os.WriteFile(absPath, []byte("def"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if got := SearchFile(absPath, nil); got != filepath.ToSlash(absPath) {
		t.Fatalf("SearchFile(abs) = %q, want %q", got, filepath.ToSlash(absPath))
	}

	missing := "missing.def"
	if got := SearchFile(missing, nil); got != missing {
		t.Fatalf("SearchFile(missing) = %q, want original %q", got, missing)
	}
}
