package main

import "testing"

func TestRound2AndDefKey(t *testing.T) {
	t.Parallel()

	if got := round2(12.345); got != 12.35 {
		t.Fatalf("round2(12.345) = %v, want 12.35", got)
	}
	if got := round2(12.344); got != 12.34 {
		t.Fatalf("round2(12.344) = %v, want 12.34", got)
	}
	if got := defKey(`chars\KFM\kfm.def`); got != "kfm" {
		t.Fatalf("defKey = %q, want %q", got, "kfm")
	}
}

func TestRankingTypeFor(t *testing.T) {
	t.Parallel()

	prevSys := sys
	sys = System{}
	t.Cleanup(func() { sys = prevSys })

	if got, ok := rankingTypeFor("arcade"); ok || got != "" {
		t.Fatalf("rankingTypeFor with nil map = (%q, %v), want (, false)", got, ok)
	}
	sys.motif.HiscoreInfo.Ranking = map[string]string{"arcade": "score"}
	if got, ok := rankingTypeFor("arcade"); !ok || got != "score" {
		t.Fatalf("rankingTypeFor = (%q, %v), want (score, true)", got, ok)
	}
}

func TestTallyRun(t *testing.T) {
	t.Parallel()

	prevSys := sys
	sys = System{}
	t.Cleanup(func() { sys = prevSys })

	sys.statsLog.Matches = []StatsMatch{
		{MatchTime: 60, WinSide: 0, Wins: [2]int32{1, 0}, TotalScore: [2]int32{1000, 0}},
		{MatchTime: 90, WinSide: 1, Wins: [2]int32{0, 1}, TotalScore: [2]int32{2500, 0}},
	}

	got := tallyRun()
	if got.timeTicks != 150 {
		t.Fatalf("timeTicks = %d, want 150", got.timeTicks)
	}
	if got.winP1 != 1 || got.loseP1 != 1 {
		t.Fatalf("wins/losses = %d/%d, want 1/1", got.winP1, got.loseP1)
	}
	if got.consecP1 != 0 {
		t.Fatalf("consecP1 = %d, want 0 after loss", got.consecP1)
	}
	if got.scoreP1 != 2500 {
		t.Fatalf("scoreP1 = %d, want 2500", got.scoreP1)
	}
}
