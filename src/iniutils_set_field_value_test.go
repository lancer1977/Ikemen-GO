package main

import (
	"reflect"
	"testing"
)

func TestSetFieldValue(t *testing.T) {
	t.Parallel()

	type nested struct {
		Value string `default:"inner"`
	}
	type sample struct {
		Name  string `ini:"name"`
		Count int    `ini:"count"`
		Flags []int  `ini:"flags" default:"1,2"`
		Fixed [3]int `ini:"fixed" default:"4,5,6"`
		Inner nested `ini:"inner"`
	}

	s := sample{}
	if err := setFieldValue(reflect.ValueOf(&s).Elem().FieldByName("Name"), "Ryu", "", "name", "", ""); err != nil {
		t.Fatalf("setFieldValue(name): %v", err)
	}
	if err := setFieldValue(reflect.ValueOf(&s).Elem().FieldByName("Count"), "12", "", "count", "", ""); err != nil {
		t.Fatalf("setFieldValue(count): %v", err)
	}
	if err := setFieldValue(reflect.ValueOf(&s).Elem().FieldByName("Flags"), "3, 4", "1,2", "flags", "", ""); err != nil {
		t.Fatalf("setFieldValue(flags): %v", err)
	}
	if err := setFieldValue(reflect.ValueOf(&s).Elem().FieldByName("Fixed"), "7,,9", "4,5,6", "fixed", "", ""); err != nil {
		t.Fatalf("setFieldValue(fixed): %v", err)
	}
	if err := setFieldValue(reflect.ValueOf(&s).Elem().FieldByName("Inner"), "updated", "", "inner", "", ""); err != nil {
		t.Fatalf("setFieldValue(inner): %v", err)
	}

	if s.Name != "Ryu" || s.Count != 12 {
		t.Fatalf("scalar fields not set: %#v", s)
	}
	if len(s.Flags) != 2 || s.Flags[0] != 3 || s.Flags[1] != 4 {
		t.Fatalf("slice field = %#v, want [3 4]", s.Flags)
	}
	if s.Fixed != [3]int{7, 5, 9} {
		t.Fatalf("array field = %#v, want [7 5 9]", s.Fixed)
	}
	if s.Inner.Value != "updated" {
		t.Fatalf("struct default field = %#v, want updated", s.Inner)
	}

	if err := setFieldValue(reflect.ValueOf(&s).Elem().FieldByName("Name"), "", "", "name", "", ""); err != nil {
		t.Fatalf("setFieldValue blank string: %v", err)
	}
	if s.Name != "" {
		t.Fatalf("blank string should clear field, got %q", s.Name)
	}

	if err := setFieldValue(reflect.ValueOf(&s).Elem().FieldByName("Flags"), "", "", "flags", "", ""); err != nil {
		t.Fatalf("setFieldValue blank slice: %v", err)
	}
	if len(s.Flags) != 0 {
		t.Fatalf("blank slice should clear field, got %#v", s.Flags)
	}

	if err := setFieldValue(reflect.ValueOf(&s).Elem().FieldByName("Count"), "", "99", "count", "", ""); err != nil {
		t.Fatalf("setFieldValue default fallback: %v", err)
	}
	if s.Count != 99 {
		t.Fatalf("default fallback = %d, want 99", s.Count)
	}
}
