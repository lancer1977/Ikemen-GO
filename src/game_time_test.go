package main

import "testing"

func TestGameTimePrefersNetworkReplayAndFallbackOffsets(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.matchTime = 120
	sys.preMatchTime = 30

	if got := sys.gameTime(); got != 150 {
		t.Fatalf("gameTime() with fallback offset = %v, want 150", got)
	}

	sys.replayFile = &ReplayFile{preMatchTime: 40}
	if got := sys.gameTime(); got != 160 {
		t.Fatalf("gameTime() with replay offset = %v, want 160", got)
	}

	sys.netConnection = &NetConnection{preMatchTime: 50}
	if got := sys.gameTime(); got != 170 {
		t.Fatalf("gameTime() with network offset = %v, want 170", got)
	}
}
