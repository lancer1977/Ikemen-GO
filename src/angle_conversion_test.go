package main

import "testing"

func TestRadAndDeg_RoundTripAndCanonicalAnglesAgain(t *testing.T) {
	for _, deg := range []float32{0, 90, 180, 270, 360} {
		if got := Deg(Rad(deg)); got != deg {
			t.Fatalf("Deg(Rad(%v)) = %v, want %v", deg, got, deg)
		}
	}

	if got := Rad(180); got != 3.1415927 {
		t.Fatalf("Rad(180) = %v, want %v", got, 3.1415927)
	}
	if got := Deg(3.1415927); got == 0 {
		t.Fatalf("Deg(pi) should not be zero")
	}
}
