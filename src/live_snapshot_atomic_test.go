package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAtomicFile_CreatesParentDirsAndReplacesContent(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "nested", "snapshot.json")

	if err := writeAtomicFile(path, []byte(`{"schema":"first"}`)); err != nil {
		t.Fatalf("first writeAtomicFile failed: %v", err)
	}
	if err := writeAtomicFile(path, []byte(`{"schema":"second"}`)); err != nil {
		t.Fatalf("second writeAtomicFile failed: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read atomic file: %v", err)
	}
	if string(raw) != `{"schema":"second"}` {
		t.Fatalf("unexpected final file content: %s", string(raw))
	}
}
