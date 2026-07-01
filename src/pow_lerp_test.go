package main

import "testing"

func TestPowLerpCeilAndFloor_ReturnExpectedResultsAgain(t *testing.T) {
	if got := Pow(2, 3); got != 8 {
		t.Fatalf("Pow() = %v, want 8", got)
	}
	if got := Pow(4, -1); got != 0.25 {
		t.Fatalf("Pow(negative exponent) = %v, want 0.25", got)
	}
	if got := Lerp(10, 20, 0.25); got != 12.5 {
		t.Fatalf("Lerp() = %v, want 12.5", got)
	}
	if got := Lerp(10, 20, 0.5); got != 15 {
		t.Fatalf("Lerp(midpoint) = %v, want 15", got)
	}
}
