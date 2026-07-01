package main

import "testing"

func TestTimeHelpersUseMatchTimerAndRoundState(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys

	sys.slowtime = 8
	sys.intro = -1
	sys.curRoundTime = 5
	if got := sys.getSlowtime(); got != 8 {
		t.Fatalf("getSlowtime() = %v, want 8", got)
	}

	sys.curRoundTime = 0
	if got := sys.getSlowtime(); got != 0 {
		t.Fatalf("getSlowtime() with zero round time = %v, want 0", got)
	}

	sys.maxRoundTime = 99
	sys.curRoundTime = 37
	if got := sys.timeElapsed(); got != 62 {
		t.Fatalf("timeElapsed() for timed round = %v, want 62", got)
	}

	sys.maxRoundTime = 0
	sys.curPlayTime = 44
	if got := sys.timeElapsed(); got != 44 {
		t.Fatalf("timeElapsed() for unlimited round = %v, want 44", got)
	}

	sys.curRoundTime = 12
	if got := sys.timeRemaining(); got != 12 {
		t.Fatalf("timeRemaining() = %v, want 12", got)
	}

	sys.curRoundTime = -1
	if got := sys.timeRemaining(); got != -1 {
		t.Fatalf("timeRemaining() after timer expiry = %v, want -1", got)
	}

	sys.timerStart = 10
	sys.timerRounds = []int32{2, 3}
	sys.fightScreen.round.timerActive = true
	sys.maxRoundTime = 99
	sys.curRoundTime = 37
	if got := sys.timeTotal(); got != 77 {
		t.Fatalf("timeTotal() with active timer = %v, want 77", got)
	}

	sys.fightScreen.round.timerActive = false
	if got := sys.timeTotal(); got != 15 {
		t.Fatalf("timeTotal() without active timer = %v, want 15", got)
	}
}
