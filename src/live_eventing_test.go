package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseLiveCommandInbox_ParsesKnownCommands(t *testing.T) {
	commands, err := parseLiveCommandInbox([]byte(`{
		"schema": "live-lancero/command-inbox/v1",
		"commands": [
			{
				"id": "cmd-power-1",
				"commandKind": "power-adjust",
				"safeTimingPolicy": "immediate",
				"side": 1,
				"amount": 150
			},
			{
				"id": "cmd-swap-1",
				"command": "players-swap",
				"safeTimingPolicy": "next-match"
			}
		]
	}`))
	if err != nil {
		t.Fatalf("parseLiveCommandInbox returned error: %v", err)
	}
	if len(commands) != 2 {
		t.Fatalf("unexpected command count: %d", len(commands))
	}
	if commands[0].ID != "cmd-power-1" || commands[0].Amount != 150 {
		t.Fatalf("unexpected first command: %#v", commands[0])
	}
	if commands[1].ID != "cmd-swap-1" || commands[1].Command != "players-swap" {
		t.Fatalf("unexpected second command: %#v", commands[1])
	}
}

func TestParseLiveCommandInbox_HandlesBlankAndSchemaMismatch(t *testing.T) {
	commands, err := parseLiveCommandInbox([]byte("   \n"))
	if err != nil {
		t.Fatalf("parseLiveCommandInbox returned error for blank input: %v", err)
	}
	if commands != nil {
		t.Fatalf("expected nil commands for blank input, got %#v", commands)
	}

	if _, err := parseLiveCommandInbox([]byte(`{"schema":"wrong","commands":[]}`)); err == nil {
		t.Fatalf("expected schema mismatch to fail")
	}
}

func TestNormalizeLiveCommandKind_FallsBackAndNormalizesSeparators(t *testing.T) {
	if got := normalizeLiveCommandKind(LiveCommandRequest{CommandKind: "power_adjust"}); got != "power-adjust" {
		t.Fatalf("unexpected normalized command kind: %q", got)
	}
	if got := normalizeLiveCommandKind(LiveCommandRequest{Command: "players-swap"}); got != "players-swap" {
		t.Fatalf("unexpected fallback command kind: %q", got)
	}
}

func TestRecordCombatDamage_WritesCombatAndThresholdEvents(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "combat_events.jsonl")
	s := &System{
		SystemStateVars: SystemStateVars{
			match:  12,
			round:  2,
			paused: false,
		},
		cmdFlags: map[string]string{
			"-combateventsfile": path,
		},
	}
	attacker := &Char{playerNo: 0, teamside: 0, name: "Ryu"}
	defender := &Char{playerNo: 1, teamside: 1, name: "Ken"}
	s.chars[0] = []*Char{attacker}
	s.chars[1] = []*Char{defender}
	s.frameCounter = 480

	s.recordCombatDamage(defender, int32(attacker.playerNo), 180, 180, 0)

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading combat events: %v", err)
	}

	lines := splitNonEmptyLines(string(raw))
	if len(lines) != 3 {
		t.Fatalf("unexpected combat event count: %d (%q)", len(lines), string(raw))
	}

	var damage, threshold, ko LiveCombatEvent
	if err := json.Unmarshal([]byte(lines[0]), &damage); err != nil {
		t.Fatalf("unmarshal damage event: %v", err)
	}
	if err := json.Unmarshal([]byte(lines[1]), &threshold); err != nil {
		t.Fatalf("unmarshal threshold event: %v", err)
	}
	if err := json.Unmarshal([]byte(lines[2]), &ko); err != nil {
		t.Fatalf("unmarshal ko event: %v", err)
	}

	if damage.Event != "damage-dealt" || damage.AttackerKey != "Ryu" || damage.DefenderKey != "Ken" || damage.Damage != 180 {
		t.Fatalf("unexpected damage event: %#v", damage)
	}
	if threshold.Event != "threshold-crossed" || threshold.ThresholdKind != "damage-threshold" || threshold.TextHint != "NICE HIT!" {
		t.Fatalf("unexpected threshold event: %#v", threshold)
	}
	if ko.Event != "ko-confirmed" || ko.ThresholdKind != "ko" || ko.TextHint != "KO" {
		t.Fatalf("unexpected ko event: %#v", ko)
	}
}

