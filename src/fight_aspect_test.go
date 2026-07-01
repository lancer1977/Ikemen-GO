package main

import "testing"

func TestGetFightAspect_UsesStageCustomOrMotifFallback(t *testing.T) {
	s := &System{}
	s.scrrect = [4]int32{0, 0, 640, 480}

	if got := s.getFightAspect(); got != CalculateAspect(640, 480) {
		t.Fatalf("expected motif aspect fallback, got %v", got)
	}

	s.cfg.Video.FightAspectWidth = 16
	s.cfg.Video.FightAspectHeight = 9
	if got := s.getFightAspect(); got != CalculateAspect(16, 9) {
		t.Fatalf("expected custom fight aspect, got %v", got)
	}

	s.cfg.Video.FightAspectWidth = -1
	s.cfg.Video.FightAspectHeight = -1
	s.stage = &Stage{stageCamera: stageCamera{localcoord: [2]int32{320, 240}}}
	if got := s.getFightAspect(); got != CalculateAspect(320, 240) {
		t.Fatalf("expected stage fight aspect, got %v", got)
	}

	s.stage.stageCamera.localcoord = [2]int32{0, 0}
	if got := s.getFightAspect(); got != CalculateAspect(640, 480) {
		t.Fatalf("expected invalid stage localcoord to fall back to motif aspect, got %v", got)
	}
}
