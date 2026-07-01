package main

import "testing"

func TestClamp_CoversLowMidHighAcrossTypes(t *testing.T) {
	if got := Clamp(3, 5, 9); got != 5 {
		t.Fatalf("Clamp(int low) = %v, want %v", got, 5)
	}
	if got := Clamp(7, 5, 9); got != 7 {
		t.Fatalf("Clamp(int mid) = %v, want %v", got, 7)
	}
	if got := Clamp(11, 5, 9); got != 9 {
		t.Fatalf("Clamp(int high) = %v, want %v", got, 9)
	}

	if got := Clamp(float32(1.5), 2.0, 4.0); got != 2.0 {
		t.Fatalf("Clamp(float low) = %v, want %v", got, 2.0)
	}
	if got := Clamp(float32(3.0), 2.0, 4.0); got != 3.0 {
		t.Fatalf("Clamp(float mid) = %v, want %v", got, 3.0)
	}
	if got := Clamp(float32(5.0), 2.0, 4.0); got != 4.0 {
		t.Fatalf("Clamp(float high) = %v, want %v", got, 4.0)
	}
}
