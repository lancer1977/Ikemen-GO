package main

import "testing"

func TestResolveLifebarPrefixOperators_ReturnsInputWhenNoOperatorsArePresent(t *testing.T) {
	is := IniSection{
		"p1.offset": "10, 20",
		"plain":     "value",
	}

	got := resolveLifebarPrefixOperators(is, "p1.", "p")
	if got["p1.offset"] != "10, 20" || got["plain"] != "value" {
		t.Fatalf("resolveLifebarPrefixOperators returned %#v, want unchanged input", got)
	}
	if &got == &is {
		t.Fatal("resolveLifebarPrefixOperators should return the original map value, not an aliasable pointer check")
	}
}
