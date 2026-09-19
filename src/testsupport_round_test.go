package main

import "testing"

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
