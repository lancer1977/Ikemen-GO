package main

import (
	"reflect"
	"testing"
)

func TestGetValueFromPatternMap(t *testing.T) {
	t.Parallel()

	type item struct {
		Name string
	}
	type sample struct {
		Items map[string]*item `ini:"map:^item[0-9]+$"`
		Skip  map[string]*item `ini:"map:["`
	}

	var s sample
	ok, got := getValueFromPatternMap(reflect.ValueOf(&s).Elem(), "Item12")
	if !ok {
		t.Fatal("expected pattern match")
	}
	if !got.IsValid() {
		t.Fatal("expected valid value")
	}
	if s.Items == nil {
		t.Fatal("expected map to be initialized")
	}
	if _, exists := s.Items["item12"]; !exists {
		t.Fatal("expected lower-cased key to be inserted")
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
