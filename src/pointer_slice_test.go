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

	recovered := 7
	slice := make([]*int, 1, 2)
	slice[0] = &recovered
	got := RecoverOrAppend(&slice, clearFunc, newFunc)
	if got != &recovered {
		t.Fatal("expected ghost value to be recovered")
	}
	if cleared != 1 || recovered != 0 {
		t.Fatalf("clearFunc not applied correctly: cleared=%d recovered=%d", cleared, recovered)
	}
	if len(slice) != 2 {
		t.Fatalf("slice length after recovery = %d, want 2", len(slice))
	}

	slice = []*int{&recovered}
	cleared = 0
	got = RecoverOrAppend(&slice, clearFunc, newFunc)
	if got != &newValue {
		t.Fatal("expected new value to be appended when no ghost available")
	}
	if cleared != 0 {
		t.Fatalf("clearFunc should not have run on append path, got %d", cleared)
	}
	if len(slice) != 2 {
		t.Fatalf("slice length after append = %d, want 2", len(slice))
	}
}
