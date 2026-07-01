package main

import "testing"

func TestNewBGDef(t *testing.T) {
	bg := newBGDef("stages/test.def")
	if bg == nil {
		t.Fatal("expected newBGDef to allocate")
	}
	if bg.def != "stages/test.def" {
		t.Fatalf("newBGDef def = %q", bg.def)
	}
	if bg.localcoord != [2]int32{320, 240} {
		t.Fatalf("newBGDef localcoord = %#v", bg.localcoord)
	}
	if !bg.resetbg || bg.localscl != 1 || bg.scale != [2]float32{1, 1} || bg.lastTick != -1 {
		t.Fatalf("newBGDef defaults = %#v", bg)
	}
}
