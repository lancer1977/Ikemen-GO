package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFile_StripsQuotesAndReturnsResolvedPath(t *testing.T) {
	dir := t.TempDir()
	asset := filepath.Join(dir, "chars", "ryu.def")
	if err := os.MkdirAll(filepath.Dir(asset), 0o755); err != nil {
		t.Fatalf("mkdir asset dir: %v", err)
	}
	if err := os.WriteFile(asset, []byte("def"), 0o644); err != nil {
		t.Fatalf("write asset: %v", err)
	}

	file := `  "` + filepath.Base(asset) + `"  `
	var loaded string
	err := LoadFile(&file, []string{dir}, "chars", func(fp string) error {
		loaded = fp
		return nil
	})
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}
	if file != filepath.ToSlash(asset) {
		t.Fatalf("file = %q, want %q", file, filepath.ToSlash(asset))
	}
	if loaded != filepath.ToSlash(asset) {
		t.Fatalf("loaded = %q, want %q", loaded, filepath.ToSlash(asset))
	}
}

func TestLoadFile_FormatsErrorsWithSearchContext(t *testing.T) {
	file := "missing.def"
	err := LoadFile(&file, []string{"chars"}, "font", func(string) error {
		return os.ErrNotExist
	})
	if err == nil {
		t.Fatal("expected error")
	}
	got := err.Error()
	if !strings.Contains(got, "chars:\n") || !strings.Contains(got, "missing.def") || !strings.Contains(got, "file does not exist") {
		t.Fatalf("LoadFile() error = %q, want search context and original error", got)
	}
}
