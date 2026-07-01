package main

import (
	"reflect"
	"testing"
)

func TestIniutilsFieldLookupHelpers(t *testing.T) {
	type embedded struct {
		Child string `ini:"child"`
	}
	type sample struct {
		Direct string `ini:"direct"`
		embedded
		Default string
		Mapping map[string]int `ini:"mapping"`
	}

	v := reflect.ValueOf(&sample{Direct: "a", embedded: embedded{Child: "b"}, Default: "c", Mapping: map[string]int{}}).Elem()

	if got, field, ok := findFieldByINITag(v, "direct"); !ok || got.String() != "a" || field.Name != "Direct" {
		t.Fatalf("findFieldByINITag(direct) = %#v %v %v", got, field.Name, ok)
	}
	if got, field, ok := findFieldByINITag(v, "child"); !ok || got.String() != "b" || field.Name != "Child" {
		t.Fatalf("findFieldByINITag(child) = %#v %v %v", got, field.Name, ok)
	}
	if got, field, ok := findFieldByINITag(v, "default"); !ok || got.String() != "c" || field.Name != "Default" {
		t.Fatalf("findFieldByINITag(default) = %#v %v %v", got, field.Name, ok)
	}
	if got, _, ok := findFieldByINITag(v, "missing"); ok || got.IsValid() {
		t.Fatalf("findFieldByINITag(missing) = %#v %v", got, ok)
	}

	if got, field, ok := findMapFieldWithTag(v, "mapping"); !ok || got.Kind() != reflect.Map || field.Name != "Mapping" {
		t.Fatalf("findMapFieldWithTag(mapping) = %#v %v %v", got, field.Name, ok)
	}
	if got, _, ok := findMapFieldWithTag(v, "direct"); ok || got.IsValid() {
		t.Fatalf("findMapFieldWithTag(direct) = %#v %v", got, ok)
	}

	if got, field, ok := findDefaultField(v); !ok || got.String() != "c" || field.Name != "Default" {
		t.Fatalf("findDefaultField = %#v %v %v", got, field.Name, ok)
	}
}
