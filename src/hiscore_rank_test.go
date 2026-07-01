package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseRankingRows_ReadsConfiguredStatsFile(t *testing.T) {
	tempDir := t.TempDir()
	statsPath := filepath.Join(tempDir, "stats.json")
	if err := os.WriteFile(statsPath, []byte(`{
		"modes": {
			"arcade": {
				"ranking": [
					{"score": 9000, "time": 3.5, "win": 2, "name": "RYU", "chars": ["ryu", "ken"]},
					{"score": 4500, "time": 4.25, "win": 1, "name": "KEN", "chars": ["ken"]}
				]
			}
		}
	}`), 0o644); err != nil {
		t.Fatalf("write stats: %v", err)
	}

	rows := parseRankingRows(statsPath, "arcade")
	if len(rows) != 2 {
		t.Fatalf("unexpected ranking row count: %d", len(rows))
	}
	if rows[0].score != 9000 || rows[0].name != "RYU" || len(rows[0].chars) != 2 {
		t.Fatalf("unexpected first ranking row: %#v", rows[0])
	}
	if rows[1].time != 4.25 || rows[1].chars[0] != "ken" {
		t.Fatalf("unexpected second ranking row: %#v", rows[1])
	}
}

func TestParseRankingRows_ReturnsNilForMissingOrBlankInputs(t *testing.T) {
	if got := parseRankingRows("", "arcade"); got != nil {
		t.Fatalf("expected nil for blank path, got %#v", got)
	}
	if got := parseRankingRows("/definitely/missing.json", "arcade"); got != nil {
		t.Fatalf("expected nil for missing path, got %#v", got)
	}
	if got := parseRankingRows("/definitely/missing.json", ""); got != nil {
		t.Fatalf("expected nil for blank mode, got %#v", got)
	}
}
