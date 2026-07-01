package main

import "testing"

func TestUniqueStrings(t *testing.T) {
	t.Parallel()

	got := uniqueStrings([]string{"a", "", "b", "a", "c", "b", ""})
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("uniqueStrings len = %d, want %d (%#v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("uniqueStrings[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
