package main

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestSetValueUpdate(t *testing.T) {
	t.Parallel()

	type sample struct {
		Enabled bool     `ini:"enabled"`
		Tags    []string `ini:"tags"`
	}

	s := &sample{}
	f := ini.Empty()
	sec, err := f.NewSection("sample")
	if err != nil {
		t.Fatalf("NewSection: %v", err)
	}

	if err := SetValueUpdate(s, f, "enabled", true); err != nil {
		t.Fatalf("SetValueUpdate(enabled): %v", err)
	}
	if err := SetValueUpdate(s, f, "tags", []string{"one", "two"}); err != nil {
		t.Fatalf("SetValueUpdate(tags): %v", err)
	}

	if !s.Enabled {
		t.Fatal("expected struct field to be updated")
	}
	if len(s.Tags) != 2 || s.Tags[0] != "one" || s.Tags[1] != "two" {
		t.Fatalf("unexpected tags: %#v", s.Tags)
	}
	if got := sec.Key("enabled").String(); got != "1" {
		t.Fatalf("INI enabled = %q, want 1", got)
	}
	if got := sec.Key("tags").String(); got != "one, two" {
		t.Fatalf("INI tags = %q, want %q", got, "one, two")
	}
}
