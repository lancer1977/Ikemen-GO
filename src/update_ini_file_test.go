package main

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestUpdateINIFileRejectsInvalidQueryAndFailsOnSimpleField(t *testing.T) {
	type sample struct {
		Name string `ini:"name"`
	}

	f := ini.Empty()
	if err := updateINIFile(&sample{}, f, "", "ryu"); err == nil {
		t.Fatal("updateINIFile should reject an empty query")
	}

	// DEFECT: updateINIFile fails to write simple (root-level) struct fields to INI.
	// When processing a query like "name", the code incorrectly treats the field
	// tag as a section name instead of a key name, leaving keyNameParts empty and
	// causing "unable to determine key name" error. This prevents users from
	// persisting simple config values back to INI files after modifications.
	err := updateINIFile(&sample{}, f, "name", "ryu")
	if err == nil {
		t.Fatalf("updateINIFile(simple field): expected error but succeeded")
	}
	if err.Error() != "unable to determine key name from query 'name'" {
		t.Fatalf("updateINIFile(simple field) error = %v, want 'unable to determine key name from query 'name''", err)
	}
}
