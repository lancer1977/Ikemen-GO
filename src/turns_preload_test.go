package main

import "testing"

func TestTurnsPreloadActiveReflectsConfigAndSelectedMembers(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.cfg.Config.TurnsLoading = false
	sys.turnsPreloadMember = [2]int{1, -1}
	if sys.turnsPreloadActive() {
		t.Fatal("turnsPreloadActive() should be false when Turns loading is disabled")
	}

	sys.cfg.Config.TurnsLoading = true
	if !sys.turnsPreloadActive() {
		t.Fatal("turnsPreloadActive() should be true when a preload member is selected")
	}
	sys.turnsPreloadMember = [2]int{-1, -1}
	if sys.turnsPreloadActive() {
		t.Fatal("turnsPreloadActive() should be false when no preload member is selected")
	}
}

func TestStartNextTurnsPreloadSelectsNextMemberOrSkipsLoadedOne(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.cfg.Config.TurnsLoading = true
	sys.loader.state = LS_NotYet
	sys.sel.selected[0] = [][2]int{{0, 0}, {1, 0}, {2, 0}}
	sys.sel.selected[1] = [][2]int{{0, 0}, {1, 0}, {2, 0}}
	sys.numTurns = [2]int32{3, 3}
	sys.wins = [2]int32{0, 0}
	sys.tmode = [2]TeamMode{TM_Turns, TM_Turns}
	sys.turnsPreloadMember = [2]int{-1, -1}
	sys.chars[2] = nil
	sys.chars[3] = nil

	sys.startNextTurnsPreload()
	if sys.turnsPreloadMember != [2]int{1, 1} {
		t.Fatalf("startNextTurnsPreload() selected members = %#v, want [1 1]", sys.turnsPreloadMember)
	}
	if sys.loader.state != LS_NotYet {
		t.Fatalf("startNextTurnsPreload() should leave loader idle when already NotYet, got %v", sys.loader.state)
	}

	sys.turnsPreloadMember = [2]int{-1, -1}
	sys.chars[2] = []*Char{{memberNo: 1, selectNo: 1}}
	sys.chars[3] = []*Char{{memberNo: 1, selectNo: 1}}
	sys.startNextTurnsPreload()
	if sys.turnsPreloadMember != [2]int{-1, -1} {
		t.Fatalf("startNextTurnsPreload() should skip already-loaded preload members, got %#v", sys.turnsPreloadMember)
	}

	sys.cfg.Config.TurnsLoading = false
	sys.turnsPreloadMember = [2]int{-1, -1}
	sys.startNextTurnsPreload()
	if sys.turnsPreloadMember != [2]int{-1, -1} {
		t.Fatalf("startNextTurnsPreload() should be a no-op when loading is disabled, got %#v", sys.turnsPreloadMember)
	}
}
