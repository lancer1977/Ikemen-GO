package main

import (
	"os"
	"testing"
)

// TestMain gives the package-level sys a non-nil fightScreen.round before any
// test runs.
//
// src is a single package, so every test shares one global sys. A zero-value
// System leaves fightScreen.round nil, and the System and Char methods that
// dereference it panic rather than fail -- which aborts the whole test binary
// and hides every test ordered after the one that tripped it. Establishing the
// pointer once up front matches the real precondition (these methods are only
// reachable once a fight exists) and keeps a single bad assumption from
// masking the rest of the suite.
//
// Tests that replace sys wholesale still need ensureGlobalRound, since that
// assignment drops this pointer.
func TestMain(m *testing.M) {
	sys.fightScreen.round = &FightScreenRound{}
	sys.sel.gameParams = newGameParams()
	os.Exit(m.Run())
}

// ensureGlobalGameParams gives the package-level sys a non-nil sel.gameParams
// for the duration of a test.
//
// Production treats a nil gameParams as a lazily-initialised state and guards
// every read with a nil check (see char.go and script.go). Tests that assign
// into its fields directly have no such guard, so they need the pointer
// established first. The previous value is restored on cleanup.
func ensureGlobalGameParams(t *testing.T) {
	t.Helper()
	prev := sys.sel.gameParams
	sys.sel.gameParams = newGameParams()
	t.Cleanup(func() {
		sys.sel.gameParams = prev
	})
}

// ensureGlobalRound gives the package-level sys a non-nil fightScreen.round for
// the duration of a test.
//
// Several System methods (roundState, roundOver) dereference
// sys.fightScreen.round, and a number of Char methods reach the global sys
// rather than the receiver's own System. The zero-value sys leaves that pointer
// nil, so any test that reaches one of those paths panics unless it is set up
// here first. The previous value is restored on cleanup.
func ensureGlobalRound(t *testing.T) {
	t.Helper()
	prev := sys.fightScreen.round
	sys.fightScreen.round = &FightScreenRound{}
	t.Cleanup(func() {
		sys.fightScreen.round = prev
	})
}
