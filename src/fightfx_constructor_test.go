package main

import "testing"

func TestNewFightFxInitializesDefaults(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys.fightScreen.localcoord = [2]int32{640, 480}
	ffx := newFightFx()
	if ffx == nil {
		t.Fatal("newFightFx returned nil")
	}
	if ffx.sff == nil {
		t.Fatal("newFightFx should allocate an SFF")
	}
	if ffx.fx_scale != 1.0 {
		t.Fatalf("newFightFx fx_scale = %v, want 1.0", ffx.fx_scale)
	}
	if ffx.localcoord != [2]int32{640, 480} {
		t.Fatalf("newFightFx localcoord = %#v, want %#v", ffx.localcoord, [2]int32{640, 480})
	}
}
