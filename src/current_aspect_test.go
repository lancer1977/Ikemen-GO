package main

import "testing"

func TestGetCurrentAspect_UsesMotifAspectOutsideMatch(t *testing.T) {
	s := &System{}
	s.scrrect = [4]int32{0, 0, 640, 480}
	if got := s.getCurrentAspect(); got != CalculateAspect(640, 480) {
		t.Fatalf("expected motif aspect outside match, got %v", got)
	}
}

func TestGetCurrentAspect_UsesFightAspectDuringMatchWhenMotifScalingIsSkipped(t *testing.T) {
	s := &System{}
	s.scrrect = [4]int32{0, 0, 640, 480}
	s.cfg.Video.FightAspectWidth = 16
	s.cfg.Video.FightAspectHeight = 9
	s.matchTime = 1
	s.motif.me.active = false
	s.motif.di.active = false
	if got := s.getCurrentAspect(); got != CalculateAspect(16, 9) {
		t.Fatalf("expected fight aspect during match when motif scaling is skipped, got %v", got)
	}
}

func TestGetCurrentAspect_UsesFightAspectAfterMatchWhenSkipEnabled(t *testing.T) {
	s := &System{}
	s.scrrect = [4]int32{0, 0, 640, 480}
	s.cfg.Video.FightAspectWidth = 16
	s.cfg.Video.FightAspectHeight = 9
	s.postMatchFlg = true
	s.cfg.Video.KeepAspect = true
	if got := s.getCurrentAspect(); got != CalculateAspect(16, 9) {
		t.Fatalf("expected fight aspect post-match when skip is enabled, got %v", got)
	}
}
