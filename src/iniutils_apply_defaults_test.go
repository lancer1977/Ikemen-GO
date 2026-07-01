package main

import (
	"reflect"
	"testing"
)

func TestApplyDefaultsToValue(t *testing.T) {
	t.Parallel()

	type nested struct {
		Name string `default:"nested"`
	}
	type sample struct {
		Title   string   `default:"main"`
		Count   int      `default:"7"`
		Enabled bool     `default:"true"`
		Inner   *nested  `default:""`
		Items   []nested `default:""`
	}

	s := sample{
		Inner: &nested{},
		Items: []nested{{}},
	}

	applyDefaultsToValue(reflect.ValueOf(&s).Elem())

	if s.Title != "main" || s.Count != 7 || !s.Enabled {
		t.Fatalf("scalar defaults not applied: %#v", s)
	}
	if s.Inner == nil || s.Inner.Name != "nested" {
		t.Fatalf("pointer defaults not applied: %#v", s.Inner)
	}
	if len(s.Items) != 1 || s.Items[0].Name != "nested" {
		t.Fatalf("slice defaults not applied: %#v", s.Items)
	}
}
