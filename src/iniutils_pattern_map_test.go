package main

import (
	"reflect"
	"testing"
)

func TestAssignToPatternMap(t *testing.T) {
	t.Parallel()

	type entry struct {
		Name string `default:"seed"`
	}
	type sample struct {
		Items map[string]*entry `ini:"map:^item[0-9]+$"`
	}

	s := sample{}
	v := reflect.ValueOf(&s).Elem()

	matched, mapElem, err := assignToPatternMap(v, "Item1", "alpha", true, "")
	if err != nil {
		t.Fatalf("assignToPatternMap(final insert): %v", err)
	}
	if !matched {
		t.Fatal("assignToPatternMap should match item key")
	}
	if mapElem.Kind() != reflect.Ptr || mapElem.IsNil() {
		t.Fatalf("expected pointer map element, got %#v", mapElem)
	}
	if s.Items == nil || len(s.Items) != 1 {
		t.Fatalf("map not initialized: %#v", s.Items)
	}
	if _, ok := s.Items["item1"]; !ok {
		t.Fatalf("map key should be normalized to lowercase: %#v", s.Items)
	}
	if s.Items["item1"].Name != "alpha" {
		t.Fatalf("map value = %#v, want alpha", s.Items["item1"])
	}

	matched, mapElem, err = assignToPatternMap(v, "ITEM1", "beta", true, "")
	if err != nil {
		t.Fatalf("assignToPatternMap(final update): %v", err)
	}
	if !matched {
		t.Fatal("assignToPatternMap should match existing key")
	}
	if s.Items["item1"].Name != "beta" {
		t.Fatalf("map value not updated: %#v", s.Items["item1"])
	}

	matched, _, err = assignToPatternMap(v, "other", "noop", true, "")
	if err != nil {
		t.Fatalf("assignToPatternMap(unmatched): %v", err)
	}
	if matched {
		t.Fatal("assignToPatternMap should not match non-matching key")
	}

	_ = mapElem
}
