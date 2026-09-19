package main

import (
	"reflect"
	"testing"
)

func TestAssignToPatternMap_CaseInsensitiveMatchingDefect(t *testing.T) {
	t.Parallel()

	type entry struct {
		Name string `default:"seed"`
	}
	type sample struct {
		Items map[string]*entry `ini:"map:^item[0-9]+$"`
	}

	s := sample{}
	v := reflect.ValueOf(&s).Elem()

	// lancer1977/Ikemen-GO#21: assignToPatternMap should match case-insensitively.
	// The regex pattern "^item[0-9]+$" (lowercase) should match against "Item1" (mixed case)
	// by normalizing the key to lowercase before pattern matching.
	matched, mapElem, err := assignToPatternMap(v, "Item1", "alpha", true, "")
	if err != nil {
		t.Fatalf("assignToPatternMap(final insert mixed-case): %v", err)
	}
	if !matched {
		t.Fatalf("assignToPatternMap should match mixed-case key 'Item1' against lowercase pattern '^item[0-9]+$'")
	}
	if s.Items == nil || len(s.Items) != 1 {
		t.Fatalf("map not initialized: %#v", s.Items)
	}
	if _, ok := s.Items["item1"]; !ok {
		t.Fatalf("map should have key 'item1' (normalized from 'Item1'): %#v", s.Items)
	}
	if s.Items["item1"].Name != "alpha" {
		t.Fatalf("map value = %#v, want alpha", s.Items["item1"])
	}

	// Lowercase should also match (same key, already normalized)
	matched, mapElem, err = assignToPatternMap(v, "item1", "beta", true, "")
	if err != nil {
		t.Fatalf("assignToPatternMap(final update lowercase): %v", err)
	}
	if !matched {
		t.Fatal("assignToPatternMap should match lowercase key")
	}
	if mapElem.Kind() != reflect.Ptr || mapElem.IsNil() {
		t.Fatalf("expected pointer map element, got %#v", mapElem)
	}
	if s.Items["item1"].Name != "beta" {
		t.Fatalf("map value not updated: %#v", s.Items["item1"])
	}

	// Uppercase should also match and update the same key (case-insensitive)
	matched, mapElem, err = assignToPatternMap(v, "ITEM1", "gamma", true, "")
	if err != nil {
		t.Fatalf("assignToPatternMap(final update uppercase): %v", err)
	}
	if !matched {
		t.Fatalf("assignToPatternMap should match uppercase key 'ITEM1' against lowercase pattern")
	}
	if s.Items["item1"].Name != "gamma" {
		t.Fatalf("map value not updated from uppercase key: %#v", s.Items["item1"])
	}

	// Non-matching key should not match
	matched, _, err = assignToPatternMap(v, "other", "noop", true, "")
	if err != nil {
		t.Fatalf("assignToPatternMap(unmatched): %v", err)
	}
	if matched {
		t.Fatal("assignToPatternMap should not match non-matching key")
	}

	_ = mapElem
}
