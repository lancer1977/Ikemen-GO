package main

import "testing"

func TestBuildLiveOverlayRuntime(t *testing.T) {
	s := &System{}
	s.frameCounter = 123

	if got := s.buildLiveOverlayRuntime(LiveOverlayEffect{Kind: "unknown"}); got != nil {
		t.Fatalf("expected unknown kind to return nil, got %#v", got)
	}

	if got := s.buildLiveOverlayRuntime(LiveOverlayEffect{Kind: "text", Text: "hello"}); got != nil {
		t.Fatalf("expected text overlay without debug font to return nil, got %#v", got)
	}
}
