package main

import "testing"

func TestTimeOverFreezeAndRoundNoDamageFollowMatchState(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	// roundNoDamage() reads sys.fightScreen.round, which is nil on a System that
	// has not entered a fight.
	ensureGlobalRound(t)

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

	// roundNoDamage returns true when intro is within the no-damage window,
	// defined as: over_hittime <= -intro <= over_waittime
	// Set over_hittime=20 (20 frames of hit intro) and over_waittime=60 (total intro time).
	// Then set intro=-20 to be at the boundary of the no-damage window (still no damage).
	sys.fightScreen.round.over_hittime = 20
	sys.fightScreen.round.over_waittime = 60
	sys.intro = -20
	if !sys.roundNoDamage() {
		t.Fatal("roundNoDamage should be true during the no-damage intro window")
	}
}
