package main

import "testing"

func TestCommandStepEqualSteps(t *testing.T) {
	a := CommandStep{greater: true, keys: []CommandStepKey{{key: CK_U}, {key: CK_a}}}
	b := CommandStep{greater: true, keys: []CommandStepKey{{key: CK_U}, {key: CK_a}}}
	if !a.EqualSteps(b) {
		t.Fatal("expected identical command steps to compare equal")
	}

	if a.EqualSteps(CommandStep{greater: false, keys: []CommandStepKey{{key: CK_U}, {key: CK_a}}}) {
		t.Fatal("expected greater flag mismatch to compare unequal")
	}
	if a.EqualSteps(CommandStep{greater: true, keys: []CommandStepKey{{key: CK_a}, {key: CK_U}}}) {
		t.Fatal("expected key order mismatch to compare unequal")
	}
	if a.EqualSteps(CommandStep{greater: true, keys: []CommandStepKey{{key: CK_U}}}) {
		t.Fatal("expected length mismatch to compare unequal")
	}
}
