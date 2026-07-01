package main

import "testing"

func TestBtoi_ReturnsOneForTrueAndZeroForFalse(t *testing.T) {
	if got := Btoi(true); got != 1 {
		t.Fatalf("Btoi(true) = %v, want 1", got)
	}
	if got := Btoi(false); got != 0 {
		t.Fatalf("Btoi(false) = %v, want 0", got)
	}
}
