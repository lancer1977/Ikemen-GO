package main

import (
	"reflect"
	"testing"
)

func TestSliceMoveIntAndHelpers_ReorderAndInsertValues(t *testing.T) {
	if got := sliceInsertInt([]int{1, 3}, 2, 1); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("sliceInsertInt() = %#v, want %#v", got, []int{1, 2, 3})
	}

	if got := sliceRemoveInt([]int{1, 2, 3}, 1); !reflect.DeepEqual(got, []int{1, 3}) {
		t.Fatalf("sliceRemoveInt() = %#v, want %#v", got, []int{1, 3})
	}

	if got := sliceMoveInt([]int{1, 2, 3, 4}, 1, 3); !reflect.DeepEqual(got, []int{1, 3, 4, 2}) {
		t.Fatalf("sliceMoveInt() = %#v, want %#v", got, []int{1, 3, 4, 2})
	}

	orig := []int{5, 6, 7}
	if got := sliceMoveInt(orig, 2, 2); !reflect.DeepEqual(got, orig) {
		t.Fatalf("sliceMoveInt(no-op) = %#v, want %#v", got, orig)
	}
}
