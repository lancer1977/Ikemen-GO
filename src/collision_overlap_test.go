package main

import "testing"

func TestZAxisOverlapHonorsEnabledBounds(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.zmin = 0
	sys.zmax = 1
	if !sys.zAxisOverlap(0, 1, 1, 1, 0.5, 1, 1, 1) {
		t.Fatal("zAxisOverlap() should report overlapping ranges when Z is enabled")
	}
	if sys.zAxisOverlap(0, 1, 1, 1, 5, 1, 1, 1) {
		t.Fatal("zAxisOverlap() should report separated ranges as non-overlapping")
	}

	sys.zmax = 0
	if !sys.zAxisOverlap(0, 1, 1, 1, 5, 1, 1, 1) {
		t.Fatal("zAxisOverlap() should ignore Z separation when Z is disabled")
	}
}

func TestCollisionOverlapHandlesAxisAlignedAndRotatedBoxes(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys

	box1 := [][4]float32{{-1, -1, 1, 1}}
	box2 := [][4]float32{{-1, -1, 1, 1}}
	if !sys.clsnOverlap(box1, [2]float32{1, 1}, [2]float32{0, 0}, 1, 0, box2, [2]float32{1, 1}, [2]float32{0.5, 0}, -1, 0) {
		t.Fatal("clsnOverlap() should report overlapping axis-aligned boxes")
	}
	if sys.clsnOverlap(box1, [2]float32{1, 1}, [2]float32{0, 0}, 1, 0, box2, [2]float32{1, 1}, [2]float32{5, 5}, -1, 0) {
		t.Fatal("clsnOverlap() should report separated axis-aligned boxes as non-overlapping")
	}
	if !sys.clsnOverlap(box1, [2]float32{1, 1}, [2]float32{0, 0}, 1, 45, box2, [2]float32{1, 1}, [2]float32{0, 0}, -1, 0) {
		t.Fatal("clsnOverlap() should still detect overlap when one side is rotated")
	}
}
