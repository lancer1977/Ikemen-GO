package main

import "testing"

func TestParseQueryPathTreatsUnclosedBracketAsPlainName(t *testing.T) {
	parts := parseQueryPath("root.child[2.leaf")
	if len(parts) != 2 {
		t.Fatalf("parseQueryPath length = %d, want 2", len(parts))
	}
	if parts[1].name != "child[2" || parts[1].index != nil {
		t.Fatalf("parseQueryPath malformed part = %#v, want plain name child[2", parts[1])
	}
}
