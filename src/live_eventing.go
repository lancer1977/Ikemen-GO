package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	liveCombatEventSchema   = "live-lancero/combat-event/v1"
	liveCommandInboxSchema  = "live-lancero/command-inbox/v1"
	liveCommandResultSchema = "live-lancero/command-result/v1"
)

type LiveCombatEvent struct {
	Schema        string `json:"schema"`
	EventID       string `json:"eventId"`
	Match         int32  `json:"match"`
	Round         int32  `json:"round"`
	Event         string `json:"event"`
	AttackerSide  int32  `json:"attackerSide,omitempty"`
	AttackerKey   string `json:"attackerKey,omitempty"`
	DefenderSide  int32  `json:"defenderSide,omitempty"`
	DefenderKey   string `json:"defenderKey,omitempty"`
	Damage        int32  `json:"damage,omitempty"`
	LifeBefore    int32  `json:"lifeBefore,omitempty"`
	LifeAfter     int32  `json:"lifeAfter,omitempty"`
	ComboCount    int32  `json:"comboCount,omitempty"`
	ThresholdKind string `json:"thresholdKind,omitempty"`
	Side          int32  `json:"side,omitempty"`
	RosterKey     string `json:"rosterKey,omitempty"`
	Value         int32  `json:"value,omitempty"`
	TextHint      string `json:"textHint,omitempty"`
	TimestampUTC  string `json:"timestampUtc"`
}

type LiveCommandInboxFile struct {
	Schema   string               `json:"schema"`
	Commands []LiveCommandRequest `json:"commands"`
}

type LiveCommandRequest struct {
	ID               string `json:"id"`
	Command          string `json:"command,omitempty"`
	CommandKind      string `json:"commandKind,omitempty"`
	SafeTimingPolicy string `json:"safeTimingPolicy,omitempty"`
	Side             int32  `json:"side,omitempty"`
	Amount           int32  `json:"amount,omitempty"`
}

type LiveCommandResult struct {
	Schema         string `json:"schema"`
	EventID        string `json:"eventId"`
	CommandID      string `json:"commandId"`
	CommandKind    string `json:"commandKind"`
	Status         string `json:"status"`
	Reason         string `json:"reason,omitempty"`
	AppliedMatchID string `json:"appliedMatchId,omitempty"`
	AppliedRound   int32  `json:"appliedRound,omitempty"`
	AppliedFrame   int32  `json:"appliedFrame,omitempty"`
	TimestampUTC   string `json:"timestampUtc"`
}

func (s *System) combatEventsPath() string {
	if path := strings.TrimSpace(s.cmdFlags["-combateventsfile"]); path != "" {
		return path
	}
	if livePath := strings.TrimSpace(s.cmdFlags["-livedatafile"]); livePath != "" {
		return filepath.Join(filepath.Dir(livePath), "combat_events.jsonl")
	}
	return ""
}

func (s *System) commandInboxPath() string {
	return strings.TrimSpace(s.cmdFlags["-commandinboxfile"])
}

func (s *System) commandResultsPath() string {
	if path := strings.TrimSpace(s.cmdFlags["-commandresultsfile"]); path != "" {
		return path
	}
	if inboxPath := s.commandInboxPath(); inboxPath != "" {
		return filepath.Join(filepath.Dir(inboxPath), "command_results.jsonl")
	}
	return ""
}

func appendJSONLine(path string, value any) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, err := file.Write(append(data, '\n')); err != nil {
		return err
	}
	return file.Sync()
}

func parseLiveCommandInbox(data []byte) ([]LiveCommandRequest, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, nil
	}
	var inbox LiveCommandInboxFile
	if err := json.Unmarshal(data, &inbox); err != nil {
		return nil, err
	}
	if inbox.Schema != "" && inbox.Schema != liveCommandInboxSchema {
		return nil, fmt.Errorf("unsupported command inbox schema: %s", inbox.Schema)
	}
	return inbox.Commands, nil
}

func normalizeLiveCommandKind(request LiveCommandRequest) string {
	kind := request.CommandKind
	if strings.TrimSpace(kind) == "" {
		kind = request.Command
	}
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(kind), "_", "-"))
}

func liveTimestampUTC() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func liveSideFromChar(c *Char) int32 {
	if c == nil || c.teamside < 0 || c.teamside > 1 {
		return 0
	}
	return int32(c.teamside + 1)
}

func liveSideFromPlayerNo(playerNo int32) int32 {
	if playerNo < 0 {
		return 0
	}
	return int32(playerNo&1) + 1
}

func liveRosterKey(c *Char) string {
	if c == nil {
		return ""
	}
	return c.name
}

func (s *System) recordCombatDamage(defender *Char, attackerPlayerNo int32, requestedDamage int32, lifeBefore int32, lifeAfter int32) {
	path := s.combatEventsPath()
	if path == "" || defender == nil {
		return
	}
	damage := lifeBefore - lifeAfter
	if damage <= 0 {
		return
	}
	attacker := s.playerByNo(attackerPlayerNo)
	attackerSide := liveSideFromPlayerNo(attackerPlayerNo)
	comboCount := int32(0)
	if attackerSide >= 1 && attackerSide <= 2 && s.fightScreen.combos[attackerSide-1] != nil {
		comboCount = s.fightScreen.combos[attackerSide-1].trueHits
	}
	event := LiveCombatEvent{
		Schema:       liveCombatEventSchema,
		EventID:      fmt.Sprintf("match-%d-r%d-f%d-damage-p%d", s.match, s.round, s.frameCounter, defender.playerNo+1),
		Match:        s.match,
		Round:        s.round,
		Event:        "damage-dealt",
		AttackerSide: attackerSide,
		AttackerKey:  liveRosterKey(attacker),
		DefenderSide: liveSideFromChar(defender),
		DefenderKey:  liveRosterKey(defender),
		Damage:       damage,
		LifeBefore:   lifeBefore,
		LifeAfter:    lifeAfter,
		ComboCount:   comboCount,
		TimestampUTC: liveTimestampUTC(),
	}
	if err := appendJSONLine(path, event); err != nil {
		LogMessage("live combat event write failed: %v", err)
	}
	s.recordCombatThresholds(path, event, requestedDamage)
}

