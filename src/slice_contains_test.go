package main

import "testing"

func TestSliceContains_HandlesLowercaseMatching(t *testing.T) {
	items := []string{"Alpha", "Bravo"}
	if !sliceContains(items, "alpha", true) {
		t.Fatal("expected lowercase search to match ignoring case")
	}
	if sliceContains(items, "charlie", true) {
		t.Fatal("expected missing lowercase search to fail")
	}
	if !sliceContains(items, "Bravo", false) {
		t.Fatal("expected exact search to match")
	}
}
