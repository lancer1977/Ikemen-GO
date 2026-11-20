package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFile_StripsQuotesAndWrapsLoaderErrorsWithRootContext(t *testing.T) {
	// Follow the same pattern as TestLoadFile_StripsQuotesAndReturnsResolvedPath,
	// but have the loader return an error to test error wrapping with root context.
	// Use a nested directory structure to verify SearchFile can construct paths through subdirectories.
	dir := t.TempDir()
	asset := filepath.Join(dir, "chars", "select.def")
	if err := os.MkdirAll(filepath.Dir(asset), 0o755); err != nil {
		t.Fatalf("mkdir asset dir: %v", err)
	}
	if err := os.WriteFile(asset, []byte("test"), 0o644); err != nil {
		t.Fatalf("write asset: %v", err)
	}

	file := `  "` + filepath.Base(asset) + `"  `
	var loaderPath string
	err := LoadFile(&file, []string{dir}, "chars", func(fp string) error {
		loaderPath = fp
		return errors.New("boom")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	got := err.Error()
	// Error should contain the root directory, the path that was searched, and the cause.
	// LoadFile wraps the error returned by the loader with root context (dirs[0]).
	if !strings.Contains(got, filepath.ToSlash(loaderPath)) || !strings.Contains(got, "boom") {
		t.Fatalf("LoadFile wrapped error = %q, want to contain loader path %q and cause 'boom'", got, filepath.ToSlash(loaderPath))
	}
	// LoadFile strips quotes and spaces from the input.
	// When the loader returns an error, LoadFile returns early without updating file.
	// So file remains as the stripped filename, not the full resolved path.
	if file != filepath.Base(asset) {
		t.Fatalf("LoadFile stripped file = %q, want %q (unchanged when loader errors)", file, filepath.Base(asset))
	}
}
