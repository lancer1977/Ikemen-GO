package main

import "testing"

func TestGetKeyboardState(t *testing.T) {
	t.Parallel()

	prevSys := sys
	sys = System{}
	t.Cleanup(func() { sys = prevSys })

	sys.keyState[Key(1)] = true
	sys.keyState[Key(4)] = true
	sys.keyState[Key(9)] = true

	kc := KeyConfig{Joy: -1, dU: 1, dD: 2, dL: 3, dR: 4, bA: 5, bB: 6, bC: 7, bX: 8, bY: 9, bZ: 10, bS: 11, bD: 12, bW: 13, bM: 14}
	got := GetKeyboardState(kc)

	if !got[0] || !got[3] || !got[8] {
		t.Fatalf("unexpected keyboard state: %#v", got)
	}
	if got[1] || got[2] || got[4] || got[5] || got[6] || got[7] || got[9] || got[10] || got[11] || got[12] || got[13] {
		t.Fatalf("expected only seeded keys to be true: %#v", got)
	}
}