func TestRecordCombatDamage_IgnoresNilDefenderAndZeroDamage(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "combat_events.jsonl")
	s := &System{
		cmdFlags: map[string]string{
			"-combateventsfile": path,
		},
	}

	s.recordCombatDamage(nil, 0, 0, 100, 100)
	s.recordCombatDamage(&Char{playerNo: 0, teamside: 0, name: "Ryu"}, 0, 0, 100, 100)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected no combat events to be written, got err=%v", err)
	}
}

func TestLiveSideAndRosterKeyHelpers_MapSidesAndNames(t *testing.T) {
	if got := liveSideFromChar(nil); got != 0 {
		t.Fatalf("expected nil char to map to 0, got %d", got)
	}
	if got := liveSideFromChar(&Char{teamside: -1}); got != 0 {
		t.Fatalf("expected invalid teamside to map to 0, got %d", got)
	}
	if got := liveSideFromChar(&Char{teamside: 0}); got != 1 {
		t.Fatalf("expected side 0 to map to 1, got %d", got)
	}
	if got := liveSideFromPlayerNo(-1); got != 0 {
		t.Fatalf("expected negative playerNo to map to 0, got %d", got)
	}
	if got := liveSideFromPlayerNo(0); got != 1 {
		t.Fatalf("expected player 0 to map to 1, got %d", got)
	}
	if got := liveSideFromPlayerNo(1); got != 2 {
		t.Fatalf("expected player 1 to map to 2, got %d", got)
	}
	if got := liveRosterKey(nil); got != "" {
		t.Fatalf("expected nil roster key to be blank, got %q", got)
	}
	if got := liveRosterKey(&Char{name: "Ryu"}); got != "Ryu" {
		t.Fatalf("unexpected roster key: %q", got)
	}

	s := &System{}
	if got := s.liveRosterKeyForSide(-1); got != "" {
		t.Fatalf("expected invalid side to return blank roster key, got %q", got)
	}
	s.chars[0] = []*Char{{name: "Ryu"}}
	if got := s.liveRosterKeyForSide(0); got != "ryu" {
		t.Fatalf("unexpected side roster key: %q", got)
	}
}

func TestWriteLiveStatus_WritesCurrentFightSnapshot(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "live_status.json")
	s := &System{
		SystemStateVars: SystemStateVars{
			match:     12,
			round:     2,
			matchTime: 180,
		},
		stage: &Stage{name: "training_ground", displayname: "Training Ground"},
		cmdFlags: map[string]string{
			"-livestatusfile": path,
		},
	}
	s.chars[0] = []*Char{{name: "Ryu"}}
	s.chars[1] = []*Char{{name: "Ken"}}

	s.writeLiveStatus()

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading live status: %v", err)
	}

	var status LiveStatusSnapshot
	if err := json.Unmarshal(raw, &status); err != nil {
		t.Fatalf("unmarshal live status: %v", err)
	}

	if status.Schema != liveStatusSchema {
		t.Fatalf("unexpected schema: %#v", status)
	}
	if status.Mode != "fight" || status.Match != 12 || status.Round != 2 {
		t.Fatalf("unexpected status fields: %#v", status)
	}
	if status.Stage != "Training Ground" {
		t.Fatalf("unexpected stage: %#v", status)
	}
	if status.P1 == nil || status.P1.Key != "ryu" || status.P2 == nil || status.P2.Key != "ken" {
		t.Fatalf("unexpected fighters: %#v", status)
	}
}

func TestWriteLiveStatus_UsesResultModeAfterMatchOver(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "live_status.json")
	s := &System{
		SystemStateVars: SystemStateVars{
			match:     12,
			round:     2,
			matchTime: 180,
			wins:      [2]int32{1, 0},
			matchWins: [2]int32{1, 1},
		},
		cmdFlags: map[string]string{
			"-livestatusfile": path,
		},
	}
	s.chars[0] = []*Char{{name: "Ryu"}}
	s.chars[1] = []*Char{{name: "Ken"}}

	s.writeLiveStatus()

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading live status: %v", err)
	}

	var status LiveStatusSnapshot
	if err := json.Unmarshal(raw, &status); err != nil {
		t.Fatalf("unmarshal live status: %v", err)
	}

	if status.Mode != "result" {
		t.Fatalf("unexpected status mode: %#v", status)
	}
}

