package main

import "testing"

func TestRoundFloat_RoundsToRequestedPrecisionLegacy(t *testing.T) {
	if got := RoundFloat(12.3456, 2); got != 12.35 {
		t.Fatalf("RoundFloat() = %v, want %v", got, 12.35)
	}
	if got := RoundFloat(12.5, 0); got != 13 {
		t.Fatalf("RoundFloat() = %v, want %v", got, 13)
	}
	if got := RoundFloat(-12.5, 0); got != -13 {
		t.Fatalf("RoundFloat() = %v, want %v", got, -13)
	}
}
