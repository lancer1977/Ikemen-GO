package main

import "testing"

func TestActivateNextTurnsFightersPromotesPreloadedMemberIntoActiveSlot(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.charList.idMap = make(map[int32]*Char)
	sys.charList.creationOrder = nil
	sys.charList.runOrder = nil
	sys.tmode = [2]TeamMode{TM_Turns, TM_Single}
	sys.effectiveLoss = [2]bool{true, false}
	sys.numTurns = [2]int32{2, 0}
	// activateNextTurnsFighters checks wins[team^1] (= wins[1]) against the preloaded char's memberNo.
	// With memberNo=1 at slot 2, wins[1] must equal 1 for promotion to proceed.
	sys.wins = [2]int32{0, 1}
	sys.turnsPreloadMember = [2]int{1, -1}
	sys.sel.selected[0] = [][2]int{{0, 0}, {1, 7}}
	sys.cfg.Config.TurnsLoading = true
	sys.loader.state = LS_NotYet
	sys.cgi[0].states = map[int32]StateBytecode{10: {playerNo: 2}}
	sys.cgi[2].states = map[int32]StateBytecode{20: {playerNo: 0}}

	active := &Char{id: 11, memberNo: 0, selectNo: 0, helperIndex: 0, controller: 0, life: 10, lifeMax: 10}
	active.ss.sb.playerNo = 0
	active.setSCF(SCF_disabled)
	active.setSCF(SCF_standby)
	preloaded := &Char{id: 22, memberNo: 1, selectNo: 7, helperIndex: 0, controller: 2, life: 0, lifeMax: 100}
	preloaded.ss.sb.playerNo = 2
	preloaded.setSCF(SCF_disabled)
	preloaded.setSCF(SCF_standby)

	sys.chars[0] = []*Char{active}
	sys.chars[2] = []*Char{preloaded}
	sys.charList.add(active)
	sys.charList.add(preloaded)

	sys.activateNextTurnsFighters()

	if sys.chars[0][0] != preloaded || sys.chars[2][0] != active {
		t.Fatalf("activateNextTurnsFighters() did not swap slots: slot0=%#v slot2=%#v", sys.chars[0][0], sys.chars[2][0])
	}
	if preloaded.playerNo != 0 || preloaded.ss.sb.playerNo != 0 || preloaded.controller != 0 {
		t.Fatalf("promoted fighter did not get remapped to active slot: %#v", preloaded)
	}
	if active.playerNo != 2 || active.ss.sb.playerNo != 2 {
		t.Fatalf("demoted fighter did not get remapped to preload slot: %#v", active)
	}
	if preloaded.teamside != 0 || preloaded.scf(SCF_disabled) || preloaded.scf(SCF_standby) {
		t.Fatalf("promoted fighter should be active and enabled, got %#v", preloaded)
	}
	if active.teamside != -1 || !active.scf(SCF_disabled) || !active.scf(SCF_standby) {
		t.Fatalf("demoted fighter should be inactive and disabled, got %#v", active)
	}
	// After swapping chars[0] and chars[2], the CGI structs also swap (system.go:5803).
	// So cgi[0] now contains the old cgi[2].states (which had key 20),
	// and cgi[2] now contains the old cgi[0].states (which had key 10).
	// rebindCgiStateOwners then updates playerNo for all states to match their slot.
	if sys.cgi[0].states[20].playerNo != 0 || sys.cgi[2].states[10].playerNo != 2 {
		t.Fatalf("state owners were not rebound to their slots: cgi0=%#v cgi2=%#v", sys.cgi[0].states[20], sys.cgi[2].states[10])
	}
	if len(sys.charList.creationOrder) == 0 || sys.charList.idMap[preloaded.id] != preloaded {
		t.Fatalf("promoted fighter should be present in charList, got creationOrder=%#v idMap=%#v", sys.charList.creationOrder, sys.charList.idMap)
	}
}