func TestMaybeWriteTerminalLiveArtifacts_WritesResultAndMatchComplete(t *testing.T) {
	tempDir := t.TempDir()
	resultPath := filepath.Join(tempDir, "result.json")
	matchEventsPath := filepath.Join(tempDir, "match_events.jsonl")
	s := &System{
		SystemStateVars: SystemStateVars{
			match:     12,
			round:     2,
			matchTime: 180,
			winTeam:   1,
			wins:      [2]int32{0, 1},
			matchWins: [2]int32{1, 1},
		},
		cmdFlags: map[string]string{
			"-resultfile":      resultPath,
			"-matcheventsfile": matchEventsPath,
		},
	}
	s.statsLog.Matches = []StatsMatch{{WinSide: 1, Ended: true}}
	origSys := sys
	sys = *s
	t.Cleanup(func() {
		sys = origSys
	})

	s.maybeWriteTerminalLiveArtifacts()
	s.maybeWriteTerminalLiveArtifacts()

	resultRaw, err := os.ReadFile(resultPath)
	if err != nil {
		t.Fatalf("reading result file: %v", err)
	}
	var result GameLiveSnapshot
	if err := json.Unmarshal(resultRaw, &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if !result.MatchOver {
		t.Fatalf("expected match over result: %#v", result)
	}

	matchRaw, err := os.ReadFile(matchEventsPath)
	if err != nil {
		t.Fatalf("reading match events: %v", err)
	}
	lines := splitNonEmptyLines(string(matchRaw))
	if len(lines) != 1 {
		t.Fatalf("unexpected match event count: %d (%q)", len(lines), string(matchRaw))
	}
	var event LiveMatchEvent
	if err := json.Unmarshal([]byte(lines[0]), &event); err != nil {
		t.Fatalf("unmarshal match event: %v", err)
	}
	if event.Event != "match_complete" || event.WinnerSide != 2 {
		t.Fatalf("unexpected match complete event: %#v", event)
	}
}

func TestTerminalWinnerData_UsesLiveRosterKeys(t *testing.T) {
	s := &System{}
	s.winTeam = 1
	s.chars[0] = []*Char{{name: "Ryu"}}
	s.chars[1] = []*Char{{name: "Ken"}}

	winnerSide, winnerKey, loserKey := s.terminalWinnerData()

	if winnerSide != 2 || winnerKey != "ken" || loserKey != "ryu" {
		t.Fatalf("unexpected terminal winner data: side=%d winner=%q loser=%q", winnerSide, winnerKey, loserKey)
	}
}

func TestMaybeWriteTerminalLiveArtifacts_WritesRichFightWhenEnabled(t *testing.T) {
	tempDir := t.TempDir()
	richPath := filepath.Join(tempDir, "fight_history.jsonl")
	s := &System{
		SystemStateVars: SystemStateVars{
			match:      12,
			round:      2,
			matchTime:  180,
			randseed:   12345,
			winTeam:    0,
			wins:       [2]int32{1, 0},
			matchWins:  [2]int32{1, 1},
			gameWidth:  320,
			gameHeight: 240,
			scrrect:    [4]int32{0, 0, 640, 480},
		},
		stage: &Stage{
			def:         "stages/training.def",
			name:        "training_ground",
			displayname: "Training Ground",
			stageCamera: stageCamera{localcoord: [2]int32{320, 240}},
			scale:       [2]float32{1, 1},
		},
		cmdFlags: map[string]string{
			"-richfightfile": richPath,
		},
	}
	s.sel.charlist = []SelectChar{
		{def: "chars/ryu/ryu.def", name: "Ryu", lifebarname: "Ryu", author: "Capcom", localcoord: [2]int32{320, 240}, cns_scale: [2]float32{1, 1}},
		{def: "chars/ken/ken.def", name: "Ken", lifebarname: "Ken", author: "Capcom", localcoord: [2]int32{320, 240}, cns_scale: [2]float32{1, 1}},
	}
	s.sel.selected[0] = [][2]int{{0, 3}}
	s.sel.selected[1] = [][2]int{{1, 4}}
	s.sel.selectedStageNo = 1
	s.sel.stagelist = []SelectStage{{def: "stages/training.def", name: "Training Ground", localcoord: [2]int32{320, 240}}}
	origSys := sys
	sys = *s
	t.Cleanup(func() {
		sys = origSys
	})

	s.maybeWriteTerminalLiveArtifacts()
	s.maybeWriteTerminalLiveArtifacts()

	raw, err := os.ReadFile(richPath)
	if err != nil {
		t.Fatalf("reading rich fight history: %v", err)
	}
	lines := splitNonEmptyLines(string(raw))
	if len(lines) != 1 {
		t.Fatalf("unexpected rich fight count: %d (%q)", len(lines), string(raw))
	}
	var record LiveRichFightRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("unmarshal rich fight record: %v", err)
	}
	if record.Schema != liveRichFightSchema || record.Mode != "result" || record.Stage != "Training Ground" {
		t.Fatalf("unexpected rich fight record: %#v", record)
	}
	if record.WinnerSide != 1 || !record.Snapshot.MatchOver {
		t.Fatalf("unexpected rich fight winner/snapshot: %#v", record)
	}
	if record.Context.RandSeed != 12345 || record.Context.Stage.Def != "stages/training.def" {
		t.Fatalf("unexpected rich fight context: %#v", record.Context)
	}
	if got := record.Context.SelectedTeams[0][0]; got.Def != "chars/ryu/ryu.def" || got.PaletteNo != 3 {
		t.Fatalf("unexpected p1 rich fight selection context: %#v", got)
	}
	if got := record.Context.SelectedTeams[1][0]; got.Def != "chars/ken/ken.def" || got.PaletteNo != 4 {
		t.Fatalf("unexpected p2 rich fight selection context: %#v", got)
	}
}

