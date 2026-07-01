package main

import "testing"

func TestGetValueAndSetValue(t *testing.T) {
	t.Parallel()

	type sample struct {
		Name string `ini:"name"`
		Meta struct {
			Count int `ini:"count"`
		} `ini:"meta"`
	}

	s := &sample{}
	if err := SetValue(s, "name", "Ryu"); err != nil {
		t.Fatalf("SetValue(name): %v", err)
	}
	if err := SetValue(s, "meta.count", 12); err != nil {
		t.Fatalf("SetValue(meta.count): %v", err)
	}

	if got, err := GetValue(s, "name"); err != nil || got != "Ryu" {
		t.Fatalf("GetValue(name) = (%v, %v), want (Ryu, nil)", got, err)
	}
	if got, err := GetValue(s, "meta.count"); err != nil || got != int64(12) {
		t.Fatalf("GetValue(meta.count) = (%v, %v), want (12, nil)", got, err)
	}
	if _, err := GetValue(s, "missing"); err == nil {
		t.Fatal("expected missing field lookup to fail")
	}
}
