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

	// First test: append when no ghost available (reserved capacity is empty/nil)
	// When a slice is created with reserved capacity but no ghost is present,
	// RecoverOrAppend appends a new item
	val := 7
	slice := make([]*int, 1, 2)
	slice[0] = &val
	got := RecoverOrAppend(&slice, clearFunc, newFunc)
	if got != &newValue {
		t.Fatalf("expected new value to be appended when no ghost, got %v", got)
	}
	if cleared != 0 {
		t.Fatalf("clearFunc should not have run on append path, got %d", cleared)
	}
	if len(slice) != 2 {
		t.Fatalf("slice length after append = %d, want 2", len(slice))
	}

	// Second test: recover a valid ghost from reserved capacity
	// Create slice with an item in reserved capacity, then shrink and grow
	ghost := 5
	slice = make([]*int, 1, 2)
	slice[0] = &val
	// Manually create a ghost by using the underlying array
	// We need to put a non-nil value in position 1 before shrinking
	ghostVal := 99
	oldSlice := slice
	slice = append(slice, &ghostVal)  // len=2, cap=2, [0]=&val, [1]=&ghostVal
	slice = oldSlice              // Reset to len=1, but underlying array still has &ghostVal at [1]

	// Actually, append creates a new reference. Let me use a direct approach instead.
	// Create the ghost by direct array manipulation through a larger slice
	largeSlice := make([]*int, 2, 2)
	largeSlice[0] = &val
	largeSlice[1] = &ghost
	slice = largeSlice[:1]  // Now len=1, cap=2, with &ghost as ghost at [1]

	cleared = 0
	got = RecoverOrAppend(&slice, clearFunc, newFunc)
	if got != &ghost {
		t.Fatalf("expected ghost value to be recovered, got %v", got)
	}
	if cleared != 1 || ghost != 0 {
		t.Fatalf("clearFunc not applied correctly: cleared=%d ghost=%d", cleared, ghost)
	}
	if len(slice) != 2 {
		t.Fatalf("slice length after recovery = %d, want 2", len(slice))
	}
}
