package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMaybeWriteLiveSnapshot_WritesSnapshotDuringMatch(t *testing.T) {
	ensureGlobalRound(t)
	tempDir := t.TempDir()
	livePath := filepath.Join(tempDir, "live_data.json")
	statusPath := filepath.Join(tempDir, "live_status.json")

	s := &System{
		SystemStateVars: SystemStateVars{
			match:        12,
			round:        2,
			matchTime:    180,
			curRoundTime: 90,
			fightLoopEnd: false,
			postMatchFlg: false,
		},
		frameCounter: 6,
		cmdFlags: map[string]string{
			"-livedatafile":   livePath,
			"-livestatusfile": statusPath,
		},
	}
	s.chars[0] = []*Char{{name: "Ryu", teamside: 0}}
	s.chars[1] = []*Char{{name: "Ken", teamside: 1}}
	s.fightScreen.scores[0] = &FightScreenScore{scorePoints: 1}
	s.fightScreen.scores[1] = &FightScreenScore{scorePoints: 0}

	s.maybeWriteLiveSnapshot()

	raw, err := os.ReadFile(livePath)
	if err != nil {
		t.Fatalf("reading live snapshot: %v", err)
	}
	var snap GameLiveSnapshot
	if err := json.Unmarshal(raw, &snap); err != nil {
		t.Fatalf("unmarshal live snapshot: %v", err)
	}
	if snap.CurrentRound.Index != 2 || snap.CurrentRound.Timer != 90 {
		t.Fatalf("unexpected snapshot round data: %#v", snap.CurrentRound)
	}
	if snap.CurrentRound.Score != [2]int32{1, 0} {
		t.Fatalf("unexpected snapshot score: %#v", snap.CurrentRound.Score)
	}
	if snap.CurrentRound.Fighters[0][0].Name != "Ryu" || snap.CurrentRound.Fighters[1][0].Name != "Ken" {
		t.Fatalf("unexpected snapshot fighters: %#v", snap.CurrentRound.Fighters)
	}
}

func TestMaybeWriteLiveSnapshot_SkipsOddFramesButWritesStatus(t *testing.T) {
	ensureGlobalRound(t)
	tempDir := t.TempDir()
	livePath := filepath.Join(tempDir, "live_data.json")
	statusPath := filepath.Join(tempDir, "live_status.json")

	s := &System{
		SystemStateVars: SystemStateVars{
			match:        12,
			round:        2,
			matchTime:    180,
			curRoundTime: 90,
		},
		frameCounter: 11,
		cmdFlags: map[string]string{
			"-livedatafile":   livePath,
			"-livestatusfile": statusPath,
		},
	}

	s.maybeWriteLiveSnapshot()

	if _, err := os.Stat(livePath); !os.IsNotExist(err) {
		t.Fatalf("expected no live snapshot on odd frame, got err=%v", err)
	}
	if _, err := os.Stat(statusPath); err != nil {
		t.Fatalf("expected live status fallback to be written: %v", err)
	}
}

func TestMaybeWriteLiveSnapshot_OmitsStatusWhenNotInMatchDueToDeadFallback(t *testing.T) {
	ensureGlobalRound(t)
	tempDir := t.TempDir()
	livePath := filepath.Join(tempDir, "live_data.json")
	statusPath := filepath.Join(tempDir, "live_status.json")

	s := &System{
		SystemStateVars: SystemStateVars{
			match:        12,
			round:        2,
			matchTime:    180,
			curRoundTime: 90,
			// middleOfMatch() is !fightLoopEnd && matchTime != 0 &&
			// !postMatchFlg, so a non-zero matchTime alone still counts as
			// in-match. fightLoopEnd is what actually ends it.
			fightLoopEnd: true,
		},
		frameCounter: 6,
		cmdFlags: map[string]string{
			"-livedatafile":   livePath,
			"-livestatusfile": statusPath,
		},
	}

	s.maybeWriteLiveSnapshot()

	// KNOWN BUG: maybeWriteLiveSnapshot calls writeLiveStatus as a fallback when
	// !middleOfMatch() && !matchOver(), but writeLiveStatus guards with the same
	// condition, so the fallback status file is never written. This is an engine
	// defect that causes status files to be silently omitted in this scenario.
	// See: the duplicate guard conditions in live_snapshot.go:95 and live_eventing.go:362
	if _, err := os.Stat(livePath); !os.IsNotExist(err) {
		t.Fatalf("expected no live snapshot when not in match, got err=%v", err)
	}
	if _, err := os.Stat(statusPath); err == nil {
		t.Fatalf("expected no live status fallback due to dead code path, but file was written")
	} else if !os.IsNotExist(err) {
		t.Fatalf("unexpected error checking status path: %v", err)
	}
}
