package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestComputeAndSaveRanking_WritesUpdatedStatsFile(t *testing.T) {
	tempDir := t.TempDir()
	statsPath := filepath.Join(tempDir, "stats.json")
	if err := os.WriteFile(statsPath, []byte(`{
		"playtime": 12.5,
		"modes": {
			"arcade": {
				"playtime": 3.25,
				"clear": 1,
				"ranking": [
					{"score": 5000, "time": 2.0, "win": 1, "lose": 0, "consecutive": 1, "name": "AAA", "chars": ["old"], "tmode": 0, "ailevel": 4}
				]
			}
		}
	}`), 0o644); err != nil {
		t.Fatalf("write stats: %v", err)
	}

	prevSys := sys
	sys = System{
		SystemStateVars: SystemStateVars{
			match:        12,
			winTeam:      0,
			maxRoundTime: 120,
		},
		cmdFlags: map[string]string{
			"-stats": statsPath,
		},
	}
	sys.statsLog.Matches = []StatsMatch{
		{
			MatchTime:  120,
			WinSide:    0,
			Wins:       [2]int32{1, 0},
			TotalScore: [2]int32{9000, 0},
		},
	}
	sys.timerRounds = []int32{120}
	sys.scoreStart = [2]float32{100, 0}
	sys.scoreRounds = [][2]float32{{8900, 0}}
	sys.sel.gameParams.RankingCondition = false
	sys.sel.charlist = []SelectChar{{def: "chars/ryu/ryu.def"}}
	sys.sel.selected[0] = [][2]int{{0, 3}}
	sys.tmode[0] = TM_Single
	sys.cfg.Options.Difficulty = 6
	sys.motif.HiscoreInfo.Ranking = map[string]string{"arcade": "score"}
	sys.motif.HiscoreInfo.Window.VisibleItems = 10
	t.Cleanup(func() {
		sys = prevSys
	})

	cleared, place := computeAndSaveRanking("arcade")
	if !cleared || place != 1 {
		t.Fatalf("unexpected ranking result: cleared=%v place=%d", cleared, place)
	}

	raw, err := os.ReadFile(statsPath)
	if err != nil {
		t.Fatalf("read stats: %v", err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("unmarshal stats: %v", err)
	}
	if got := doc["playtime"].(float64); got != 14.5 {
		t.Fatalf("unexpected playtime: %v", got)
	}
	modes := doc["modes"].(map[string]any)
	arcade := modes["arcade"].(map[string]any)
	if got := arcade["playtime"].(float64); got != 5.25 {
		t.Fatalf("unexpected mode playtime: %v", got)
	}
	if got := arcade["clear"].(float64); got != 2 {
		t.Fatalf("unexpected clear count: %v", got)
	}
	ranking := arcade["ranking"].([]any)
	if len(ranking) != 2 {
		t.Fatalf("unexpected ranking length: %d", len(ranking))
	}
	newEntry := ranking[0].(map[string]any)
	if newEntry["name"].(string) != "" || int(newEntry["score"].(float64)) != 9000 {
		t.Fatalf("unexpected new ranking entry: %#v", newEntry)
	}
}

func TestComputeAndSaveRanking_BootstrapsMissingStatsFile(t *testing.T) {
	tempDir := t.TempDir()
	statsPath := filepath.Join(tempDir, "missing", "stats.json")

	prevSys := sys
	sys = System{
		SystemStateVars: SystemStateVars{
			match:        12,
			winTeam:      0,
			maxRoundTime: 120,
		},
		cmdFlags: map[string]string{
			"-stats": statsPath,
		},
	}
	sys.statsLog.Matches = []StatsMatch{
		{
			MatchTime:  120,
			WinSide:    0,
			Wins:       [2]int32{1, 0},
			TotalScore: [2]int32{9000, 0},
		},
	}
	sys.timerRounds = []int32{120}
	sys.scoreStart = [2]float32{0, 0}
	sys.scoreRounds = [][2]float32{{9000, 0}}
	sys.sel.gameParams.RankingCondition = false
	sys.sel.charlist = []SelectChar{{def: "chars/ryu/ryu.def"}}
	sys.sel.selected[0] = [][2]int{{0, 3}}
	sys.tmode[0] = TM_Single
	sys.cfg.Options.Difficulty = 6
	sys.motif.HiscoreInfo.Ranking = map[string]string{"arcade": "score"}
	sys.motif.HiscoreInfo.Window.VisibleItems = 10
	t.Cleanup(func() {
		sys = prevSys
	})

	cleared, place := computeAndSaveRanking("arcade")
	if !cleared || place != 1 {
		t.Fatalf("unexpected ranking result: cleared=%v place=%d", cleared, place)
	}

	if _, err := os.Stat(statsPath); err != nil {
		t.Fatalf("expected stats file to be created: %v", err)
	}
}
