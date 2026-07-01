package main

import "testing"

func TestResetRemapInputRestoresIdentityMapping(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	for i := range sys.inputRemap {
		sys.inputRemap[i] = len(sys.inputRemap) - i
	}
	sys.resetRemapInput()
	for i, v := range sys.inputRemap {
		if v != i {
			t.Fatalf("inputRemap[%d] = %d, want %d", i, v, i)
		}
	}
}

func TestLoaderResetRestoresRoundAndPreloadState(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.round = 9
	sys.wins = [2]int32{1, 2}
	sys.roundsExisted = [2]int32{3, 4}
	sys.decisiveRound = [2]bool{true, true}
	sys.turnsPreloadMember = [2]int{2, 3}
	sys.loader.state = LS_NotYet
	sys.loaderReset()

	if sys.round != 1 || sys.wins != [2]int32{} || sys.roundsExisted != [2]int32{} || sys.decisiveRound != [2]bool{} {
		t.Fatalf("loaderReset() did not restore match state: round=%d wins=%#v roundsExisted=%#v decisiveRound=%#v", sys.round, sys.wins, sys.roundsExisted, sys.decisiveRound)
	}
	if sys.turnsPreloadMember != [2]int{-1, -1} {
		t.Fatalf("loaderReset() turnsPreloadMember = %#v, want [-1 -1]", sys.turnsPreloadMember)
	}
	if sys.loader.state != LS_NotYet || sys.loader.err != nil || sys.loader.cancelCh != nil {
		t.Fatalf("loaderReset() should leave loader reset, got %#v", sys.loader)
	}
}
