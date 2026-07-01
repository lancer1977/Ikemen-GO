package main

import (
	"reflect"
	"testing"
)

func TestSortedKeys_UsesNaturalOrderingForMapKeys(t *testing.T) {
	got := SortedKeys(map[string]int{
		"file10": 1,
		"file2":  2,
		"file1":  3,
	})
	want := []string{"file1", "file2", "file10"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SortedKeys() = %#v, want %#v", got, want)
	}
}
