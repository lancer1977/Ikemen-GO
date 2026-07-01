package main

import "testing"

func TestNewInputBuffer(t *testing.T) {
	t.Parallel()

	ib := NewInputBuffer()
	if ib == nil {
		t.Fatal("NewInputBuffer returned nil")
	}
	if ib.InputReader == nil {
		t.Fatal("expected InputReader to be initialized")
	}
	if ib.Bb != 0 || ib.Db != 0 || ib.Nb != 0 || ib.ab != 0 || ib.mp != 0 {
		t.Fatalf("expected zeroed buffer state, got %#v", ib)
	}
}
