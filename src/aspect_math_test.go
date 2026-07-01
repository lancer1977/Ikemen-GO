package main

import "testing"

func TestCalculateAspect_CastsIntegerAndFloatInputsConsistently(t *testing.T) {
	if got := CalculateAspect(320, 240); got != 4.0/3.0 {
		t.Fatalf("unexpected int aspect: %v", got)
	}
	if got := CalculateAspect(int32(16), int32(9)); got != 16.0/9.0 {
		t.Fatalf("unexpected int32 aspect: %v", got)
	}
	if got := CalculateAspect(float32(21), float32(9)); got != 21.0/9.0 {
		t.Fatalf("unexpected float32 aspect: %v", got)
	}
}
