package main

import "testing"

func TestMinAndMax_WorkAcrossIntegerAndFloatTypes(t *testing.T) {
	if got := Min(3, 1, 2); got != 1 {
		t.Fatalf("unexpected int min: %v", got)
	}
	if got := Max(int64(-5), int64(9), int64(2)); got != 9 {
		t.Fatalf("unexpected int64 max: %v", got)
	}
	if got := Min(float32(3.5), float32(1.25), float32(2.75)); got != 1.25 {
		t.Fatalf("unexpected float32 min: %v", got)
	}
	if got := Max(float64(-2.5), float64(4.75), float64(0.125)); got != 4.75 {
		t.Fatalf("unexpected float64 max: %v", got)
	}
}

func TestMinAndMax_ReturnSingleArgument(t *testing.T) {
	if got := Min(42); got != 42 {
		t.Fatalf("unexpected single-arg min: %v", got)
	}
	if got := Max(42); got != 42 {
		t.Fatalf("unexpected single-arg max: %v", got)
	}
}
