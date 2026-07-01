package main

import "testing"

func TestTextSpriteMutators_UpdateLocalcoordPositionScaleAndWindow(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.gameWidth = 640
	sys.gameHeight = 480
	sys.widthScale = 2
	sys.heightScale = 3

	ts := NewTextSprite()
	ts.fnt = &Fnt{}

	ts.SetLocalcoord(640, 480)
	if ts.localScale != 0.5 || ts.offsetX != 0 {
		t.Fatalf("SetLocalcoord() = scale %v offset %d", ts.localScale, ts.offsetX)
	}

	ts.SetLocalcoordChar([2]int32{320, 240})
	if ts.localScale != 2 || ts.offsetX != 0 {
		t.Fatalf("SetLocalcoordChar() = scale %v offset %d", ts.localScale, ts.offsetX)
	}

	ts.SetPos(10, 20)
	if ts.offsetInit != [2]float32{10, 20} || ts.x != 20 || ts.y != 40 {
		t.Fatalf("SetPos() = %#v", ts)
	}

	ts.AddPos(1, -2)
	if ts.x != 22 || ts.y != 36 {
		t.Fatalf("AddPos() = x %v y %v", ts.x, ts.y)
	}

	ts.SetScale(1.5, 0.5)
	if ts.scaleInit != [2]float32{1.5, 0.5} || ts.xscl != 3 || ts.yscl != 1 {
		t.Fatalf("SetScale() = %#v", ts)
	}

	ts.SetTextSpacing(4, 5)
	if ts.textSpacing != [2]float32{4, 5} {
		t.Fatalf("SetTextSpacing() = %#v", ts.textSpacing)
	}

	ts.SetVelocity(2, 3)
	if ts.velocityInit != [2]float32{2, 3} || ts.xvel != 4 || ts.yvel != 6 || ts.vel != [2]float32{} {
		t.Fatalf("SetVelocity() = %#v", ts)
	}

	ts.SetMaxDist(7, 8)
	if ts.maxDist != [2]float32{14, 16} {
		t.Fatalf("SetMaxDist() = %#v", ts.maxDist)
	}

	ts.SetAccel(9, 10)
	if ts.accel != [2]float32{18, 20} {
		t.Fatalf("SetAccel() = %#v", ts.accel)
	}

	ts.SetWindow([4]float32{10, 20, 110, 220})
	if ts.windowInit != [4]float32{10, 20, 110, 220} {
		t.Fatalf("SetWindow() windowInit = %#v", ts.windowInit)
	}
	if ts.window[2] != 200 || ts.window[3] != 1200 {
		t.Fatalf("SetWindow() window = %#v", ts.window)
	}
}

func TestTextSpriteSetWindow_UsesTruetypeYBranchAndIgnoresEmptyWindow(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.gameWidth = 640
	sys.gameHeight = 480
	sys.widthScale = 1
	sys.heightScale = 1

	ts := NewTextSprite()
	ts.fnt = &Fnt{Type: "truetype"}
	ts.SetWindow([4]float32{})
	if ts.window != sys.scrrect {
		t.Fatalf("empty SetWindow should keep default window, got %#v", ts.window)
	}

	ts.SetWindow([4]float32{10, 20, 30, 40})
	if ts.window[1] != 20 {
		t.Fatalf("truetype SetWindow should use direct Y, got %#v", ts.window)
	}
}
