package main

import (
	"reflect"
	"testing"
)

// Embedded must be exported: findFieldByINITag skips fields whose PkgPath is
// set, and an embedded unexported type produces an unexported field.
type iniLookupEmbedded struct {
	Child string `ini:"child"`
}

type IniLookupEmbedded struct {
	Child string `ini:"child"`
}

func TestIniutilsFieldLookupHelpers(t *testing.T) {
	// Default is declared before the embedded struct because findDefaultField
	// returns the first untagged field, and an anonymous field is untagged too.
	type sample struct {
		Direct  string `ini:"direct"`
		Default string
		IniLookupEmbedded
		Mapping map[string]int `ini:"mapping"`
	}

	v := reflect.ValueOf(&sample{
		Direct:            "a",
		Default:           "c",
		IniLookupEmbedded: IniLookupEmbedded{Child: "b"},
		Mapping:           map[string]int{},
	}).Elem()

	if got, field, ok := findFieldByINITag(v, "direct"); !ok || got.String() != "a" || field.Name != "Direct" {
		t.Fatalf("findFieldByINITag(direct) = %q %v %v", got, field.Name, ok)
	}
	if got, field, ok := findFieldByINITag(v, "child"); !ok || got.String() != "b" || field.Name != "Child" {
		t.Fatalf("findFieldByINITag(child) = %q %v %v", got, field.Name, ok)
	}
	// No ini tag on Default, so it resolves through the field-name fallback.
	if got, field, ok := findFieldByINITag(v, "default"); !ok || got.String() != "c" || field.Name != "Default" {
		t.Fatalf("findFieldByINITag(default) = %q %v %v", got, field.Name, ok)
	}
	if got, _, ok := findFieldByINITag(v, "missing"); ok || got.IsValid() {
		t.Fatalf("findFieldByINITag(missing) = %v %v", got, ok)
	}

	if got, field, ok := findMapFieldWithTag(v, "mapping"); !ok || got.Kind() != reflect.Map || field.Name != "Mapping" {
		t.Fatalf("findMapFieldWithTag(mapping) = %v %v %v", got, field.Name, ok)
	}
	// "direct" is tagged but is not a map, so the map lookup must reject it.
	if got, _, ok := findMapFieldWithTag(v, "direct"); ok || got.IsValid() {
		t.Fatalf("findMapFieldWithTag(direct) = %v %v", got, ok)
	}

	if got, field, ok := findDefaultField(v); !ok || got.String() != "c" || field.Name != "Default" {
		t.Fatalf("findDefaultField = %q %v %v", got, field.Name, ok)
	}

	// An unexported embedded type is skipped entirely: its fields are not
	// settable through reflect, so the helper must not hand them back.
	type shadowed struct {
		iniLookupEmbedded
	}
	sv := reflect.ValueOf(&shadowed{}).Elem()
	if got, _, ok := findFieldByINITag(sv, "child"); ok || got.IsValid() {
		t.Fatalf("findFieldByINITag through unexported embed = %v %v", got, ok)
	}
}
