package main

import "testing"

func TestSignAndAbs_WorkAcrossIntegerAndFloatTypes(t *testing.T) {
	if got := Sign(-3); got != -1 {
		t.Fatalf("unexpected int sign: %v", got)
	}
	if got := Sign(0); got != 0 {
		t.Fatalf("unexpected zero sign: %v", got)
	}
	if got := Sign(int64(7)); got != 1 {
		t.Fatalf("unexpected int64 sign: %v", got)
	}

	if got := Abs(-3); got != 3 {
		t.Fatalf("unexpected int abs: %v", got)
	}
	if got := Abs(int32(-12)); got != 12 {
		t.Fatalf("unexpected int32 abs: %v", got)
	}
	if got := Abs(float32(-1.5)); got != 1.5 {
		t.Fatalf("unexpected float32 abs: %v", got)
	}
	if got := Abs(float64(2.5)); got != 2.5 {
		t.Fatalf("unexpected float64 abs: %v", got)
	}
}
