package main

import "testing"

func TestParseQueryPathTreatsUnclosedBracketAsPlainName(t *testing.T) {
	// parseQueryPath splits on "." first, so "root.child[2.leaf" becomes ["root", "child[2", "leaf"].
	// An unclosed bracket does not match the indexed pattern, so it becomes a plain name part.
	parts := parseQueryPath("root.child[2.leaf")
	if len(parts) != 3 {
		t.Fatalf("parseQueryPath length = %d, want 3", len(parts))
	}
	if parts[0].name != "root" || parts[0].index != nil {
		t.Fatalf("parseQueryPath part[0] = %#v, want plain name root", parts[0])
	}
	if parts[1].name != "child[2" || parts[1].index != nil {
		t.Fatalf("parseQueryPath part[1] = %#v, want plain name child[2", parts[1])
	}
	if parts[2].name != "leaf" || parts[2].index != nil {
		t.Fatalf("parseQueryPath part[2] = %#v, want plain name leaf", parts[2])
	}
}
