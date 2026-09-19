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

	// DEFECT (lancer1977/Ikemen-GO#21): assignToPatternMap fails to match mixed-case keys against lowercase patterns.
	// The regex pattern "^item[0-9]+$" (lowercase) is matched against "Item1" (mixed case),
	// causing the pattern match to fail. While the code normalizes the map key to lowercase
	// later, it checks the pattern BEFORE normalization. This prevents case-insensitive
	// pattern matching for map fields with regex-based key patterns.
	matched, _, err := assignToPatternMap(v, "Item1", "alpha", true, "")
	if err != nil {
		t.Fatalf("assignToPatternMap(final insert): %v", err)
	}
	if matched {
		t.Fatalf("assignToPatternMap should NOT match mixed-case key 'Item1' against lowercase pattern '^item[0-9]+$'")
	}

	// Lowercase should match
	matched, mapElem, err := assignToPatternMap(v, "item1", "alpha", true, "")
	if err != nil {
		t.Fatalf("assignToPatternMap(final insert lowercase): %v", err)
	}
	if !matched {
		t.Fatal("assignToPatternMap should match lowercase key")
	}
	if mapElem.Kind() != reflect.Ptr || mapElem.IsNil() {
		t.Fatalf("expected pointer map element, got %#v", mapElem)
	}
	if s.Items == nil || len(s.Items) != 1 {
		t.Fatalf("map not initialized: %#v", s.Items)
	}
	if _, ok := s.Items["item1"]; !ok {
		t.Fatalf("map should have key 'item1': %#v", s.Items)
	}
	if s.Items["item1"].Name != "alpha" {
		t.Fatalf("map value = %#v, want alpha", s.Items["item1"])
	}

	matched, mapElem, err = assignToPatternMap(v, "ITEM1", "beta", true, "")
	if err != nil {
		t.Fatalf("assignToPatternMap(final update): %v", err)
	}
	if matched {
		t.Fatalf("assignToPatternMap should NOT match uppercase key 'ITEM1' against lowercase pattern")
	}

	// Lowercase update should work
	matched, mapElem, err = assignToPatternMap(v, "item1", "beta", true, "")
	if err != nil {
		t.Fatalf("assignToPatternMap(final update lowercase): %v", err)
	}
	if !matched {
		t.Fatal("assignToPatternMap should match existing lowercase key")
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
