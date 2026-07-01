package main

import "testing"

func TestNormalizeAxes(t *testing.T) {
	axes := [6]int8{-128, -64, 0, 64, 127, 1}
	got := NormalizeAxes(&axes)

	if got[0] != -1 {
		t.Fatalf("NormalizeAxes[0] = %v, want -1", got[0])
	}
	if got[1] >= 0 {
		t.Fatalf("NormalizeAxes[1] = %v, want negative", got[1])
	}
	if got[2] != 0 {
		t.Fatalf("NormalizeAxes[2] = %v, want 0", got[2])
	}
	if got[3] <= 0 {
		t.Fatalf("NormalizeAxes[3] = %v, want positive", got[3])
	}
	if got[4] != 1 {
		t.Fatalf("NormalizeAxes[4] = %v, want 1", got[4])
	}
	if got[5] <= 0 {
		t.Fatalf("NormalizeAxes[5] = %v, want positive", got[5])
	}
}
