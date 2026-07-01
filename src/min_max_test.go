package main

import "testing"

func TestMinAndMax_WorkAcrossIntegerAndFloatTypesAgain(t *testing.T) {
	if got := Min(5, 2, 9, 2); got != 2 {
		t.Fatalf("Min(int) = %v, want 2", got)
	}
	if got := Max(5, 2, 9, 2); got != 9 {
		t.Fatalf("Max(int) = %v, want 9", got)
	}

	if got := Min(float32(4.5), float32(-1.25), float32(0)); got != -1.25 {
		t.Fatalf("Min(float) = %v, want -1.25", got)
	}
	if got := Max(float32(4.5), float32(-1.25), float32(0)); got != 4.5 {
		t.Fatalf("Max(float) = %v, want 4.5", got)
	}
}

func TestMinAndMax_ReturnSingleArgumentAgain(t *testing.T) {
	if got := Min(7); got != 7 {
		t.Fatalf("Min(single) = %v, want 7", got)
	}
	if got := Max(7); got != 7 {
		t.Fatalf("Max(single) = %v, want 7", got)
	}
}
