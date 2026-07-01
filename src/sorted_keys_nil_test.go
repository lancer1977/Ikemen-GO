package main

import "testing"

func TestSortedKeys_HandlesNilMap(t *testing.T) {
	var m map[string]int
	got := SortedKeys(m)
	if len(got) != 0 {
		t.Fatalf("SortedKeys(nil) = %#v, want empty", got)
	}
}
