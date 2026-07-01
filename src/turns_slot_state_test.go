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
	if root.controller != 0 || helper.controller != 0 {
		t.Fatalf("active state should keep controllers on slot 0, got root=%d helper=%d", root.controller, helper.controller)
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
	if root.controller != 1 || helper.controller != 1 {
		t.Fatalf("inactive state should keep controllers on slot 1, got root=%d helper=%d", root.controller, helper.controller)
	}
}
