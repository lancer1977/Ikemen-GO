package main

import "testing"

func TestRoundFloat_RoundsToRequestedPrecisionAgain(t *testing.T) {
	if got := RoundFloat(12.3456, 2); got != 12.35 {
		t.Fatalf("RoundFloat(2) = %v, want %v", got, 12.35)
	}
	if got := RoundFloat(12.5, 0); got != 13 {
		t.Fatalf("RoundFloat(0) = %v, want %v", got, 13)
	}
	if got := RoundFloat(-12.5, 0); got != -13 {
		t.Fatalf("RoundFloat(0 neg) = %v, want %v", got, -13)
	}
	if got := RoundFloat(1234.0, -2); got != 1200 {
		t.Fatalf("RoundFloat(-2) = %v, want 1200", got)
	}
}
