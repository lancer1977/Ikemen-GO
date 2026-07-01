package main

import "testing"

func TestMiddleOfMatchDependsOnFightLoopMatchTimeAndPostMatch(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys.fightLoopEnd = false
	sys.matchTime = 1
	sys.postMatchFlg = false
	if !sys.middleOfMatch() {
		t.Fatal("middleOfMatch should be true when the fight is active")
	}

	sys.fightLoopEnd = true
	if sys.middleOfMatch() {
		t.Fatal("middleOfMatch should be false after fightLoopEnd")
	}

	sys.fightLoopEnd = false
	sys.matchTime = 0
	if sys.middleOfMatch() {
		t.Fatal("middleOfMatch should be false when matchTime is zero")
	}

	sys.matchTime = 1
	sys.postMatchFlg = true
	if sys.middleOfMatch() {
		t.Fatal("middleOfMatch should be false during post-match")
	}
}
