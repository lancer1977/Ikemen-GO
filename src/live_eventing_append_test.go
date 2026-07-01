package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendJSONLine_CreatesDirectoryAndAppendsRecords(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "nested", "events.jsonl")

	if err := appendJSONLine(path, map[string]any{"schema": "first"}); err != nil {
		t.Fatalf("first appendJSONLine failed: %v", err)
	}
	if err := appendJSONLine(path, map[string]any{"schema": "second"}); err != nil {
		t.Fatalf("second appendJSONLine failed: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read appended file: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	if len(lines) != 2 {
		t.Fatalf("unexpected line count: %d", len(lines))
	}
	if !strings.Contains(lines[0], `"schema":"first"`) || !strings.Contains(lines[1], `"schema":"second"`) {
		t.Fatalf("unexpected appended content: %q", string(raw))
	}
}

func TestAppendJSONLine_NoOpForBlankPath(t *testing.T) {
	if err := appendJSONLine("", map[string]any{"schema": "noop"}); err != nil {
		t.Fatalf("expected blank path to be a no-op, got %v", err)
	}
}
