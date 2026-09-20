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
	// Belt-and-braces: if the test exits early (e.g. via t.Fatal) before the
	// explicit sys.loader.reset() call below runs, make sure the Loader's
	// background goroutine is still stopped before sys gets swapped back.
	// Registered after the sys restore defer so it runs first (defers are
	// LIFO). See the longer explanation at the reset() call site below.
	defer func() { sys.loader.reset() }()

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
	// When startNextTurnsPreload finds work to do and loader is NotYet, it calls runTread()
	// which immediately transitions to LS_Loading. See system.go:6471.
	if sys.loader.state != LS_Loading {
		t.Fatalf("startNextTurnsPreload() should start loader when work found, got %v", sys.loader.state)
	}
	// runTread() started the Loader's load() goroutine (SafeGo, system.go:6473),
	// which keeps reading/writing sys fields (sys.selMutex, sys.turnsPreloadMember
	// via turnsPreloadActive(), sys.chars, ...) on its own 10ms-poll schedule
	// until it observes a cancel or finishes. The rest of this test mutates
	// those same fields directly, so the goroutine must be stopped here —
	// not only at teardown — or it races/corrupts sys concurrently with the
	// test body itself. Loader.reset() cancels it and blocks on <-l.loadExit
	// until it has actually exited, which previously showed up as a ~1-in-8
	// fatal "RUnlock of unlocked RWMutex" panic when a leaked instance of
	// this goroutine outlived the test and raced a later `sys = *newTestSystem()`
	// reset in some other test (lancer1977/Ikemen-GO#25), and -race also
	// flags data races on the same fields within this test body itself.
	sys.loader.reset()

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
