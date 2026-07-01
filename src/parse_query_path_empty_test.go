package main

import "testing"

func TestParseQueryPathReturnsSingleEmptyPartForBlankInput(t *testing.T) {
	parts := parseQueryPath("")
	if len(parts) != 1 {
		t.Fatalf("parseQueryPath empty length = %d, want 1", len(parts))
	}
	if parts[0].name != "" || parts[0].index != nil {
		t.Fatalf("parseQueryPath empty part = %#v, want zero value", parts[0])
	}
}
