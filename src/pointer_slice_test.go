package main

import "testing"

func TestPointerSliceReset_ClearsContentsAndPreservesCapacity(t *testing.T) {
	a, b := 1, 2
	s := []*int{&a, &b}
	origCap := cap(s)

	got := PointerSliceReset(s)
	if len(got) != 0 {
		t.Fatalf("PointerSliceReset len = %d, want 0", len(got))
	}
	if cap(got) != origCap {
		t.Fatalf("PointerSliceReset cap = %d, want %d", cap(got), origCap)
	}
}

func TestRecoverOrAppend_ReusesGhostAndAppendsWhenNeeded(t *testing.T) {
	cleared := 0
	clearFunc := func(v *int) {
		cleared++
		*v = 0
	}
	newValue := 99
	newFunc := func() *int { return &newValue }

	// RecoverOrAppend reuses a "ghost": a non-nil pointer still sitting in the
	// slice's reserved capacity past its length, left behind by an earlier
	// shrink. With nothing there it falls through to newFunc and appends.

	// Append path: capacity is reserved but the slot past len is nil.
	val := 7
	slice := make([]*int, 1, 2)
	slice[0] = &val
	got := RecoverOrAppend(&slice, clearFunc, newFunc)
	if got != &newValue {
		t.Fatalf("expected new value to be appended when no ghost, got %v", got)
	}
	if cleared != 0 {
		t.Fatalf("clearFunc should not have run on the append path, got %d", cleared)
	}
	if len(slice) != 2 {
		t.Fatalf("slice length after append = %d, want 2", len(slice))
	}

	// Recovery path: build a real ghost by filling both slots and then reslicing
	// down, which leaves the second pointer live in the reserved capacity.
	ghost := 5
	backing := make([]*int, 2, 2)
	backing[0] = &val
	backing[1] = &ghost
	slice = backing[:1]

	cleared = 0
	got = RecoverOrAppend(&slice, clearFunc, newFunc)
	if got != &ghost {
		t.Fatalf("expected the ghost to be recovered, got %v", got)
	}
	if cleared != 1 {
		t.Fatalf("clearFunc should run exactly once on the recovery path, got %d", cleared)
	}
	if ghost != 0 {
		t.Fatalf("clearFunc should have reset the recovered value, got %d", ghost)
	}
	if len(slice) != 2 {
		t.Fatalf("slice length after recovery = %d, want 2", len(slice))
	}
}
