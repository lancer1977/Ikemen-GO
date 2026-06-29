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
	liveMatchEventSchema    = "live-lancero/match-event/v1"
	liveStatusSchema        = "live-lancero/live-status/v1"
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

type LiveStatusSnapshot struct {
	Schema       string               `json:"schema"`
	Mode         string               `json:"mode"`
	Match        int32                `json:"match,omitempty"`
	Round        int32                `json:"round,omitempty"`
	Stage        string               `json:"stage,omitempty"`
	TimestampUTC string               `json:"timestampUtc"`
	P1           *LiveStatusFighter   `json:"p1,omitempty"`
	P2           *LiveStatusFighter   `json:"p2,omitempty"`
}

type LiveStatusFighter struct {
	Key         string `json:"key"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Wins        *int32 `json:"wins,omitempty"`
	Losses      *int32 `json:"losses,omitempty"`
	Tier        *int32 `json:"tier,omitempty"`
}

type LiveMatchEvent struct {
	Schema       string `json:"schema"`
	Line         int32  `json:"line"`
	Match        int32  `json:"match"`
	Round        int32  `json:"round"`
	Event        string `json:"event"`
	WinnerSide   int32  `json:"winnerSide,omitempty"`
	WinnerKey    string `json:"winnerKey,omitempty"`
	LoserKey     string `json:"loserKey,omitempty"`
	TimestampUTC string `json:"timestampUtc"`
}

func (s *System) combatEventsPath() string {
	if path := strings.TrimSpace(s.cmdFlags["-combateventsfile"]); path != "" {
		return path
	}
	if path := strings.TrimSpace(s.cmdFlags["-livedatafile"]); path != "" {
		return filepath.Join(filepath.Dir(path), "combat_events.jsonl")
	}
	if path := strings.TrimSpace(s.cmdFlags["-livestatusfile"]); path != "" {
		return filepath.Join(filepath.Dir(path), "combat_events.jsonl")
	}
	return filepath.Join(s.baseDir, "save", "combat_events.jsonl")
}

func (s *System) liveStatusPath() string {
	if path := strings.TrimSpace(s.cmdFlags["-livestatusfile"]); path != "" {
		return path
	}
	return filepath.Join(s.baseDir, "save", "live_status.json")
}

func (s *System) matchEventsPath() string {
	if path := strings.TrimSpace(s.cmdFlags["-matcheventsfile"]); path != "" {
		return path
	}
	return filepath.Join(s.baseDir, "save", "match_events.jsonl")
}

func (s *System) commandInboxPath() string {
	if path := strings.TrimSpace(s.cmdFlags["-commandinboxfile"]); path != "" {
		return path
	}
	return filepath.Join(s.baseDir, "save", "command_inbox.json")
}

func (s *System) commandResultsPath() string {
	if path := strings.TrimSpace(s.cmdFlags["-commandresultsfile"]); path != "" {
		return path
	}
	if inboxPath := s.commandInboxPath(); inboxPath != "" {
		return filepath.Join(filepath.Dir(inboxPath), "command_results.jsonl")
	}
	return filepath.Join(s.baseDir, "save", "command_results.jsonl")
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

func (s *System) writeLiveStatus() {
	path := s.liveStatusPath()
	if path == "" || (!s.middleOfMatch() && !s.matchOver()) {
		return
	}

	data, err := json.Marshal(s.buildLiveStatusSnapshot())
	if err != nil {
		LogMessage("live status marshal failed: %v", err)
		return
	}
	if err := writeAtomicFile(path, data); err != nil {
		LogMessage("live status write failed: %v", err)
	}
}

func (s *System) buildLiveStatusSnapshot() LiveStatusSnapshot {
	p1, p2 := s.liveStatusFighters()
	return LiveStatusSnapshot{
		Schema:       liveStatusSchema,
		Mode:         "fight",
		Match:        s.match,
		Round:        s.round,
		Stage:        s.liveStatusStageName(),
		TimestampUTC: liveTimestampUTC(),
		P1:           p1,
		P2:           p2,
	}
}

func (s *System) liveStatusFighters() (*LiveStatusFighter, *LiveStatusFighter) {
	return s.liveStatusFighterForSide(0), s.liveStatusFighterForSide(1)
}

func (s *System) liveStatusFighterForSide(side int) *LiveStatusFighter {
	if side < 0 || side >= len(s.chars) || len(s.chars[side]) == 0 || s.chars[side][0] == nil {
		return nil
	}

	c := s.chars[side][0]
	key := strings.TrimSpace(c.name)
	if key == "" {
		key = fmt.Sprintf("player-%d", side+1)
	}
	name := strings.TrimSpace(c.name)
	if name == "" {
		name = key
	}
	displayName := name

	return &LiveStatusFighter{
		Key:         strings.ToLower(key),
		Name:        name,
		DisplayName: displayName,
	}
}

func (s *System) liveStatusStageName() string {
	if s.stage == nil {
		return ""
	}
	if name := strings.TrimSpace(s.stage.displayname); name != "" {
		return name
	}
	if name := strings.TrimSpace(s.stage.name); name != "" {
		return name
	}
	if def := strings.TrimSpace(s.stage.def); def != "" {
		return def
	}
	return ""
}

func (s *System) appendMatchEvent(event string, winnerSide int32, winnerKey string, loserKey string) {
	path := s.matchEventsPath()
	if path == "" {
		return
	}

	s.matchEventLine++
	record := LiveMatchEvent{
		Schema:       liveMatchEventSchema,
		Line:         s.matchEventLine,
		Match:        s.match,
		Round:        s.round,
		Event:        event,
		WinnerSide:   winnerSide,
		WinnerKey:    winnerKey,
		LoserKey:     loserKey,
		TimestampUTC: liveTimestampUTC(),
	}

	if err := appendJSONLine(path, record); err != nil {
		LogMessage("live match event write failed: %v", err)
	}
}

func (s *System) recordRoundStart() {
	s.appendMatchEvent("round_start", 0, "", "")
}

func (s *System) recordRoundOutcome() {
	if s.winTeam < 0 {
		s.appendMatchEvent("round_draw", 0, "", "")
		return
	}

	winnerSide := int32(s.winTeam + 1)
	loserSide := 1 - s.winTeam
	s.appendMatchEvent("round_win", winnerSide, s.liveRosterKeyForSide(s.winTeam), s.liveRosterKeyForSide(loserSide))
}

func (s *System) liveRosterKeyForSide(side int) string {
	if side < 0 || side >= len(s.chars) || len(s.chars[side]) == 0 || s.chars[side][0] == nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(s.chars[side][0].name))
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
			if c != nil && int32(c.playerNo) == playerNo {
				return c
			}
		}
	}
	return nil
}
