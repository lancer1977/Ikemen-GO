package main

import "testing"

func TestSign_ReturnsNegativeZeroAndPositiveForIntsAndFloats(t *testing.T) {
	if got := Sign(-7); got != -1 {
		t.Fatalf("Sign(int negative) = %v, want -1", got)
	}
	if got := Sign(0); got != 0 {
		t.Fatalf("Sign(int zero) = %v, want 0", got)
	}
	if got := Sign(9); got != 1 {
		t.Fatalf("Sign(int positive) = %v, want 1", got)
	}

	if got := Sign(float32(-3.5)); got != -1 {
		t.Fatalf("Sign(float negative) = %v, want -1", got)
	}
	if got := Sign(float32(0)); got != 0 {
		t.Fatalf("Sign(float zero) = %v, want 0", got)
	}
	if got := Sign(float32(4.2)); got != 1 {
		t.Fatalf("Sign(float positive) = %v, want 1", got)
	}
}
