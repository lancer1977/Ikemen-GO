package main

import "testing"

func TestShouldPersistMotifAspect_RequiresKeepAspectAndSkippedMotifScaling(t *testing.T) {
	s := newTestSystem()
	// getMotifAspect() reads from scrrect, so initialize it
	s.scrrect = [4]int32{0, 0, 320, 240}

	if s.shouldPersistMotifAspect() {
		t.Fatalf("expected default motif aspect persistence to be disabled")
	}

	s.cfg.Video.KeepAspect = true
	// Production code: shouldPersistMotifAspect = KeepAspect && !skipMotifScaling
	// With default state (no match time, stage==nil): skipMotifScaling() == false
	// So shouldPersistMotifAspect = true && !false = true
	if !s.shouldPersistMotifAspect() {
		t.Fatalf("expected keep-aspect with no skip to enable persistence")
	}

	s.matchTime = 1
	s.motif.Info.Localcoord = [2]int32{320, 240}
	s.stage = &Stage{stageCamera: stageCamera{localcoord: [2]int32{640, 240}}}
	if !s.skipMotifScaling() {
		t.Fatalf("expected wide stage to trigger motif scaling skip")
	}
	// Production code: shouldPersistMotifAspect = KeepAspect && !skipMotifScaling
	// With wide stage: skipMotifScaling() == true
	// So shouldPersistMotifAspect = true && !true = false (persistence is disabled when skip is active)
	if s.shouldPersistMotifAspect() {
		t.Fatalf("expected keep-aspect with skipped motif scaling to disable persistence")
	}
}

func TestEnterAndLeaveMotifAspect_NoOpWhenPersistenceDisabled(t *testing.T) {
	s := newTestSystem()
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
