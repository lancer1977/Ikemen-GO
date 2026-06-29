package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var pauseProofSeconds = parsePauseProofSeconds()

func parsePauseProofSeconds() int32 {
	v := strings.TrimSpace(os.Getenv("IKEMEN_PAUSE_PROOF_SECONDS"))
	if v == "" {
		return 0
	}
	seconds, err := strconv.Atoi(v)
	if err != nil || seconds <= 0 {
		return 0
	}
	if seconds > 30 {
		seconds = 30
	}
	return int32(seconds)
}

func (s *System) maybeStartPauseProof() {
	if pauseProofSeconds <= 0 || s.pauseProofStarted || s.postMatchFlg || !s.fightScreen.active {
		return
	}
	if s.roundState() != 2 || s.matchTime < s.gameLogicSpeed() {
		return
	}

	frames := pauseProofSeconds * s.gameLogicSpeed()
	s.pauseProofStarted = true
	s.pauseProofFrames = frames
	s.pausebg = true
	if s.pausetime < frames {
		s.pausetime = frames
	}
}

func (s *System) drawPauseProofOverlay() {
	if pauseProofSeconds <= 0 || !s.pauseProofStarted || s.pausetime <= 0 || !s.fightScreen.active || s.postMatchFlg {
		return
	}

	logicSpeed := s.gameLogicSpeed()
	remaining := int32(1)
	if logicSpeed > 0 {
		remaining = (s.pausetime + logicSpeed - 1) / logicSpeed
	}
	drawRenderProbeBlockMode("pause-proof", fmt.Sprintf("PAUSE %d", remaining), 96, 88, 128, 48, 255, 255, 64)
}