func (s *System) recordCombatThresholds(path string, damageEvent LiveCombatEvent, requestedDamage int32) {
	if damageEvent.Damage >= 150 {
		s.writeCombatFact(path, damageEvent, "threshold-crossed", "damage-threshold", damageEvent.AttackerSide, damageEvent.AttackerKey, damageEvent.Damage, "NICE HIT!")
	}
	if damageEvent.LifeAfter > 0 && damageEvent.LifeBefore > damageEvent.LifeAfter && damageEvent.LifeAfter*4 <= damageEvent.LifeBefore {
		s.writeCombatFact(path, damageEvent, "ko-near", "low-health", damageEvent.DefenderSide, damageEvent.DefenderKey, damageEvent.LifeAfter, "LOW HEALTH")
	}
	if damageEvent.LifeAfter <= 0 {
		s.writeCombatFact(path, damageEvent, "ko-confirmed", "ko", damageEvent.DefenderSide, damageEvent.DefenderKey, requestedDamage, "KO")
	}
}

func (s *System) writeCombatFact(path string, source LiveCombatEvent, eventName string, thresholdKind string, side int32, rosterKey string, value int32, textHint string) {
	event := LiveCombatEvent{
		Schema:        liveCombatEventSchema,
		EventID:       fmt.Sprintf("match-%d-r%d-f%d-%s-p%d", s.match, s.round, s.frameCounter, eventName, side),
		Match:         source.Match,
		Round:         source.Round,
		Event:         eventName,
		ThresholdKind: thresholdKind,
		Side:          side,
		RosterKey:     rosterKey,
		Value:         value,
		TextHint:      textHint,
		TimestampUTC:  liveTimestampUTC(),
	}
	if err := appendJSONLine(path, event); err != nil {
		LogMessage("live combat fact write failed: %v", err)
	}
}

func (s *System) maybeProcessLiveCommandInbox() {
	inboxPath := s.commandInboxPath()
	resultsPath := s.commandResultsPath()
	if inboxPath == "" || resultsPath == "" {
		return
	}
	if !s.middleOfMatch() {
		return
	}
	if s.lastLiveCommandPollFrame == s.frameCounter || s.frameCounter%6 != 0 {
		return
	}
	s.lastLiveCommandPollFrame = s.frameCounter
	data, err := os.ReadFile(inboxPath)
	if err != nil {
		if !os.IsNotExist(err) {
			LogMessage("live command inbox read failed: %v", err)
		}
		return
	}
	commands, err := parseLiveCommandInbox(data)
	if err != nil {
		LogMessage("live command inbox parse failed: %v", err)
		return
	}
	if s.liveProcessedCommandIDs == nil {
		s.liveProcessedCommandIDs = make(map[string]bool)
	}
	for _, command := range commands {
		if command.ID == "" || s.liveProcessedCommandIDs[command.ID] {
			continue
		}
		s.liveProcessedCommandIDs[command.ID] = true
		result := s.applyLiveCommand(command)
		if err := appendJSONLine(resultsPath, result); err != nil {
			LogMessage("live command result write failed: %v", err)
		}
	}
}

func (s *System) applyLiveCommand(command LiveCommandRequest) LiveCommandResult {
	kind := normalizeLiveCommandKind(command)
	result := LiveCommandResult{
		Schema:         liveCommandResultSchema,
		EventID:        fmt.Sprintf("%s-%d", command.ID, s.frameCounter),
		CommandID:      command.ID,
		CommandKind:    kind,
		Status:         "failed",
		AppliedMatchID: fmt.Sprintf("match-%d", s.match),
		AppliedRound:   s.round,
		AppliedFrame:   s.frameCounter,
		TimestampUTC:   liveTimestampUTC(),
	}
	switch kind {
	case "power-adjust":
		if command.Side != 1 && command.Side != 2 {
			result.Reason = "invalid-side"
			return result
		}
		if command.Amount < -3000 || command.Amount > 3000 {
			result.Reason = "invalid-amount"
			return result
		}
		if !s.middleOfMatch() {
			result.Reason = "not-in-fight-phase"
			return result
		}
		root := s.teamRoot(command.Side - 1)
		if root == nil {
			result.Reason = "missing-team-root"
			return result
		}
		root.powerAdd(command.Amount)
		result.Status = "applied"
		return result
	case "players-swap":
		result.Status = "queued"
		result.Reason = "next-match-hook-pending"
		return result
	default:
		result.CommandKind = "unknown"
		result.Reason = "unsupported-command"
		return result
	}
}

func (s *System) teamRoot(teamSide int32) *Char {
	if teamSide < 0 || int(teamSide) >= len(s.chars) || len(s.chars[teamSide]) == 0 {
		return nil
	}
	return s.chars[teamSide][0]
}

func (s *System) playerByNo(playerNo int32) *Char {
	for _, side := range s.chars {
		for _, c := range side {
			if c != nil && c.playerNo == playerNo {
				return c
			}
		}
	}
	return nil
}
