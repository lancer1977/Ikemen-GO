package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteStatsPretty_PrettyPrintsValidJsonAndFallsBackOnInvalidInput(t *testing.T) {
	tempDir := t.TempDir()
	prettyPath := filepath.Join(tempDir, "pretty.json")
	if err := writeStatsPretty(prettyPath, []byte(`{"a":1,"b":{"c":2}}`)); err != nil {
		t.Fatalf("writeStatsPretty failed: %v", err)
	}
	raw, err := os.ReadFile(prettyPath)
	if err != nil {
		t.Fatalf("read pretty stats: %v", err)
	}
	if string(raw) != "{\n  \"a\": 1,\n  \"b\": {\n    \"c\": 2\n  }\n}" {
		t.Fatalf("unexpected pretty output: %q", string(raw))
	}

	compactPath := filepath.Join(tempDir, "compact.json")
	if err := writeStatsPretty(compactPath, []byte(`{"a":1,`)); err != nil {
		t.Fatalf("writeStatsPretty fallback failed: %v", err)
	}
	raw, err = os.ReadFile(compactPath)
	if err != nil {
		t.Fatalf("read compact stats: %v", err)
	}
	if string(raw) != `{"a":1,` {
		t.Fatalf("unexpected fallback output: %q", string(raw))
	}
}
