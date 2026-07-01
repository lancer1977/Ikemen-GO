package main

import "testing"

func TestResolveLifebarPrefixOperators_RewritesWildcardAndGroupedSelectors(t *testing.T) {
	is := IniSection{
		"p*.font":              "wild",
		"p1|p3|p5.fgcolor":     "group",
		"p2.background":        "other",
		"plain":                "plain",
		" P1 | P2 .offset ":    "group2",
		"notprefixed.value":    "stay",
		"p1|p3|p5.bad spacing": "bad",
	}

	got := resolveLifebarPrefixOperators(is, "p1.", "p")
	if got["p1.font"] != "wild" {
		t.Fatalf("expected wildcard selector to map to p1.font, got %#v", got["p1.font"])
	}
	if got["p1.fgcolor"] != "group" {
		t.Fatalf("expected grouped selector to map to p1.fgcolor, got %#v", got["p1.fgcolor"])
	}
	if got["p2.background"] != "other" {
		t.Fatalf("expected unrelated selector to remain unchanged, got %#v", got["p2.background"])
	}
	if got["plain"] != "plain" {
		t.Fatalf("expected plain key to remain unchanged, got %#v", got["plain"])
	}
	if got["p1.offset"] != "group2" {
		t.Fatalf("expected spaced grouped selector to map to p1.offset, got %#v", got["p1.offset"])
	}
	if got["notprefixed.value"] != "stay" {
		t.Fatalf("expected non-target key to remain unchanged, got %#v", got["notprefixed.value"])
	}
}
