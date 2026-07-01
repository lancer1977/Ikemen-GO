package main

import "testing"

func TestBuildLookupDirs(t *testing.T) {
	t.Parallel()

	if got := buildLookupDirs("", "base"); got != nil {
		t.Fatalf("empty lookup tag = %#v, want nil", got)
	}

	got := buildLookupDirs("def,, data/ , alt", "base")
	want := []string{"base", "", "data/", "alt"}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (%#v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
