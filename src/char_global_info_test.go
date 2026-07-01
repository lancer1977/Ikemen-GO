package main

import "testing"

func TestNewCharGlobalInfo(t *testing.T) {
	prev := sys.cfg.Config.PaletteMax
	sys.cfg.Config.PaletteMax = 3
	defer func() { sys.cfg.Config.PaletteMax = prev }()

	gi := newCharGlobalInfo()
	if gi.localcoord != [2]int32{320, 240} {
		t.Fatalf("newCharGlobalInfo localcoord = %#v", gi.localcoord)
	}
	if gi.portraitscale != 1 {
		t.Fatalf("newCharGlobalInfo portraitscale = %v", gi.portraitscale)
	}
	if len(gi.palInfo) != 3 {
		t.Fatalf("newCharGlobalInfo palInfo len = %d, want 3", len(gi.palInfo))
	}
	for i, pal := range gi.palInfo {
		if pal.keyMap != int32(i) {
			t.Fatalf("newCharGlobalInfo palInfo[%d] = %#v", i, pal)
		}
	}
	if gi.constants == nil || gi.states == nil || gi.callFuncs == nil || gi.animTable == nil || gi.fnt == nil {
		t.Fatalf("newCharGlobalInfo should initialize maps/tables: %#v", gi)
	}
}