func TestRecordRoundOutcome_WritesMatchEvent(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "match_events.jsonl")
	s := &System{
		SystemStateVars: SystemStateVars{
			match: 12,
			round: 2,
		},
		cmdFlags: map[string]string{
			"-matcheventsfile": path,
		},
	}
	s.chars[0] = []*Char{{name: "Ryu"}}
	s.chars[1] = []*Char{{name: "Ken"}}
	s.winTeam = 0

	s.recordRoundOutcome()

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading match events: %v", err)
	}

	lines := splitNonEmptyLines(string(raw))
	if len(lines) != 1 {
		t.Fatalf("unexpected match event count: %d (%q)", len(lines), string(raw))
	}

	var event LiveMatchEvent
	if err := json.Unmarshal([]byte(lines[0]), &event); err != nil {
		t.Fatalf("unmarshal match event: %v", err)
	}

	if event.Schema != liveMatchEventSchema || event.Event != "round_win" || event.WinnerSide != 1 {
		t.Fatalf("unexpected match event: %#v", event)
	}
	if event.WinnerKey != "ryu" || event.LoserKey != "ken" {
		t.Fatalf("unexpected winner data: %#v", event)
	}
}

func TestRecordRoundStart_WritesMatchEvent(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "match_events.jsonl")
	s := &System{
		SystemStateVars: SystemStateVars{
			match: 12,
			round: 1,
		},
		cmdFlags: map[string]string{
			"-matcheventsfile": path,
		},
	}

	s.recordRoundStart()

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading match events: %v", err)
	}

	lines := splitNonEmptyLines(string(raw))
	if len(lines) != 1 {
		t.Fatalf("unexpected match event count: %d (%q)", len(lines), string(raw))
	}

	var event LiveMatchEvent
	if err := json.Unmarshal([]byte(lines[0]), &event); err != nil {
		t.Fatalf("unmarshal match event: %v", err)
	}

	if event.Event != "round_start" || event.Match != 12 || event.Round != 1 {
		t.Fatalf("unexpected round start event: %#v", event)
	}
}

