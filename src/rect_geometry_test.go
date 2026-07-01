package main

import (
	"math"
	"testing"
)

func TestRectRotate_RotatesAroundCenterAgain(t *testing.T) {
	got := RectRotate(0, 0, 2, 1, 1, 0.5, Rad(90))
	want := [][2]float32{
		{1.5, -0.5},
		{1.5, 1.5},
		{0.5, 1.5},
		{0.5, -0.5},
	}
	for i := range want {
		if math.Abs(float64(got[i][0]-want[i][0])) > 1e-6 || math.Abs(float64(got[i][1]-want[i][1])) > 1e-6 {
			t.Fatalf("RectRotate() corner %d = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestRectIntersect_DetectsOverlapAndSeparationAgain(t *testing.T) {
	if !RectIntersect(0, 0, 2, 2, 1, 1, 2, 2, 1, 1, 2, 2, 0, 0) {
		t.Fatal("RectIntersect should report overlapping rectangles")
	}
	if RectIntersect(0, 0, 2, 2, 5, 5, 2, 2, 1, 1, 6, 6, 0, 0) {
		t.Fatal("RectIntersect should report separated rectangles as non-intersecting")
	}
}
