package main

import "testing"

func TestRectMutators_UpdatePackedColorAlphaAndWindowState(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.gameWidth = 640
	sys.gameHeight = 480
	sys.widthScale = 2
	sys.heightScale = 3

	r := &Rect{}
	r.SetColor([3]int32{1, 2, 3})
	if r.col != 0x030201 {
		t.Fatalf("SetColor() = %#x, want 0x030201", r.col)
	}

	r.SetAlpha([2]int32{10, 20})
	if r.autoAlpha {
		t.Fatal("SetAlpha should disable autoAlpha")
	}

	r.SetAlphaPulse(128, 64, 120)
	if !r.autoAlpha || r.pulseMid != 128 || r.pulseAmp != 64 || r.pulsePhaseStep == 0 {
		t.Fatalf("SetAlphaPulse() = %#v", r)
	}

	r.SetLocalcoord(640, 480)
	if r.localScale != 2 || r.offsetX != 0 {
		t.Fatalf("SetLocalcoord() = scale %v offset %d", r.localScale, r.offsetX)
	}

	r.SetWindow([4]float32{10, 20, 110, 220})
	if r.windowInit != [4]float32{10, 20, 110, 220} {
		t.Fatalf("SetWindow() windowInit = %#v", r.windowInit)
	}
	if r.window[2] != 100 || r.window[3] != 200 {
		t.Fatalf("SetWindow() size = %#v", r.window)
	}
}

func TestRectMutators_IgnoreInvalidInputs(t *testing.T) {
	r := &Rect{window: [4]int32{1, 2, 3, 4}}
	r.SetAlphaPulse(12, 0, 0)
	if r.pulseAmp != 0 || r.pulsePhaseStep != 0 || r.autoAlpha {
		t.Fatalf("SetAlphaPulse invalid input = %#v", r)
	}

	r.SetLocalcoord(0, 480)
	if r.localScale != 0 {
		t.Fatalf("SetLocalcoord invalid input changed scale: %v", r.localScale)
	}

	r.SetWindow([4]float32{})
	if r.window != [4]int32{1, 2, 3, 4} {
		t.Fatalf("SetWindow zero input should no-op: %#v", r.window)
	}
}

func TestCameraBounds_HonorEnableDebugAndClamp(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.turbo = 0.5
	sys.paused = false
	sys.frameStepFlag = false
	sys.oldTickCount = 0
	sys.tickCount = 1

	c := &Camera{
		ZoomEnable: true,
		MinScale:   1.5,
		zoomin:     3,
	}
	if got := c.ScaleBound(2, 4); got != 2 {
		t.Fatalf("ScaleBound enabled = %v, want 2", got)
	}

	sys.paused = true
	sys.frameStepFlag = false
	sys.oldTickCount = 0
	sys.tickCount = 1
	if got := c.ScaleBound(2, 4); got != 2 {
		t.Fatalf("ScaleBound debug paused = %v, want 2", got)
	}

	c.ZoomEnable = false
	if got := c.ScaleBound(2, 4); got != 1 {
		t.Fatalf("ScaleBound disabled = %v, want 1", got)
	}

	c.ZoomEnable = true
	c.boundL = -100
	c.boundR = 100
	c.halfWidth = 50
	if got := c.XBound(2, 999); got != 125 {
		t.Fatalf("XBound upper clamp = %v, want 125", got)
	}
	if got := c.XBound(2, -999); got != -125 {
		t.Fatalf("XBound lower clamp = %v, want -125", got)
	}
}

func TestCameraBaseScaleAndGroundLevel(t *testing.T) {
	c := &Camera{
		ztopscale:            1.75,
		zoff:                 240,
		aspectcorrection:     12,
		zoomanchorcorrection: 8,
	}
	if got := c.BaseScale(); got != 1.75 {
		t.Fatalf("BaseScale() = %v, want 1.75", got)
	}
	if got := c.GroundLevel(); got != 220 {
		t.Fatalf("GroundLevel() = %v, want 220", got)
	}
}
