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
	// compareNatural is a less-than comparator for use in sort.Slice.
	// For stable sorting, equal strings must return false (not true).
	// If equal strings returned true, the sort would be unstable (claiming a < a).
	if compareNatural("alpha", "alpha") {
		t.Fatal("equal strings should compare as false for stable sorting")
	}
}
