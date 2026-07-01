package main

import "testing"

func TestGameParamsConstructorsInitializeDefaults(t *testing.T) {
	gp := newGameParams()
	if gp == nil {
		t.Fatal("newGameParams returned nil")
	}
	if !gp.Continue || gp.Order != -1 || gp.Time != -1 || gp.AI != -1 || !gp.VsScreen || !gp.VictoryScreen || !gp.WinScreen {
		t.Fatalf("unexpected newGameParams defaults: %#v", gp)
	}
	if len(gp.musicEntries) != 0 || len(gp.Raw) != 0 {
		t.Fatalf("newGameParams should start with empty slices: %#v", gp)
	}

	if got := newGameParamsFromMotif(nil); got == nil || got.Continue != gp.Continue || got.VictoryScreen != gp.VictoryScreen || got.WinScreen != gp.WinScreen {
		t.Fatalf("newGameParamsFromMotif(nil) should preserve defaults: %#v", got)
	}

	if got := newGameParamsFromMotif(&Motif{}); got.Continue || got.VictoryScreen || got.WinScreen {
		t.Fatalf("newGameParamsFromMotif should honor zero-value motif flags: %#v", got)
	}
}
