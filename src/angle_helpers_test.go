package main

import (
	"math"
	"testing"
)

func TestRadAndDeg_RoundTripAndCanonicalAngles(t *testing.T) {
	if got := Rad(180); math.Abs(float64(got-math.Pi)) > 1e-6 {
		t.Fatalf("unexpected radians for 180 degrees: %v", got)
	}
	if got := Deg(float32(math.Pi)); math.Abs(float64(got-180)) > 1e-5 {
		t.Fatalf("unexpected degrees for pi radians: %v", got)
	}
	if got := Deg(Rad(90)); math.Abs(float64(got-90)) > 1e-5 {
		t.Fatalf("unexpected round-trip degrees for 90: %v", got)
	}
}
