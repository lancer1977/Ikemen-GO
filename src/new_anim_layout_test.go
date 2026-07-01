package main

import "testing"

func TestNewAnimLayout_InitializesNestedDefaults(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.scrrect = [4]int32{0, 0, 640, 480}

	sff := &Sff{}
	al := newAnimLayout(sff, 5)
	if al.anim == nil || al.palfx == nil {
		t.Fatal("newAnimLayout should initialize animation and palfx")
	}
	if al.lay.layerno != 5 || al.lay.facing != 1 || al.lay.vfacing != 1 {
		t.Fatalf("newAnimLayout layout defaults = %#v", al.lay)
	}
	if al.lay.scale != [2]float32{1, 1} || al.lay.projection != Projection_Orthographic {
		t.Fatalf("newAnimLayout layout scale/projection = %#v", al.lay)
	}
	if al.anim.window != sys.scrrect {
		t.Fatalf("newAnimLayout animation window = %#v, want %#v", al.anim.window, sys.scrrect)
	}
}
