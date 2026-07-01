package main

import "testing"

func TestTickFrameHelpersRespectPauseAndInterpolationState(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys

	sys.paused = false
	sys.frameStepFlag = false
	sys.oldTickCount = 0
	sys.tickCount = 1
	if !sys.tickFrame() {
		t.Fatal("tickFrame() should advance when not paused")
	}
	if !sys.tickNextFrame() {
		t.Fatal("tickNextFrame() should advance when not paused")
	}

	sys.paused = true
	sys.frameStepFlag = false
	sys.oldTickCount = 0
	sys.tickCount = 1
	if sys.tickFrame() {
		t.Fatal("tickFrame() should stop while paused")
	}
	if sys.tickNextFrame() {
		t.Fatal("tickNextFrame() should stop while paused")
	}

	sys.frameStepFlag = true
	if !sys.tickFrame() || !sys.tickNextFrame() {
		t.Fatal("tickFrame() and tickNextFrame() should advance during frame step")
	}
}

func TestPerspectiveHelpersRespectZStateAndProjection(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.zmin = 0
	sys.zmax = 1
	if !sys.zEnabled() {
		t.Fatal("zEnabled() should report enabled for distinct Z bounds")
	}
	sys.zmax = 0
	if sys.zEnabled() {
		t.Fatal("zEnabled() should report disabled when Z bounds match")
	}

	sys.cam.Pos[0] = 10
	sys.stage = &Stage{}
	sys.stage.stageCamera.depthtoscreen = 0.5
	if got := sys.posZtoYoffset(4, 2); got != 4 {
		t.Fatalf("posZtoYoffset() = %v, want 4", got)
	}

	out := sys.drawposXYfromZ([2]float32{14, 3}, 2, 4, 0.5)
	if out[0] != 12 || out[1] != 5 {
		t.Fatalf("drawposXYfromZ() = %#v, want [12 5]", out)
	}
}
