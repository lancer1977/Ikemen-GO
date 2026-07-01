package main

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestUpdateINIFileRejectsInvalidQueryAndWritesSimpleField(t *testing.T) {
	type sample struct {
		Name string `ini:"name"`
	}

	f := ini.Empty()
	if err := updateINIFile(&sample{}, f, "", "ryu"); err == nil {
		t.Fatal("updateINIFile should reject an empty query")
	}

	if err := updateINIFile(&sample{}, f, "name", "ryu"); err != nil {
		t.Fatalf("updateINIFile(simple field) error = %v", err)
	}
	if got := f.Section("").Key("name").String(); got != "ryu" {
		t.Fatalf("updateINIFile wrote %q, want ryu", got)
	}
}
