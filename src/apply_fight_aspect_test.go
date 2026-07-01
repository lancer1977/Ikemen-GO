package main

import "testing"

func TestApplyFightAspect_UsesSelectedStageOrMotifFallback(t *testing.T) {
	s := &System{}
	s.scrrect = [4]int32{0, 0, 640, 480}
	s.sel.selectedStageNo = 1
	s.sel.stagelist = []SelectStage{{def: "stages/kaiser.def"}}
	s.stageLocalcoords = map[string][2]int32{
		"kaiser.def": {400, 300},
	}

	s.applyFightAspect()
	if s.gameWidth != 400 || s.gameHeight != 240 {
		t.Fatalf("expected stage aspect to drive game size, got width=%v height=%v", s.gameWidth, s.gameHeight)
	}

	s.stageLocalcoords = map[string][2]int32{}
	s.applyFightAspect()
	if s.gameWidth != 320 || s.gameHeight != 240 {
		t.Fatalf("expected motif fallback to drive game size, got width=%v height=%v", s.gameWidth, s.gameHeight)
	}
}
