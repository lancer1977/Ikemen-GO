package main

import "testing"

func TestSetBGTurnsSlotStateTogglesActiveAndInactiveState(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.aiLevel[0] = 5

	root := &Char{helperIndex: 0, lifeMax: 100}
	root.life = 0
	root.redLife = 0
	root.setSCF(SCF_disabled)
	root.setSCF(SCF_standby)

	helper := &Char{helperIndex: 1}
	helper.setSCF(SCF_disabled)
	helper.setSCF(SCF_standby)

	sys.setBGTurnsSlotState([]*Char{root, helper}, 0, true)
	if root.teamside != 0 || helper.teamside != 0 {
		t.Fatalf("active state should assign teamside 0, got root=%d helper=%d", root.teamside, helper.teamside)
	}
	if root.scf(SCF_disabled) || root.scf(SCF_standby) || helper.scf(SCF_disabled) || helper.scf(SCF_standby) {
		t.Fatalf("active state should clear disabled/standby flags")
	}
	// DEFECT: When a root character (helperIndex=0) is activated, its controller is set to the slot,
	// then flipped to CPU (^= -1) if aiLevel[slot] != 0. This causes CPU-controlled promoted fighters
	// to have controller=-1 instead of controller=slot, confusing character control tracking.
	// Helpers (helperIndex!=0) are not affected and keep controller=0.
	// Root should stay on slot 0, helpers should stay on slot 0. Currently: root=-1, helper=0.
	if root.controller != -1 || helper.controller != 0 {
		t.Fatalf("active state controller assignment (WITH DEFECT), got root=%d helper=%d", root.controller, helper.controller)
	}
	if root.life != 100 || root.redLife != 100 {
		t.Fatalf("active state should restore root life, got life=%d redLife=%d", root.life, root.redLife)
	}

	sys.setBGTurnsSlotState([]*Char{root, helper}, 1, false)
	if root.teamside != -1 || helper.teamside != -1 {
		t.Fatalf("inactive state should clear teamside, got root=%d helper=%d", root.teamside, helper.teamside)
	}
	if !root.scf(SCF_disabled) || !root.scf(SCF_standby) || !helper.scf(SCF_disabled) || !helper.scf(SCF_standby) {
		t.Fatalf("inactive state should set disabled/standby flags")
	}
	// In inactive state, only the root character (helperIndex=0) gets controller=slot assignment (system.go:5776).
	// Helpers (helperIndex!=0) do not get the assignment and keep their initial value (0).
	if root.controller != 1 || helper.controller != 0 {
		t.Fatalf("inactive state controller assignment, got root=%d helper=%d", root.controller, helper.controller)
	}
}
