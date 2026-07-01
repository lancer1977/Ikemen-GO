package main

import (
	"errors"
	"strings"
	"testing"
)

func TestLoadFile_WrapsLoaderErrorsWithoutRootContext(t *testing.T) {
	file := "select.def"
	err := LoadFile(&file, nil, "font", func(path string) error {
		if path != "select.def" {
			t.Fatalf("loader path = %q, want %q", path, "select.def")
		}
		return errors.New("boom")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	got := err.Error()
	if strings.Contains(got, ":\n") || !strings.Contains(got, "select.def") || !strings.Contains(got, "boom") {
		t.Fatalf("LoadFile no-root wrapped error = %q, want path and cause without root prefix", got)
	}
	if file != "select.def" {
		t.Fatalf("LoadFile updated file = %q, want %q", file, "select.def")
	}
}
