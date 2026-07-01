package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExist_ReturnsZipPathForArchiveRoot(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "chars.zip")
	if err := os.WriteFile(zipPath, []byte("not a real zip but the path exists"), 0o644); err != nil {
		t.Fatalf("write zip path: %v", err)
	}

	got := FileExist(zipPath)
	if got != filepath.ToSlash(zipPath) {
		t.Fatalf("FileExist(zip root) = %q, want %q", got, filepath.ToSlash(zipPath))
	}
}
