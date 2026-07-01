package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRecordCombatThresholds_WritesExpectedThresholdFacts(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "combat_events.jsonl")
	s := &System{
		SystemStateVars: SystemStateVars{
			match: 12,
			round: 2,
		},
	}

	damageEvent := LiveCombatEvent{
		Match:        12,
		Round:        2,
		AttackerSide: 1,
		AttackerKey:  "ryu",
		DefenderSide: 2,
		DefenderKey:  "ken",
		Damage:       180,
		LifeBefore:   200,
		LifeAfter:    40,
	}
	s.recordCombatThresholds(path, damageEvent, 180)

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read combat threshold facts: %v", err)
	}
	lines := splitNonEmptyLines(string(raw))
	if len(lines) != 2 {
		t.Fatalf("expected two threshold facts, got %d (%q)", len(lines), string(raw))
	}

	var first, second LiveCombatEvent
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatalf("unmarshal first threshold fact: %v", err)
	}
	if err := json.Unmarshal([]byte(lines[1]), &second); err != nil {
		t.Fatalf("unmarshal second threshold fact: %v", err)
	}
	if first.Event != "threshold-crossed" || first.ThresholdKind != "damage-threshold" || first.TextHint != "NICE HIT!" {
		t.Fatalf("unexpected first threshold fact: %#v", first)
	}
	if second.Event != "ko-near" || second.ThresholdKind != "low-health" || second.TextHint != "LOW HEALTH" {
		t.Fatalf("unexpected second threshold fact: %#v", second)
	}
}

func TestRecordCombatThresholds_WritesKoConfirmationFact(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "combat_events.jsonl")
	s := &System{}

	damageEvent := LiveCombatEvent{
		Match:        12,
		Round:        2,
		DefenderSide: 2,
		DefenderKey:  "ken",
		LifeBefore:   50,
		LifeAfter:    0,
	}
	s.recordCombatThresholds(path, damageEvent, 999)

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read combat threshold facts: %v", err)
	}
	lines := splitNonEmptyLines(string(raw))
	if len(lines) != 1 {
		t.Fatalf("expected one ko fact, got %d (%q)", len(lines), string(raw))
	}

	var fact LiveCombatEvent
	if err := json.Unmarshal([]byte(lines[0]), &fact); err != nil {
		t.Fatalf("unmarshal ko fact: %v", err)
	}
	if fact.Event != "ko-confirmed" || fact.ThresholdKind != "ko" || fact.TextHint != "KO" || fact.Value != 999 {
		t.Fatalf("unexpected ko fact: %#v", fact)
	}
}
