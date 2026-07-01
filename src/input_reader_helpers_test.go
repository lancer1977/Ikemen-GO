package main

import "testing"

func TestInputReaderResetClearsState(t *testing.T) {
	ir := &InputReader{
		SocdAllow:          [4]bool{true, true, true, true},
		SocdFirst:          [4]bool{true, true, true, true},
		ButtonAssistBuffer: [9]bool{true, true, true, true, true, true, true, true, true},
	}
	ir.Reset()
	if ir.SocdAllow != [4]bool{} || ir.SocdFirst != [4]bool{} || ir.ButtonAssistBuffer != [9]bool{} {
		t.Fatalf("Reset() = %#v", ir)
	}
}

func TestInputReaderSocdResolutionHonorsPriorityAndAllowMask(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.cfg.Input.SOCDResolution = 1
	ir := &InputReader{SocdAllow: [4]bool{true, true, true, true}}

	U, D, B, F := ir.SocdResolution(true, true, true, true)
	if U || D || !B || !F {
		t.Fatalf("SOCD first-input resolution should clear U/D while preserving B/F maskable state, got %v %v %v %v", U, D, B, F)
	}
	if ir.SocdFirst[0] || !ir.SocdFirst[1] {
		t.Fatalf("SOCD first-input tracking should prefer the last held direction, got %#v", ir.SocdFirst)
	}

	sys.cfg.Input.SOCDResolution = 0
	ir = &InputReader{SocdAllow: [4]bool{true, true, true, true}}
	U, D, B, F = ir.SocdResolution(true, true, false, false)
	if U || D || B || F {
		t.Fatalf("SOCD neutral resolution should clear opposing directions, got %v %v %v %v", U, D, B, F)
	}
}

func TestButtonAssistCheckResetsOutOfMatchAndBuffersInMatch(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	ir := &InputReader{ButtonAssistBuffer: [9]bool{true, true, true, true, true, true, true, true, true}}

	sys.paused = true
	current := [9]bool{true, false, true, false, true, false, true, false, true}
	if got := ir.ButtonAssistCheck(current); got != current {
		t.Fatalf("ButtonAssistCheck() paused path = %#v, want %#v", got, current)
	}
	if ir.ButtonAssistBuffer != [9]bool{} {
		t.Fatalf("ButtonAssistCheck() paused path should clear buffer, got %#v", ir.ButtonAssistBuffer)
	}

	sys.paused = false
	sys.match = 1
	sys.intro = 0
	sys.finishType = FT_NotYet
	sys.fightScreen.round.ctrl_time = 0
	sys.fightScreen.round.over_waittime = 0
	sys.fightScreen.round.over_hittime = 0
	sys.winposetime = 1
	sys.curRoundTime = 1
	sys.timerStart = 1
	ir.ButtonAssistBuffer = [9]bool{true, false, false, false, false, false, false, false, false}
	current = [9]bool{true, true, false, false, false, false, false, false, false}
	got := ir.ButtonAssistCheck(current)
	if !got[0] || got[1] || got[2] {
		t.Fatalf("ButtonAssistCheck() in-match path = %#v", got)
	}
	if ir.ButtonAssistBuffer != current {
		t.Fatalf("ButtonAssistCheck() should store current frame, got %#v", ir.ButtonAssistBuffer)
	}
}
