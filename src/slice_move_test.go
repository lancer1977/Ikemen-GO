package main

import (
	"reflect"
	"testing"
)

func TestSliceMove_ReordersElementsAndPreservesNoOpMoves(t *testing.T) {
	got := sliceMove([]string{"a", "b", "c", "d"}, 1, 3)
	want := []string{"a", "c", "d", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sliceMove(forward) = %#v, want %#v", got, want)
	}

	got = sliceMove([]string{"a", "b", "c", "d"}, 3, 1)
	want = []string{"a", "d", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sliceMove(backward) = %#v, want %#v", got, want)
	}

	orig := []string{"x", "y"}
	got = sliceMove(orig, 1, 1)
	if !reflect.DeepEqual(got, orig) {
		t.Fatalf("sliceMove(no-op) = %#v, want %#v", got, orig)
	}
}
