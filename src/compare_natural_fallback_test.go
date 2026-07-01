package main

import "testing"

func TestCompareNatural_FallsBackToLexicographicForNonMatchingKeys(t *testing.T) {
	if !compareNatural("alpha", "beta!") {
		t.Fatal("expected non-matching keys to fall back to lexicographic ordering")
	}
	if compareNatural("beta!", "alpha") {
		t.Fatal("expected reverse non-matching order to be false")
	}
}
