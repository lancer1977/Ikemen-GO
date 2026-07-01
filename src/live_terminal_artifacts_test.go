package main

import "testing"

func TestTerminalArtifactKey_ChangesWithMatchState(t *testing.T) {
	s := &System{}
	key1 := s.terminalArtifactKey()

	s.match = 12
	s.round = 2
	s.winTeam = 1
	s.wins = [2]int32{1, 0}
	s.statsLog.Matches = []StatsMatch{{}}
	key2 := s.terminalArtifactKey()

	if key1 == key2 {
		t.Fatalf("expected terminal artifact key to change with match state")
	}
}

func TestTerminalWinnerData_ReturnsBlankWhenNoWinnerAndKeysWhenPresent(t *testing.T) {
	s := &System{}
	if side, winner, loser := s.terminalWinnerData(); side != 0 || winner != "" || loser != "" {
		t.Fatalf("expected blank winner data when winTeam is unset, got side=%d winner=%q loser=%q", side, winner, loser)
	}

	s.winTeam = 1
	s.chars[0] = []*Char{{name: "Ryu"}}
	s.chars[1] = []*Char{{name: "Ken"}}
	side, winner, loser := s.terminalWinnerData()
	if side != 2 || winner != "ken" || loser != "ryu" {
		t.Fatalf("unexpected winner data: side=%d winner=%q loser=%q", side, winner, loser)
	}
}
