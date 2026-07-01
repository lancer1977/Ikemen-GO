package main

import "testing"

func TestRoundStateTransitionsAcrossMatchPhases(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.fightScreen.round.ctrl_time = 10
	sys.fightScreen.round.over_waittime = 4

	sys.intro = 12
	if got := sys.roundState(); got != 0 {
		t.Fatalf("roundState() before intro = %v, want 0", got)
	}

	sys.postMatchFlg = false
	sys.intro = 5
	if got := sys.roundState(); got != 1 {
		t.Fatalf("roundState() during intro = %v, want 1", got)
	}

	sys.intro = 0
	sys.finishType = FT_NotYet
	if got := sys.roundState(); got != 2 {
		t.Fatalf("roundState() at fight call = %v, want 2", got)
	}

	sys.finishType = FT_TO
	sys.intro = -2
	if got := sys.roundState(); got != 3 {
		t.Fatalf("roundState() during outro = %v, want 3", got)
	}

	sys.intro = -5
	if got := sys.roundState(); got != 4 {
		t.Fatalf("roundState() after round end = %v, want 4", got)
	}
}
