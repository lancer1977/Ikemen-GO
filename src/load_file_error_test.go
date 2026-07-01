package main

import (
	"errors"
	"strings"
	"testing"
)

func TestLoadFile_StripsQuotesAndWrapsLoaderErrorsWithRootContext(t *testing.T) {
	file := `  "select.def"  `
	dirs := []string{"chars"}
	err := LoadFile(&file, dirs, "font", func(path string) error {
		if path != "chars/font/select.def" {
			t.Fatalf("loader path = %q, want %q", path, "chars/font/select.def")
		}
		return errors.New("boom")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	got := err.Error()
	if !strings.Contains(got, "chars:") || !strings.Contains(got, "chars/font/select.def") || !strings.Contains(got, "boom") {
		t.Fatalf("LoadFile wrapped error = %q, want root, path, and cause", got)
	}
	if file != "chars/font/select.def" {
		t.Fatalf("LoadFile updated file = %q, want %q", file, "chars/font/select.def")
	}
}
