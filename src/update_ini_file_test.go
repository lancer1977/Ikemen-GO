package main

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestUpdateINIFileRejectsInvalidQueryAndSingleLevelQuery(t *testing.T) {
	type sample struct {
		Name string `ini:"name"`
	}

	f := ini.Empty()
	if err := updateINIFile(&sample{}, f, "", "ryu"); err == nil {
		t.Fatal("updateINIFile should reject an empty query")
	}

	// Not a defect: see the equivalent comment on
	// TestSetValueUpdate_RequiresTwoLevelSectionKeyQuery in
	// iniutils_set_value_update_test.go. updateINIFile always needs a query
	// with at least two segments to split into an INI [section] and key; a
	// single-segment query like "name" has nothing left over for the key.
	// No real Config/Motif/Storyboard field lives at struct root, so no
	// production caller ever passes a single-segment query.
	err := updateINIFile(&sample{}, f, "name", "ryu")
	if err == nil {
		t.Fatalf("updateINIFile(simple field): expected error but succeeded")
	}
	if err.Error() != "unable to determine key name from query 'name'" {
		t.Fatalf("updateINIFile(simple field) error = %v, want 'unable to determine key name from query 'name''", err)
	}
}
