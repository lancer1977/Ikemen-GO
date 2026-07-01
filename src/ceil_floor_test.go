package main

import "testing"

func TestPowLerpCeilAndFloor_ReturnExpectedResultsLegacy(t *testing.T) {
	if got := Pow(2, 3); got != 8 {
		t.Fatalf("Pow() = %v, want 8", got)
	}
	if got := Lerp(10, 20, 0.25); got != 12.5 {
		t.Fatalf("Lerp() = %v, want 12.5", got)
	}

	if got := Ceil(1.2); got != 2 {
		t.Fatalf("Ceil() = %v, want 2", got)
	}
	if got := Ceil(-1.2); got != -1 {
		t.Fatalf("Ceil() = %v, want -1", got)
	}
	if got := Floor(1.8); got != 1 {
		t.Fatalf("Floor() = %v, want 1", got)
	}
	if got := Floor(-1.2); got != -2 {
		t.Fatalf("Floor() = %v, want -2", got)
	}
}
