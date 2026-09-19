package main

import "testing"

func TestMatchOverAndFinalRoundConditions(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.matchWins = [2]int32{2, 2}
	sys.wins = [2]int32{0, 2}
	// matchOver checks each team independently: (wins[0]>0 && wins[0]>=threshold) || (wins[1]>0 && wins[1]>=threshold)
	// Team 1 has 2 wins >= 2 threshold, so it returns true regardless of team 0 being 0
	if !sys.matchOver() {
		t.Fatal("matchOver() should return true when team 1 reaches its threshold even if team 0 has zero wins")
	}

	sys.wins = [2]int32{2, 0}
	// Team 0 has 2 wins >= 2 threshold, so it returns true regardless of team 1 being 0
	if !sys.matchOver() {
		t.Fatal("matchOver() should return true when team 0 reaches its threshold even if team 1 has zero wins")
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
	// maxDrawsReached checks if draws >= maxDraws[team], so both teams need their limit reached
	// With draws=2, team 0 limit=1 (reached: 2>=1), but team 1 limit=3 (not reached: 2<3)
	sys.maxDraws = [2]int32{1, 3}
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
	ensureGlobalRound(t)
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
