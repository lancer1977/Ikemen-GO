package main

import "testing"

func TestSelectParamsConstructorsInitializeDefaults(t *testing.T) {
	scp := newSelectCharParams()
	if scp == nil {
		t.Fatal("newSelectCharParams returned nil")
	}
	if len(scp.musicEntries) != 0 || len(scp.Raw) != 0 {
		t.Fatalf("newSelectCharParams should start empty, got %#v", scp)
	}
	if scp.AI != -1 || !scp.VsScreen || !scp.VictoryScreen || scp.Rounds != -1 || scp.Time != -1 || scp.IncludeStage != 1 || scp.Hidden != 0 || scp.Order != -1 || scp.OrderSurvival != -1 {
		t.Fatalf("unexpected newSelectCharParams defaults: %#v", scp)
	}

	ssp := newSelectStageParams()
	if ssp == nil {
		t.Fatal("newSelectStageParams returned nil")
	}
	if len(ssp.musicEntries) != 0 || len(ssp.Order) != 0 || len(ssp.Raw) != 0 {
		t.Fatalf("newSelectStageParams should start empty, got %#v", ssp)
	}
}
