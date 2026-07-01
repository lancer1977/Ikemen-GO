package main

import "testing"

func TestRectRotate_RotatesAroundCenter(t *testing.T) {
	points := RectRotate(0, 0, 2, 2, 1, 1, 90)
	if len(points) != 4 {
		t.Fatalf("unexpected point count: %d", len(points))
	}
	if points[0] != [2]float32{2, 0} || points[2] != [2]float32{0, 2} {
		t.Fatalf("unexpected rotated rectangle points: %#v", points)
	}
}

func TestRectIntersect_DetectsOverlapAndSeparation(t *testing.T) {
	if !RectIntersect(0, 0, 2, 2, 1, 1, 2, 2, 1, 1, 2, 2, 0, 0) {
		t.Fatalf("expected overlapping rectangles to intersect")
	}
	if RectIntersect(0, 0, 2, 2, 5, 5, 2, 2, 1, 1, 6, 6, 0, 0) {
		t.Fatalf("expected separated rectangles to not intersect")
	}
}

func TestF64toI32_ClampsToInt32Bounds(t *testing.T) {
	if got := F64toI32(float64(2147483648)); got != 2147483647 {
		t.Fatalf("unexpected positive clamp: %v", got)
	}
	if got := F64toI32(float64(-2147483649)); got != -2147483648 {
		t.Fatalf("unexpected negative clamp: %v", got)
	}
	if got := F64toI32(123.75); got != 123 {
		t.Fatalf("unexpected passthrough conversion: %v", got)
	}
}
