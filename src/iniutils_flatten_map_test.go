package main

import (
	"reflect"
	"testing"
)

func TestAssignToFlattenMapField(t *testing.T) {
	t.Parallel()

	type entry struct {
		Value string `ini:"value"`
	}
	type sample struct {
		Items map[string]*entry `flatten:"true"`
	}

	s := sample{}
	fieldVal := reflect.ValueOf(&s).Elem().FieldByName("Items")
	sf, _ := reflect.TypeOf(s).FieldByName("Items")
	parts := []queryPart{
		{name: "bg"},
		{name: "menu"},
		{name: "network"},
		{name: "value"},
	}

	matched, err := assignToFlattenMapField(fieldVal, sf, parts, 0, "hello", "", func(parts []queryPart) []string {
		out := make([]string, len(parts))
		for i, part := range parts {
			out[i] = part.name
		}
		return out
	})
	if err != nil {
		t.Fatalf("assignToFlattenMapField(insert): %v", err)
	}
	if !matched {
		t.Fatal("assignToFlattenMapField should match flattened path")
	}
	if s.Items == nil || len(s.Items) != 1 {
		t.Fatalf("map not initialized: %#v", s.Items)
	}
	if _, ok := s.Items["menu_network"]; !ok {
		t.Fatalf("flattened key mismatch: %#v", s.Items)
	}
	if s.Items["menu_network"].Value != "hello" {
		t.Fatalf("flattened value = %#v, want hello", s.Items["menu_network"])
	}

	s.Items["menu_network"].Value = "updated"
	matched, err = assignToFlattenMapField(fieldVal, sf, parts, 0, "again", "", func(parts []queryPart) []string {
		out := make([]string, len(parts))
		for i, part := range parts {
			out[i] = part.name
		}
		return out
	})
	if err != nil {
		t.Fatalf("assignToFlattenMapField(update): %v", err)
	}
	if !matched {
		t.Fatal("assignToFlattenMapField should match existing key")
	}
	if s.Items["menu_network"].Value != "again" {
		t.Fatalf("flattened update = %#v, want again", s.Items["menu_network"])
	}
}
