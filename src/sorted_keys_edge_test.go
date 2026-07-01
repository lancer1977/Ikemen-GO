package main

import (
	"reflect"
	"testing"
)

func TestSortedKeys_HandlesEmptyAndSingleElementMaps(t *testing.T) {
	if got := SortedKeys(map[string]int{}); len(got) != 0 {
		t.Fatalf("SortedKeys(empty) = %#v, want empty", got)
	}

	got := SortedKeys(map[string]int{"only": 1})
	want := []string{"only"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SortedKeys(single) = %#v, want %#v", got, want)
	}
}
