package main

import "testing"

func TestCaptureAndRestoreAspectState_RoundTripsDimensions(t *testing.T) {
	s := &System{}
	s.gameWidth = 320
	s.gameHeight = 240
	s.widthScale = 1.25
	s.heightScale = 1.5

	state := s.captureAspectState()
	s.gameWidth = 640
	s.gameHeight = 480
	s.widthScale = 2
	s.heightScale = 3
	s.restoreAspectState(state)

	if s.gameWidth != 320 || s.gameHeight != 240 || s.widthScale != 1.25 || s.heightScale != 1.5 {
		t.Fatalf("unexpected restored aspect state: %#v", s.captureAspectState())
	}
}

func TestWrapDrawWithAspectState_RestoresOuterStateAroundCallback(t *testing.T) {
	s := &System{}
	s.gameWidth = 320
	s.gameHeight = 240
	s.widthScale = 1
	s.heightScale = 1

	wrapped := s.wrapDrawWithAspectState(func() {
		if s.gameWidth != 320 || s.gameHeight != 240 {
			t.Fatalf("expected wrapped callback to see captured aspect state, got width=%v height=%v", s.gameWidth, s.gameHeight)
		}
		s.gameWidth = 800
		s.gameHeight = 600
		s.widthScale = 2
		s.heightScale = 2.5
	})
	if wrapped == nil {
		t.Fatalf("expected wrapped function")
	}

	s.gameWidth = 640
	s.gameHeight = 480
	s.widthScale = 4
	s.heightScale = 5
	wrapped()

	if s.gameWidth != 640 || s.gameHeight != 480 || s.widthScale != 4 || s.heightScale != 5 {
		t.Fatalf("expected outer aspect state to be restored, got %#v", s.captureAspectState())
	}
}

func TestWrapDrawWithAspectState_ReturnsNilForNilCallback(t *testing.T) {
	s := &System{}
	if got := s.wrapDrawWithAspectState(nil); got != nil {
		t.Fatalf("expected nil callback to return nil wrapper, got %#v", got)
	}
}

func TestLuaDrawQueueRoutesAndFlushesCallbacks(t *testing.T) {
	s := &System{}
	s.gameWidth = 320
	s.gameHeight = 240
	s.widthScale = 1
	s.heightScale = 1

	var calls []string
	s.luaQueuePreDraw(func() { calls = append(calls, "pre") })
	s.luaQueueLayerDraw(-1, func() { calls = append(calls, "layer-neg") })
	s.luaQueueLayerDraw(1, func() { calls = append(calls, "layer-1") })
	s.luaQueueLayerDraw(9, func() { calls = append(calls, "layer-clamped") })

	if len(s.luaDrawPreOps) != 2 {
		t.Fatalf("expected two pre ops, got %d", len(s.luaDrawPreOps))
	}
	if len(s.luaDrawLayerOps[1]) != 1 || len(s.luaDrawLayerOps[len(s.luaDrawLayerOps)-1]) != 1 {
		t.Fatalf("unexpected layer queue state: %#v", s.luaDrawLayerOps)
	}

	s.luaFlushDrawQueue()
	if got := calls; len(got) != 4 || got[0] != "pre" || got[1] != "layer-neg" || got[2] != "layer-1" || got[3] != "layer-clamped" {
		t.Fatalf("unexpected flush order: %#v", got)
	}
	if len(s.luaDrawPreOps) != 0 {
		t.Fatalf("expected pre queue to be cleared after flush, got %d", len(s.luaDrawPreOps))
	}
	if len(s.luaDrawLayerOps[1]) != 0 || len(s.luaDrawLayerOps[len(s.luaDrawLayerOps)-1]) != 0 {
		t.Fatalf("expected layer queues to be cleared after flush, got %#v", s.luaDrawLayerOps)
	}

	s.luaQueuePreDraw(func() { calls = append(calls, "discarded") })
	s.luaQueueLayerDraw(0, func() { calls = append(calls, "discarded-layer") })
	s.luaDiscardDrawQueue()
	if len(s.luaDrawPreOps) != 0 {
		t.Fatalf("expected discard to clear pre queue, got %d", len(s.luaDrawPreOps))
	}
	for i := range s.luaDrawLayerOps {
		if len(s.luaDrawLayerOps[i]) != 0 {
			t.Fatalf("expected discard to clear layer queue %d, got %#v", i, s.luaDrawLayerOps[i])
		}
	}
	if len(calls) != 4 {
		t.Fatalf("discard should not execute queued callbacks, got %#v", calls)
	}
}
