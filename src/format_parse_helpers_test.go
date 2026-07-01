package main

import "testing"

func TestAtobAndItoa_HandleCommonInputs(t *testing.T) {
	if !Atob("true") {
		t.Fatalf("expected true to parse as true")
	}
	if Atob("not-a-bool") {
		t.Fatalf("expected invalid bool to parse as false")
	}
	if got := Itoa(-12345); got != "-12345" {
		t.Fatalf("unexpected Itoa output: %q", got)
	}
}
