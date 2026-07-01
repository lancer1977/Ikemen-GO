package main

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestResolveSectionForWrite(t *testing.T) {
	t.Parallel()

	f := ini.Empty()
	if got := resolveSectionForWrite(f, "select_screen"); got != "select_screen" {
		t.Fatalf("missing section = %q", got)
	}

	sec, err := f.NewSection("select screen")
	if err != nil {
		t.Fatalf("NewSection: %v", err)
	}
	sec.NameMapper = nil

	if got := resolveSectionForWrite(f, "select_screen"); got != "select screen" {
		t.Fatalf("underscore fallback = %q", got)
	}
	if got := resolveSectionForWrite(f, "select screen"); got != "select screen" {
		t.Fatalf("exact section = %q", got)
	}
	if got := resolveSectionForWrite(nil, "x"); got != "x" {
		t.Fatalf("nil file = %q", got)
	}
}
