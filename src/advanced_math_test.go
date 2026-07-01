package main

import "testing"

func TestPowLerpCeilAndFloor_ReturnExpectedResults(t *testing.T) {
	if got := Pow(2, 3); got != 8 {
		t.Fatalf("unexpected pow result: %v", got)
	}
	if got := Lerp(10, 20, 0.25); got != 12.5 {
		t.Fatalf("unexpected lerp result: %v", got)
	}
	if got := Ceil(1.1); got != 2 {
		t.Fatalf("unexpected ceil result: %v", got)
	}
	if got := Ceil(-1.1); got != -1 {
		t.Fatalf("unexpected negative ceil result: %v", got)
	}
	if got := Floor(1.9); got != 1 {
		t.Fatalf("unexpected floor result: %v", got)
	}
	if got := Floor(-1.1); got != -2 {
		t.Fatalf("unexpected negative floor result: %v", got)
	}
}
