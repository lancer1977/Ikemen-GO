package main

import "testing"

func TestCompareNatural(t *testing.T) {
	t.Parallel()

	if !compareNatural("file2", "file10") {
		t.Fatal("expected file2 to sort before file10")
	}
	if compareNatural("file10", "file2") {
		t.Fatal("expected file10 to sort after file2")
	}
	if !compareNatural("alpha", "alpha") {
		t.Fatal("expected equal strings to compare as true for stable ordering")
	}
}
