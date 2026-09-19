package main

import (
	"reflect"
	"testing"
)

// sliceMove is the generic replacement for the former sliceInsertInt /
// sliceRemoveInt / sliceMoveInt helpers (now commented out in common.go).
// These cases carry over the reordering expectations of the old
// sliceMoveInt tests.
func TestSliceMove_ReordersValues(t *testing.T) {
	if got := sliceMove([]int{1, 2, 3, 4}, 1, 3); !reflect.DeepEqual(got, []int{1, 3, 4, 2}) {
		t.Fatalf("sliceMove() = %#v, want %#v", got, []int{1, 3, 4, 2})
	}

	if got := sliceMove([]int{1, 2, 3, 4}, 3, 0); !reflect.DeepEqual(got, []int{4, 1, 2, 3}) {
		t.Fatalf("sliceMove() backwards = %#v, want %#v", got, []int{4, 1, 2, 3})
	}

	orig := []int{5, 6, 7}
	if got := sliceMove(orig, 2, 2); !reflect.DeepEqual(got, []int{5, 6, 7}) {
		t.Fatalf("sliceMove(no-op) = %#v, want %#v", got, []int{5, 6, 7})
	}
}

// The generic accepts any element type, not just int.
func TestSliceMove_NonIntElements(t *testing.T) {
	if got := sliceMove([]string{"a", "b", "c"}, 0, 2); !reflect.DeepEqual(got, []string{"b", "c", "a"}) {
		t.Fatalf("sliceMove(strings) = %#v, want %#v", got, []string{"b", "c", "a"})
	}
}
