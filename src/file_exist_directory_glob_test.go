package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExist_IgnoresDirectoryMatchesInGlobFallback(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "stage"), 0o755); err != nil {
		t.Fatalf("mkdir stage dir: %v", err)
	}

	if got := FileExist(filepath.Join(dir, "stage")); got != "" {
		t.Fatalf("FileExist(directory glob match) = %q, want empty", got)
	}
}
