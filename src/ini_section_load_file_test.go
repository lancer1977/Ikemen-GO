package main

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestIniSectionLoadFile_ForwardsThroughLoadFileAndNoOpsOnMissingValues(t *testing.T) {
	is := IniSection{
		"asset": `  "select.def"  `,
	}
	var gotPath string
	load := func(path string) error {
		gotPath = path
		return nil
	}
	if err := is.LoadFile("asset", []string{t.TempDir()}, "chars", load); err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}
	if filepath.Base(gotPath) != "select.def" {
		t.Fatalf("LoadFile forwarded path = %q, want select.def", gotPath)
	}

	if err := (IniSection{}).LoadFile("missing", nil, "", func(string) error {
		return errors.New("should not be called")
	}); err != nil {
		t.Fatalf("missing LoadFile should be no-op, got error %v", err)
	}
}
