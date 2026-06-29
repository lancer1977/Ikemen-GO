package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// GameLiveSnapshot captures a live in-progress match state.
// It reuses the same fighter/round structures as the final stats contract
// and adds the current round state so consumers can poll during a match.
type GameLiveSnapshot struct {
	StatsLog          StatsLog   `json:"statsLog"`
	ContinueFlg       bool       `json:"continueFlg"`
	PersistRoundCount int32      `json:"persistRoundCount"`
	MatchOver         bool       `json:"matchOver"`
	FrameCounter      int32      `json:"frameCounter"`
	MatchTime         int32      `json:"matchTime"`
	CurRoundTime      int32      `json:"curRoundTime"`
	CurrentRound      StatsRound `json:"currentRound"`
}

func (s *System) liveMatchSnapshot() GameLiveSnapshot {
	fighters := s.liveMatchFighters()
	score := [2]int32{
		int32(s.fightScreen.scores[0].scorePoints),
		int32(s.fightScreen.scores[1].scorePoints),
	}

	return GameLiveSnapshot{
		StatsLog:          s.statsLog,
		ContinueFlg:       s.continueFlg,
		PersistRoundCount: s.persistRoundCount,
		MatchOver:         s.matchOver(),
		FrameCounter:      s.frameCounter,
		MatchTime:         s.matchTime,
		CurRoundTime:      s.curRoundTime,
		CurrentRound: StatsRound{
			Index:    s.round,
			Timer:    s.curRoundTime,
			Score:    score,
			Fighters: fighters,
		},
	}
}

func (s *System) liveMatchFighters() [2][]StatsFighterState {
	var fighters [2][]StatsFighterState
	for _, p := range s.chars {
		if len(p) == 0 || p[0].teamside == -1 {
			continue
		}
		side := int(p[0].teamside)
		if side < 0 || side >= len(fighters) {
			continue
		}
		fighters[side] = append(fighters[side], StatsFighterState{
			Name:       p[0].name,
			ID:         p[0].id,
			MemberNo:   int(p[0].memberNo),
			SelectNo:   int(p[0].selectNo),
			AILevel:    p[0].getAILevel(),
			PalNo:      p[0].gi().palno,
			Life:       p[0].life,
			LifeMax:    p[0].lifeMax,
			Power:      p[0].power,
			PowerMax:   p[0].powerMax,
			WinQuote:   p[0].winquote,
			Win:        p[0].win(),
			WinKO:      p[0].winKO(),
			WinTime:    p[0].winTime(),
			WinPerfect: p[0].winPerfect(),
			WinSpecial: p[0].winType(WT_Special),
			WinHyper:   p[0].winType(WT_Hyper),
			DrawGame:   p[0].drawgame(),
			KO:         p[0].scf(SCF_ko),
			OverKO:     p[0].scf(SCF_over_ko),
		})
	}
	return fighters
}

func (s *System) maybeWriteLiveSnapshot() {
	s.maybeProcessLiveCommandInbox()

	path, ok := s.cmdFlags["-livedatafile"]
	if !ok || path == "" {
		s.writeLiveStatus()
		return
	}
	if !s.middleOfMatch() && !s.matchOver() {
		s.writeLiveStatus()
		return
	}
	if s.middleOfMatch() && s.frameCounter%2 != 0 {
		s.writeLiveStatus()
		return
	}

	data, err := json.Marshal(s.liveMatchSnapshot())
	if err != nil {
		LogMessage("live snapshot marshal failed: %v", err)
		s.writeLiveStatus()
		return
	}
	if err := writeAtomicFile(path, data); err != nil {
		LogMessage("live snapshot write failed: %v", err)
	}
	s.writeLiveStatus()
}

func writeAtomicFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	tmp, err := os.CreateTemp(dir, ".live-snapshot-*.json")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
	}()
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return err
	}
	return nil
}
