package main

import "testing"

func TestAbs_WorksForIntegersAndFloats(t *testing.T) {
	if got := Abs(-7); got != 7 {
		t.Fatalf("Abs(int negative) = %v, want 7", got)
	}
	if got := Abs(0); got != 0 {
		t.Fatalf("Abs(int zero) = %v, want 0", got)
	}
	if got := Abs(9); got != 9 {
		t.Fatalf("Abs(int positive) = %v, want 9", got)
	}

	if got := Abs(float32(-3.5)); got != 3.5 {
		t.Fatalf("Abs(float negative) = %v, want 3.5", got)
	}
	if got := Abs(float32(0)); got != 0 {
		t.Fatalf("Abs(float zero) = %v, want 0", got)
	}
	if got := Abs(float32(4.2)); got != 4.2 {
		t.Fatalf("Abs(float positive) = %v, want 4.2", got)
	}
}
