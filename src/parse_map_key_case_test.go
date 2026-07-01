package main

import "testing"

func TestParseMapKeyAcceptsMixedCasePrefix(t *testing.T) {
	name, ok := parseMapKey(" MAP.  stage ")
	if !ok || name != "stage" {
		t.Fatalf("parseMapKey mixed-case prefix = %q %v, want stage true", name, ok)
	}
}
