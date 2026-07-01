package main

import "testing"

func TestCosAndSin_HandleCanonicalAngles(t *testing.T) {
	if got := Cos(0); got != 1 {
		t.Fatalf("Cos(0) = %v, want 1", got)
	}
	if got := Sin(0); got != 0 {
		t.Fatalf("Sin(0) = %v, want 0", got)
	}
	if got := Cos(Rad(90)); got > 1e-6 || got < -1e-6 {
		t.Fatalf("Cos(90deg) = %v, want near 0", got)
	}
	if got := Sin(Rad(90)); got != 1 {
		t.Fatalf("Sin(90deg) = %v, want 1", got)
	}
	if got := Cos(Rad(180)); got != -1 {
		t.Fatalf("Cos(180deg) = %v, want -1", got)
	}
}
