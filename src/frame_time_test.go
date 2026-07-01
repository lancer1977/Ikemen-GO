package main

import (
	"testing"
	"time"
)

func TestFrameTimeHelpersHandlePauseAdvanceAndReset(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys

	sys.paused = true
	sys.frameStepFlag = false
	sys.oldTickCount = 0
	sys.tickCount = 1
	if !sys.debugPaused() {
		t.Fatal("debugPaused() should report paused frame advancement")
	}

	sys.nextAddTime = 0.5
	if got := sys.addFrameTime(1); !got || sys.oldNextAddTime != 0 {
		t.Fatalf("addFrameTime() paused path = got %v oldNextAddTime=%v", got, sys.oldNextAddTime)
	}

	sys.paused = false
	sys.tickCount = 0
	sys.oldTickCount = 0
	sys.tickCountF = 0
	sys.lastTick = 0
	sys.nextAddTime = 1
	if got := sys.addFrameTime(0.5); !got || sys.tickCount != 1 || sys.lastTick != 1 {
		t.Fatalf("addFrameTime() advance path = tickCount=%d lastTick=%v", sys.tickCount, sys.lastTick)
	}

	sys.resetFrameTime()
	if sys.tickCount != 0 || sys.oldTickCount != -1 || sys.tickCountF != 0 || sys.lastTick != 0 || sys.nextAddTime != 1 || sys.oldNextAddTime != 1 {
		t.Fatalf("resetFrameTime() state = %#v", sys)
	}
	if time.Since(sys.redrawWait.nextTime) < 0 || time.Since(sys.redrawWait.lastDraw) < 0 {
		t.Fatal("resetFrameTime() should refresh redraw timestamps")
	}
}
