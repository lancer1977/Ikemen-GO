package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExist_ResolvesCaseInsensitiveGlobMatchForZipPath(t *testing.T) {
	dir := t.TempDir()
	actual := filepath.Join(dir, "Assets.ZIP")
	if err := os.WriteFile(actual, []byte("not a real zip but the path exists"), 0o644); err != nil {
		t.Fatalf("write zip path: %v", err)
	}

	got := FileExist(filepath.Join(dir, "assets.zip"))
	if got != filepath.ToSlash(actual) {
		t.Fatalf("FileExist(zip glob) = %q, want %q", got, filepath.ToSlash(actual))
	}
}
