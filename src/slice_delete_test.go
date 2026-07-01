package main

import (
	"reflect"
	"testing"
)

func TestSliceDelete_RemovesElementsAndLeavesInvalidIndicesUntouched(t *testing.T) {
	got := SliceDelete([]string{"a", "b", "c"}, 1)
	want := []string{"a", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SliceDelete(middle) = %#v, want %#v", got, want)
	}

	got = SliceDelete([]string{"a", "b", "c"}, 0)
	want = []string{"b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SliceDelete(first) = %#v, want %#v", got, want)
	}

	got = SliceDelete([]string{"a", "b", "c"}, 2)
	want = []string{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SliceDelete(last) = %#v, want %#v", got, want)
	}

	orig := []string{"x", "y"}
	got = SliceDelete(orig, -1)
	if !reflect.DeepEqual(got, orig) {
		t.Fatalf("SliceDelete(invalid low) = %#v, want %#v", got, orig)
	}
	got = SliceDelete(orig, 2)
	if !reflect.DeepEqual(got, orig) {
		t.Fatalf("SliceDelete(invalid high) = %#v, want %#v", got, orig)
	}
}
