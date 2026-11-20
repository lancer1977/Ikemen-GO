package main

import "testing"

func TestRectRotate_RotatesAroundCenter(t *testing.T) {
	// RectRotate expects angle in radians; Rad() converts degrees to radians
	points := RectRotate(0, 0, 2, 2, 1, 1, Rad(90))
	if len(points) != 4 {
		t.Fatalf("unexpected point count: %d", len(points))
	}
	const tolerance = 1e-6
	p0x, p0y := points[0][0], points[0][1]
	p2x, p2y := points[2][0], points[2][1]
	if Abs(p0x-2) > tolerance || Abs(p0y-0) > tolerance || Abs(p2x-0) > tolerance || Abs(p2y-2) > tolerance {
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
