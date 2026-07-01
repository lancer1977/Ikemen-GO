package main

import "testing"

func TestTimeOverFreezeAndRoundNoDamageFollowMatchState(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.curRoundTime = 0
	sys.finishType = FT_TO
	sys.winposetime = 1
	if !sys.timeOverFreeze() {
		t.Fatal("timeOverFreeze should be true for a finished time-over round")
	}

	sys.winposetime = 0
	if sys.timeOverFreeze() {
		t.Fatal("timeOverFreeze should be false when win pose time is exhausted")
	}

	sys.intro = -sys.fightScreen.round.over_hittime
	sys.fightScreen.round.over_waittime = sys.fightScreen.round.over_hittime
	if !sys.roundNoDamage() {
		t.Fatal("roundNoDamage should be true during the no-damage intro window")
	}
}