func TestMaybeProcessLiveCommandInbox_AppliesPowerAdjustAndWritesResult(t *testing.T) {
	origSys := sys
	sys = System{}
	t.Cleanup(func() {
		sys = origSys
	})

	tempDir := t.TempDir()
	inboxPath := filepath.Join(tempDir, "command_inbox.json")
	resultPath := filepath.Join(tempDir, "command_results.jsonl")
	if err := os.WriteFile(inboxPath, []byte(`{
		"schema": "live-lancero/command-inbox/v1",
		"commands": [
			{
				"id": "cmd-power-1",
				"commandKind": "power-adjust",
				"safeTimingPolicy": "immediate",
				"side": 1,
				"amount": 150
			}
		]
	}`), 0o644); err != nil {
		t.Fatalf("write inbox: %v", err)
	}

	root := &Char{playerNo: 0, teamside: 0, name: "Ryu", power: 100, powerMax: 1000}
	sys.chars[0] = []*Char{root}
	sys.intro = 0
	sys.maxPowerMode = false

	s := &System{
		SystemStateVars: SystemStateVars{
			match:        12,
			round:        1,
			matchTime:    1,
			postMatchFlg: false,
			fightLoopEnd: false,
		},
		cmdFlags: map[string]string{
			"-commandinboxfile":   inboxPath,
			"-commandresultsfile": resultPath,
		},
	}
	s.chars[0] = []*Char{root}
	s.frameCounter = 6

	s.maybeProcessLiveCommandInbox()

	if got := root.power; got != 250 {
		t.Fatalf("unexpected power after command: %d", got)
	}

	raw, err := os.ReadFile(resultPath)
	if err != nil {
		t.Fatalf("reading command results: %v", err)
	}
	lines := splitNonEmptyLines(string(raw))
	if len(lines) != 1 {
		t.Fatalf("unexpected result count: %d (%q)", len(lines), string(raw))
	}

	var result LiveCommandResult
	if err := json.Unmarshal([]byte(lines[0]), &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	if result.CommandID != "cmd-power-1" || result.CommandKind != "power-adjust" || result.Status != "applied" {
		t.Fatalf("unexpected command result: %#v", result)
	}
	if result.Side != 1 {
		t.Fatalf("unexpected command result side: %#v", result)
	}
	if result.Schema != liveCommandResultSchema {
		t.Fatalf("unexpected command result schema: %q", result.Schema)
	}
}

func TestMaybeProcessLiveCommandInbox_AppliesSkipRoundAndWritesResult(t *testing.T) {
	origSys := sys
	sys = System{}
	t.Cleanup(func() {
		sys = origSys
	})

	tempDir := t.TempDir()
	inboxPath := filepath.Join(tempDir, "command_inbox.json")
	resultPath := filepath.Join(tempDir, "command_results.jsonl")
	if err := os.WriteFile(inboxPath, []byte(`{
		"schema": "live-lancero/command-inbox/v1",
		"commands": [
			{
				"id": "cmd-skip-1",
				"command": "skip-round",
				"safeTimingPolicy": "immediate",
				"side": 1
			}
		]
	}`), 0o644); err != nil {
		t.Fatalf("write inbox: %v", err)
	}

	p1 := &Char{playerNo: 0, teamside: 0, helperIndex: 0, life: 900, lifeMax: 1000, redLife: 900}
	p2 := &Char{playerNo: 1, teamside: 1, helperIndex: 0, life: 850, lifeMax: 1000, redLife: 850}
	sys.chars[0] = []*Char{p1}
	sys.chars[1] = []*Char{p2}

	s := &System{
		SystemStateVars: SystemStateVars{
			match:        12,
			round:        1,
			matchTime:    1,
			winTeam:      -1,
			postMatchFlg: false,
			fightLoopEnd: false,
		},
		cmdFlags: map[string]string{
			"-commandinboxfile":   inboxPath,
			"-commandresultsfile": resultPath,
		},
	}
	s.chars[0] = []*Char{p1}
	s.chars[1] = []*Char{p2}
	s.frameCounter = 6

	s.maybeProcessLiveCommandInbox()

	if p2.life != 0 || p2.redLife != 0 || p2.alive() {
		t.Fatalf("skip-round should KO P2 from inbox, got life=%d redLife=%d alive=%v", p2.life, p2.redLife, p2.alive())
	}
	if s.finishType != FT_KO || s.winTeam != 0 {
		t.Fatalf("skip-round should resolve P1 as winner, finishType=%v winTeam=%d", s.finishType, s.winTeam)
	}

	raw, err := os.ReadFile(resultPath)
	if err != nil {
		t.Fatalf("reading command results: %v", err)
	}
	lines := splitNonEmptyLines(string(raw))
	if len(lines) != 1 {
		t.Fatalf("unexpected result count: %d (%q)", len(lines), string(raw))
	}

	var result LiveCommandResult
	if err := json.Unmarshal([]byte(lines[0]), &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if result.CommandID != "cmd-skip-1" || result.CommandKind != "skip-round" || result.Status != "applied" {
		t.Fatalf("unexpected command result: %#v", result)
	}
	if result.Side != 1 {
		t.Fatalf("unexpected command result side: %#v", result)
	}
	if result.AppliedRound != 1 || result.AppliedFrame != 6 {
		t.Fatalf("unexpected application coordinates: %#v", result)
	}
}

func TestApplyLiveCommand_AutoKillSetsTargetLifeToZero(t *testing.T) {
	origSys := sys
	sys = System{}
	t.Cleanup(func() {
		sys = origSys
	})

	root := &Char{
		playerNo:    0,
		teamside:    0,
		helperIndex: 0,
		life:        750,
		lifeMax:     1000,
		redLife:     750,
	}
	sys.chars[0] = []*Char{root}

	s := &System{
		SystemStateVars: SystemStateVars{
			match:     12,
			round:     1,
			matchTime: 1,
		},
	}
	s.chars[0] = []*Char{root}
	s.frameCounter = 6

	result := s.applyLiveCommand(LiveCommandRequest{ID: "cmd-kill-p1", CommandKind: "auto-kill", Side: 1})
	if result.Status != "applied" || result.CommandKind != "auto-kill" {
		t.Fatalf("unexpected auto-kill result: %#v", result)
	}
	if root.life != 0 || root.redLife != 0 {
		t.Fatalf("auto-kill should zero target life, got life=%d redLife=%d", root.life, root.redLife)
	}
}

func TestApplyLiveCommand_SkipRoundKOsOpposingSideAndResolvesWinner(t *testing.T) {
	origSys := sys
	sys = System{}
	t.Cleanup(func() {
		sys = origSys
	})

	p1 := &Char{playerNo: 0, teamside: 0, helperIndex: 0, life: 750, lifeMax: 1000, redLife: 750}
	p2 := &Char{playerNo: 1, teamside: 1, helperIndex: 0, life: 800, lifeMax: 1000, redLife: 800}
	sys.chars[0] = []*Char{p1}
	sys.chars[1] = []*Char{p2}

	s := &System{
		SystemStateVars: SystemStateVars{
			match:     12,
			round:     2,
			matchTime: 1,
			winTeam:   -1,
		},
	}
	s.chars[0] = []*Char{p1}
	s.chars[1] = []*Char{p2}
	s.frameCounter = 12

	result := s.applyLiveCommand(LiveCommandRequest{ID: "cmd-skip-p1", CommandKind: "skip-round", Side: 1})
	if result.Status != "applied" || result.CommandKind != "skip-round" {
		t.Fatalf("unexpected skip-round result: %#v", result)
	}
	if p2.life != 0 || p2.redLife != 0 || p2.alive() {
		t.Fatalf("skip-round should KO opposing side, got life=%d redLife=%d alive=%v", p2.life, p2.redLife, p2.alive())
	}
	if p1.life == 0 || !p1.alive() {
		t.Fatalf("skip-round should not KO winning side, got life=%d alive=%v", p1.life, p1.alive())
	}
	if s.finishType != FT_KO || s.winTeam != 0 {
		t.Fatalf("skip-round should resolve as P1 KO win, finishType=%v winTeam=%d", s.finishType, s.winTeam)
	}
}

func TestApplyLiveCommand_ReportsValidationFailures(t *testing.T) {
	s := &System{
		SystemStateVars: SystemStateVars{
			match:     12,
			round:     1,
			matchTime: 1,
		},
	}

	result := s.applyLiveCommand(LiveCommandRequest{ID: "bad-side", CommandKind: "power-adjust", Side: 9, Amount: 1})
	if result.Status != "failed" || result.Reason != "invalid-side" {
		t.Fatalf("unexpected invalid-side result: %#v", result)
	}

	result = s.applyLiveCommand(LiveCommandRequest{ID: "bad-amount", CommandKind: "power-adjust", Side: 1, Amount: 4000})
	if result.Status != "failed" || result.Reason != "invalid-amount" {
		t.Fatalf("unexpected invalid-amount result: %#v", result)
	}

	result = s.applyLiveCommand(LiveCommandRequest{ID: "unsupported", CommandKind: "anything-else"})
	if result.Status != "failed" || result.Reason != "unsupported-command" || result.CommandKind != "unknown" {
		t.Fatalf("unexpected unsupported-command result: %#v", result)
	}
}

func TestApplyLiveCommand_HandlesFightPhaseAndSuccessCases(t *testing.T) {
	origSys := sys
	sys = System{}
	t.Cleanup(func() {
		sys = origSys
	})

	s := &System{
		SystemStateVars: SystemStateVars{
			match:     12,
			round:     1,
			matchTime: 0,
		},
	}

	result := s.applyLiveCommand(LiveCommandRequest{ID: "not-in-fight", CommandKind: "power-adjust", Side: 1, Amount: 1})
	if result.Status != "failed" || result.Reason != "not-in-fight-phase" {
		t.Fatalf("unexpected not-in-fight result: %#v", result)
	}

	root := &Char{playerNo: 0, teamside: 0, name: "Ryu", power: 100, powerMax: 1000}
	s.chars[0] = []*Char{root}
	sys.chars[0] = []*Char{root}
	s.matchTime = 1
	s.intro = 0
	s.maxPowerMode = false

	result = s.applyLiveCommand(LiveCommandRequest{ID: "applied", CommandKind: "power-adjust", Side: 1, Amount: 150})
	if result.Status != "applied" || result.Reason != "" {
		t.Fatalf("unexpected applied result: %#v", result)
	}
	if root.power != 250 {
		t.Fatalf("unexpected root power after command: %d", root.power)
	}

	s.chars[0] = nil
	result = s.applyLiveCommand(LiveCommandRequest{ID: "missing-root", CommandKind: "power-adjust", Side: 1, Amount: 1})
	if result.Status != "failed" || result.Reason != "missing-team-root" {
		t.Fatalf("unexpected missing-root result: %#v", result)
	}
}

func TestPathHelpers_DeriveSiblingTelemetryFilesFromConfiguredRoots(t *testing.T) {
	tempDir := t.TempDir()
	liveDataPath := filepath.Join(tempDir, "live_data.json")
	inboxPath := filepath.Join(tempDir, "nested", "command_inbox.json")
	resultPath := filepath.Join(tempDir, "result.json")
	statusPath := filepath.Join(tempDir, "status.json")

	s := &System{
		cmdFlags: map[string]string{
			"-livedatafile":     liveDataPath,
			"-commandinboxfile": inboxPath,
			"-resultfile":       resultPath,
			"-livestatusfile":   statusPath,
		},
	}

	if got := s.combatEventsPath(); got != filepath.Join(tempDir, "combat_events.jsonl") {
		t.Fatalf("unexpected combat events path: %q", got)
	}
	if got := s.matchEventsPath(); got != filepath.Join(tempDir, "match_events.jsonl") {
		t.Fatalf("unexpected match events path: %q", got)
	}
	if got := s.richFightPath(); got != filepath.Join(tempDir, "fight_history.jsonl") {
		t.Fatalf("unexpected rich fight path: %q", got)
	}
	if got := s.commandResultsPath(); got != filepath.Join(tempDir, "nested", "command_results.jsonl") {
		t.Fatalf("unexpected command results path: %q", got)
	}
	if got := s.liveStatusPath(); got != statusPath {
		t.Fatalf("unexpected live status path: %q", got)
	}
}

func splitNonEmptyLines(raw string) []string {
	lines := make([]string, 0)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}
