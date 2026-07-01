package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExist_ResolvesCaseInsensitiveGlobMatch(t *testing.T) {
	dir := t.TempDir()
	actual := filepath.Join(dir, "Stage.DEF")
	if err := os.WriteFile(actual, []byte("stage"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	got := FileExist(filepath.Join(dir, "stage.def"))
	if got != filepath.ToSlash(actual) {
		t.Fatalf("FileExist() = %q, want %q", got, filepath.ToSlash(actual))
	}
}
