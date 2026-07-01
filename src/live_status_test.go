package main

import "testing"

func TestLiveStatusFighterForSide_FallsBackForBlankNamesAndMissingSlots(t *testing.T) {
	s := &System{}

	if got := s.liveStatusFighterForSide(-1); got != nil {
		t.Fatalf("expected negative side to return nil, got %#v", got)
	}
	if got := s.liveStatusFighterForSide(0); got != nil {
		t.Fatalf("expected empty slot to return nil, got %#v", got)
	}

	s.chars[0] = []*Char{{name: ""}}
	got := s.liveStatusFighterForSide(0)
	if got == nil {
		t.Fatalf("expected populated slot to return fighter data")
	}
	if got.Key != "player-1" || got.Name != "player-1" || got.DisplayName != "player-1" {
		t.Fatalf("unexpected fallback fighter data: %#v", got)
	}
}

func TestLiveStatusStageName_UsesDisplayNameThenNameThenDef(t *testing.T) {
	s := &System{}

	if got := s.liveStatusStageName(); got != "" {
		t.Fatalf("expected nil stage to return blank name, got %q", got)
	}

	s.stage = &Stage{displayname: "Training Ground", name: "training", def: "stages/training.def"}
	if got := s.liveStatusStageName(); got != "Training Ground" {
		t.Fatalf("expected display name to win, got %q", got)
	}

	s.stage.displayname = ""
	if got := s.liveStatusStageName(); got != "training" {
		t.Fatalf("expected stage name fallback, got %q", got)
	}

	s.stage.name = ""
	if got := s.liveStatusStageName(); got != "stages/training.def" {
		t.Fatalf("expected stage def fallback, got %q", got)
	}
}

func TestBuildLiveStatusSnapshot_UsesResultModeAfterMatchOver(t *testing.T) {
	s := &System{
		SystemStateVars: SystemStateVars{
			match: 12,
			round: 2,
		},
		stage: &Stage{displayname: "Training Ground"},
	}
	s.chars[0] = []*Char{{name: "Ryu"}}
	s.chars[1] = []*Char{{name: "Ken"}}
	s.wins = [2]int32{1, 0}
	s.matchWins = [2]int32{1, 1}

	snap := s.buildLiveStatusSnapshot()
	if snap.Mode != "result" {
		t.Fatalf("expected result mode after match over, got %q", snap.Mode)
	}
	if snap.Stage != "Training Ground" {
		t.Fatalf("unexpected stage in snapshot: %q", snap.Stage)
	}
	if snap.P1 == nil || snap.P1.Key != "ryu" || snap.P2 == nil || snap.P2.Key != "ken" {
		t.Fatalf("unexpected fighters in snapshot: %#v", snap)
	}
}
