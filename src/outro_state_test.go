package main

import "testing"

func TestOutroStateCoversRoundEndAndWinPhaseTransitions(t *testing.T) {
	// outroState accesses sys.fightScreen.round, which is nil on a zero-value System
	ensureGlobalRound(t)

	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.fightScreen.round.over_waittime = 4
	sys.fightScreen.round.over_hittime = 2

	sys.intro = 1
	if got := sys.outroState(); got != 0 {
		t.Fatalf("outroState() before outro = %v, want 0", got)
	}

	sys.intro = -6
	if got := sys.outroState(); got != 5 {
		t.Fatalf("outroState() after round over = %v, want 5", got)
	}

	sys.intro = -1
	sys.winposetime = 0
	if got := sys.outroState(); got != 4 {
		t.Fatalf("outroState() in win states = %v, want 4", got)
	}

	sys.winposetime = 1
	sys.intro = -4
	if got := sys.outroState(); got != 3 {
		t.Fatalf("outroState() in pre-win control loss = %v, want 3", got)
	}

	sys.intro = -2
	// With over_hittime=2, the condition `intro < -over_hittime` checks if intro < -2
	// At intro=-2, that's false (not more negative than -2), so we fall through to `intro < 0` case
	// which returns 1 (players can still act and potentially change outcome with double KO)
	if got := sys.outroState(); got != 1 {
		t.Fatalf("outroState() in late control = %v, want 1", got)
	}

	sys.fightScreen.round.over_hittime = 3
	if got := sys.outroState(); got != 1 {
		t.Fatalf("outroState() in double-KO window = %v, want 1", got)
	}
}
