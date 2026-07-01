package main

import "testing"

func TestCalculateAspect_CastsIntegerAndFloatInputsConsistentlyAgain(t *testing.T) {
	if got := CalculateAspect(16, 9); got != 16.0/9.0 {
		t.Fatalf("CalculateAspect(int) = %v, want %v", got, 16.0/9.0)
	}
	if got := CalculateAspect(float32(4), float32(3)); got != 4.0/3.0 {
		t.Fatalf("CalculateAspect(float32) = %v, want %v", got, 4.0/3.0)
	}
	if got := CalculateAspect(float64(21), float64(9)); got != 21.0/9.0 {
		t.Fatalf("CalculateAspect(float64) = %v, want %v", got, 21.0/9.0)
	}
}
