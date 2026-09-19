package main

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestSetValueUpdate_SimpleFieldUpdatesFails(t *testing.T) {
	t.Parallel()

	type sample struct {
		Enabled bool     `ini:"enabled"`
		Tags    []string `ini:"tags"`
	}

	s := &sample{}
	f := ini.Empty()

	// DEFECT: updateINIFile fails to write simple (root-level) struct fields to INI.
	// When processing a query like "enabled", the code incorrectly treats the field
	// tag as a section name instead of a key name, leaving keyNameParts empty and
	// causing "unable to determine key name" error. This breaks SetValueUpdate for
	// simple struct fields. User impact: cannot update simple config values via INI.
	err := SetValueUpdate(s, f, "enabled", true)
	if err == nil {
		t.Fatalf("SetValueUpdate(enabled): expected error but succeeded")
	}
	if err.Error() != "unable to determine key name from query 'enabled'" {
		t.Fatalf("SetValueUpdate(enabled): wrong error: %v", err)
	}

	err = SetValueUpdate(s, f, "tags", []string{"one", "two"})
	if err == nil {
		t.Fatalf("SetValueUpdate(tags): expected error but succeeded")
	}
	if err.Error() != "unable to determine key name from query 'tags'" {
		t.Fatalf("SetValueUpdate(tags): wrong error: %v", err)
	}

	// SetValue itself works for struct fields, it's just updateINIFile that fails
	if err := SetValue(s, "enabled", true); err != nil {
		t.Fatalf("SetValue(enabled): %v", err)
	}
	if !s.Enabled {
		t.Fatal("expected struct field to be updated")
	}
	if err := SetValue(s, "tags", []string{"one", "two"}); err != nil {
		t.Fatalf("SetValue(tags): %v", err)
	}
	if len(s.Tags) != 2 || s.Tags[0] != "one" || s.Tags[1] != "two" {
		t.Fatalf("unexpected tags: %#v", s.Tags)
	}
}
