package main

import "testing"

func TestSplitAndTrim_SplitsAndTrimsEachSegment(t *testing.T) {
	got := SplitAndTrim("  a , , b  ,c ", ",")
	want := []string{"a", "", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("unexpected segment count: %#v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected split result at %d: got %q want %q", i, got[i], want[i])
		}
	}
}
