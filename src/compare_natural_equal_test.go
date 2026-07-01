package main

import "testing"

func TestCompareNatural_ReturnsFalseForEqualNaturalKeys(t *testing.T) {
	if compareNatural("file2", "file2") {
		t.Fatal("expected equal natural keys to not compare less than each other")
	}
	if compareNatural("alpha", "alpha") {
		t.Fatal("expected equal non-numeric keys to not compare less than each other")
	}
}
