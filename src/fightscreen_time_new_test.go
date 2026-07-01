package main

import "testing"

func TestNewFightScreenTime(t *testing.T) {
	t.Parallel()

	ti := newFightScreenTime()
	if ti == nil {
		t.Fatal("newFightScreenTime returned nil")
	}
	if ti.counter == nil {
		t.Fatal("newFightScreenTime should allocate counter map")
	}
	if ti.framespercount != 60 {
		t.Fatalf("framespercount = %d, want 60", ti.framespercount)
	}
}
