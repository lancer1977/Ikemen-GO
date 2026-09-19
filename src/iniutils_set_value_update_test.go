package main

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestSetValueUpdate_RequiresTwoLevelSectionKeyQuery(t *testing.T) {
	t.Parallel()

	type sample struct {
		Enabled bool     `ini:"enabled"`
		Tags    []string `ini:"tags"`
	}

	s := &sample{}
	f := ini.Empty()

	// Not a defect: updateINIFile derives an INI [section]/key pair from the
	// query, so it needs at least two path segments (the first becomes the
	// section, the rest become the key). A single-segment query like "enabled"
	// has nothing to use as the key, hence the error. This matches production
	// usage: every real Config/Motif/Storyboard field lives at least one level
	// deep (e.g. "Options.Difficulty", "Sound.MaxBGMVolume"), and every
	// production call site already passes a two-part-or-deeper query.
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
