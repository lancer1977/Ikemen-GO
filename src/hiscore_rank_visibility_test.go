package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRankingWouldPlace_ReturnsTrueWhenVisibleWindowHasRoom(t *testing.T) {
	tempDir := t.TempDir()
	statsPath := filepath.Join(tempDir, "stats.json")
	if err := os.WriteFile(statsPath, []byte(`{
		"modes": {
			"arcade": {
				"ranking": []
			}
		}
	}`), 0o644); err != nil {
		t.Fatalf("write stats: %v", err)
	}

	prevSys := sys
	sys = System{
		SystemStateVars: SystemStateVars{
			maxRoundTime: 120,
		},
		cmdFlags: map[string]string{
			"-stats": statsPath,
		},
	}
	sys.statsLog.Matches = []StatsMatch{{MatchTime: 60, WinSide: 0, Wins: [2]int32{1, 0}, TotalScore: [2]int32{4000, 0}}}
	sys.timerRounds = []int32{60}
	sys.scoreStart = [2]float32{0, 0}
	sys.scoreRounds = [][2]float32{{4000, 0}}
	sys.sel.gameParams.RankingCondition = false
	sys.motif.HiscoreInfo.Ranking = map[string]string{"arcade": "score"}
	sys.motif.HiscoreInfo.Window.VisibleItems = 1
	t.Cleanup(func() {
		sys = prevSys
	})

	if !rankingWouldPlace("arcade") {
		t.Fatalf("expected rankingWouldPlace to accept entry when visible window has room")
	}
}

func TestRankingWouldPlace_ReturnsFalseWhenEntryFallsOutOfWindow(t *testing.T) {
	tempDir := t.TempDir()
	statsPath := filepath.Join(tempDir, "stats.json")
	if err := os.WriteFile(statsPath, []byte(`{
		"modes": {
			"arcade": {
				"ranking": [
					{"score": 8000, "time": 1.0, "win": 2, "name": "AAA", "chars": ["old"], "tmode": 0, "ailevel": 4}
				]
			}
		}
	}`), 0o644); err != nil {
		t.Fatalf("write stats: %v", err)
	}

	prevSys := sys
	sys = System{
		SystemStateVars: SystemStateVars{
			maxRoundTime: 120,
		},
		cmdFlags: map[string]string{
			"-stats": statsPath,
		},
	}
	sys.statsLog.Matches = []StatsMatch{{MatchTime: 60, WinSide: 0, Wins: [2]int32{1, 0}, TotalScore: [2]int32{4000, 0}}}
	sys.timerRounds = []int32{60}
	sys.scoreStart = [2]float32{0, 0}
	sys.scoreRounds = [][2]float32{{4000, 0}}
	sys.sel.gameParams.RankingCondition = false
	sys.motif.HiscoreInfo.Ranking = map[string]string{"arcade": "score"}
	sys.motif.HiscoreInfo.Window.VisibleItems = 1
	t.Cleanup(func() {
		sys = prevSys
	})

	if rankingWouldPlace("arcade") {
		t.Fatalf("expected rankingWouldPlace to reject entry after truncation")
	}
}

func TestModeCleared_RespectsRankingConditionAndResultsScreenRoundTarget(t *testing.T) {
	prevSys := sys
	sys = System{}
	t.Cleanup(func() {
		sys = prevSys
	})

	sys.sel.gameParams.RankingCondition = true
	if modeCleared("arcade", 99) {
		t.Fatalf("expected rankingCondition to suppress clearance")
	}

	sys.sel.gameParams.RankingCondition = false
	sys.motif.WinScreen.Results = map[string]string{"arcade": "survival"}
	sys.motif.ResultsScreen = map[string]*ResultsScreenProperties{
		"survival": {
			Enabled:     true,
			RoundsToWin: 3,
		},
	}
	sys.winTeam = 0
	sys.endMatch = false
	sys.round = 4
	if modeCleared("arcade", 2) {
		t.Fatalf("expected modeCleared to reject before rounds-to-win is met")
	}
	if !modeCleared("arcade", 3) {
		t.Fatalf("expected modeCleared to accept once rounds-to-win is reached")
	}
}

func TestResultsScreenForMode_ReturnsNilForMissingOrDisabledScreens(t *testing.T) {
	prevSys := sys
	sys = System{}
	t.Cleanup(func() {
		sys = prevSys
	})

	sys.motif.WinScreen.Results = map[string]string{
		"arcade": "missing",
		"story":  "disabled",
	}
	sys.motif.ResultsScreen = map[string]*ResultsScreenProperties{
		"disabled": {Enabled: false, RoundsToWin: 3},
	}

	if got := resultsScreenForMode("arcade"); got != nil {
		t.Fatalf("expected missing results screen to return nil, got %#v", got)
	}
	if got := resultsScreenForMode("story"); got != nil {
		t.Fatalf("expected disabled results screen to return nil, got %#v", got)
	}
}
