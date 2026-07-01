package main

import "testing"

func TestMatchOverAndFinalRoundConditions(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.matchWins = [2]int32{2, 2}
	sys.wins = [2]int32{0, 2}
	if sys.matchOver() {
		t.Fatal("matchOver() should ignore a team with zero wins")
	}

	sys.wins = [2]int32{2, 0}
	if sys.matchOver() {
		t.Fatal("matchOver() should ignore a team with zero wins on the other side")
	}

	sys.wins = [2]int32{2, 1}
	if !sys.matchOver() {
		t.Fatal("matchOver() should be true once either team reaches its threshold")
	}

	sys.sel.gameParams.PersistRounds = true
	if sys.roundIsFinal() {
		t.Fatal("roundIsFinal() should be false when rounds persist")
	}

	sys.sel.gameParams.PersistRounds = false
	sys.round = 1
	sys.decisiveRound = [2]bool{true, true}
	sys.maxDraws = [2]int32{1, 1}
	sys.draws = 1
	if sys.roundIsFinal() {
		t.Fatal("roundIsFinal() should be false on the first round")
	}

	sys.round = 2
	sys.decisiveRound = [2]bool{true, false}
	if sys.roundIsFinal() {
		t.Fatal("roundIsFinal() should require both teams to be decisive")
	}

	sys.decisiveRound = [2]bool{true, true}
	sys.maxDraws = [2]int32{1, 2}
	sys.draws = 2
	if sys.roundIsFinal() {
		t.Fatal("roundIsFinal() should require both teams to reach max draws")
	}

	sys.maxDraws = [2]int32{1, 1}
	sys.draws = 1
	if !sys.roundIsFinal() {
		t.Fatal("roundIsFinal() should be true when all final-round conditions are met")
	}
}

func TestWinnerTeamResolvesFromMatchAndRoundState(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.endMatch = false
	sys.winTeam = -1
	sys.matchWins = [2]int32{2, 2}
	sys.intro = -5
	sys.fightScreen.round.over_waittime = 1
	sys.fightScreen.round.over_hittime = 1

	sys.wins = [2]int32{2, 1}
	if got := sys.winnerTeam(); got != 1 {
		t.Fatalf("winnerTeam() for team 1 match win = %v, want 1", got)
	}

	sys.wins = [2]int32{1, 2}
	if got := sys.winnerTeam(); got != 2 {
		t.Fatalf("winnerTeam() for team 2 match win = %v, want 2", got)
	}

	sys.wins = [2]int32{2, 2}
	if got := sys.winnerTeam(); got != 0 {
		t.Fatalf("winnerTeam() for tied match win = %v, want 0", got)
	}

	sys.wins = [2]int32{1, 1}
	sys.intro = -2
	sys.winTeam = 1
	if got := sys.winnerTeam(); got != 2 {
		t.Fatalf("winnerTeam() with explicit round winner = %v, want 2", got)
	}

	sys.winTeam = -1
	sys.intro = -2
	sys.finishType = FT_NotYet
	if got := sys.winnerTeam(); got != -1 {
		t.Fatalf("winnerTeam() before round resolution = %v, want -1", got)
	}
}
