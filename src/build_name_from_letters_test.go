package main

import "testing"

func TestBuildNameFromLetters(t *testing.T) {
	t.Parallel()

	mo := &Motif{
		HiscoreInfo: HiscoreInfoProperties{
			Glyphs: []string{"A", ">", "C"},
		},
	}

	if got := buildNameFromLetters(mo, []int{1, 2, 3}); got != "A C" {
		t.Fatalf("buildNameFromLetters = %q, want %q", got, "A C")
	}
	if got := buildNameFromLetters(mo, []int{0, 4, 1}); got != "A" {
		t.Fatalf("buildNameFromLetters out-of-range = %q, want %q", got, "A")
	}
}
