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

	s.recordCombatDamage(defender, attacker.playerNo, 180, 180, 0)

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
