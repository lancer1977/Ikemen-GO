package main

import (
	"reflect"
	"testing"
)

func TestGetValueFromPatternMap_CaseInsensitiveMatchingDefect(t *testing.T) {
	t.Parallel()

	type item struct {
		Name string
	}
	type sample struct {
		Items map[string]*item `ini:"map:^item[0-9]+$"`
		Skip  map[string]*item `ini:"map:["`
	}

	var s sample

	// DEFECT (lancer1977/Ikemen-GO#21): getValueFromPatternMap fails to match mixed-case keys against lowercase patterns.
	// The regex pattern "^item[0-9]+$" (lowercase) is matched against "Item12" (mixed case),
	// causing the pattern match to fail. While the code normalizes the map key to lowercase
	// later, it checks the pattern BEFORE normalization. This prevents case-insensitive
	// pattern matching for map fields with regex-based key patterns.
	ok, got := getValueFromPatternMap(reflect.ValueOf(&s).Elem(), "Item12")
	if ok {
		t.Fatalf("expected NO pattern match for mixed-case 'Item12' against lowercase pattern '^item[0-9]+$'")
	}
	if got.IsValid() {
		t.Fatal("expected invalid value when pattern doesn't match")
	}

	// Lowercase should match
	ok, got = getValueFromPatternMap(reflect.ValueOf(&s).Elem(), "item12")
	if !ok {
		t.Fatal("expected pattern match for lowercase 'item12'")
	}
	if !got.IsValid() {
		t.Fatal("expected valid value")
	}
	if s.Items == nil {
		t.Fatal("expected map to be initialized")
	}
	if _, exists := s.Items["item12"]; !exists {
		t.Fatal("expected key 'item12' to be inserted")
	}

	ok, got = getValueFromPatternMap(reflect.ValueOf(&s).Elem(), "Other")
	if ok || got.IsValid() {
		t.Fatal("expected no match for unrelated field")
	}

	ok, got = getValueFromPatternMap(reflect.ValueOf(123), "Item12")
	if ok || got.IsValid() {
		t.Fatal("expected non-struct input to be rejected")
	}
}
