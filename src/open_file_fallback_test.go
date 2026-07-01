package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenFile_FallsBackToPlainPathWhenZipArchiveIsMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.zip", "chars", "ryu.def")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir fallback path: %v", err)
	}
	if err := os.WriteFile(path, []byte("fallback text"), 0o644); err != nil {
		t.Fatalf("write fallback file: %v", err)
	}

	rc, err := OpenFile(path)
	if err != nil {
		t.Fatalf("OpenFile(fallback) error = %v", err)
	}
	defer rc.Close()

	buf, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read fallback file: %v", err)
	}
	if string(buf) != "fallback text" {
		t.Fatalf("OpenFile(fallback) = %q, want %q", string(buf), "fallback text")
	}
}
