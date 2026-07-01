package main

import "testing"

func TestShouldPersistMotifAspect_RequiresKeepAspectAndSkippedMotifScaling(t *testing.T) {
	s := &System{}

	if s.shouldPersistMotifAspect() {
		t.Fatalf("expected default motif aspect persistence to be disabled")
	}

	s.cfg.Video.KeepAspect = true
	if s.shouldPersistMotifAspect() {
		t.Fatalf("expected keep-aspect without skip to remain disabled")
	}

	s.matchTime = 1
	s.motif.Info.Localcoord = [2]int32{320, 240}
	s.stage = &Stage{stageCamera: stageCamera{localcoord: [2]int32{640, 240}}}
	if !s.skipMotifScaling() {
		t.Fatalf("expected wide stage to trigger motif scaling skip")
	}
	if !s.shouldPersistMotifAspect() {
		t.Fatalf("expected keep-aspect with skipped motif scaling to persist")
	}
}

func TestEnterAndLeaveMotifAspect_NoOpWhenPersistenceDisabled(t *testing.T) {
	s := &System{}
	s.scrrect = [4]int32{0, 0, 640, 480}
	s.gameWidth = 320
	s.gameHeight = 240
	s.widthScale = 1
	s.heightScale = 1

	s.enterMotifAspect()
	if s.gameWidth != 320 || s.gameHeight != 240 {
		t.Fatalf("expected enterMotifAspect to be a no-op when disabled")
	}

	s.leaveMotifAspect()
	if s.gameWidth != 320 || s.gameHeight != 240 {
		t.Fatalf("expected leaveMotifAspect to be a no-op when disabled")
	}
}
