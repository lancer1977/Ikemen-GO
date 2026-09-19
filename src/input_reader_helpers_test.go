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

func TestInputReaderSocdResolutionHonorsConfiguredMethod(t *testing.T) {
	oldSys := sys
	t.Cleanup(func() { sys = oldSys })
	sys = System{}

	// Method 1 is last-direction priority. With all four directions held on the
	// first frame, D and B are recorded as the first-held ones, so the opposing
	// U and F survive.
	sys.cfg.Input.SOCDResolution = 1
	ir := &InputReader{SocdAllow: [4]bool{true, true, true, true}}
	U, D, B, F := ir.SocdResolution(true, true, true, true)
	if !U || D || B || !F {
		t.Fatalf("last-direction priority = %v %v %v %v, want true false false true", U, D, B, F)
	}
	if want := [4]bool{false, true, true, false}; ir.SocdFirst != want {
		t.Fatalf("SocdFirst = %#v, want %#v", ir.SocdFirst, want)
	}

	// Method 0 allows both directions: an up+down conflict is left untouched.
	sys.cfg.Input.SOCDResolution = 0
	ir = &InputReader{SocdAllow: [4]bool{true, true, true, true}}
	U, D, B, F = ir.SocdResolution(true, true, false, false)
	if !U || !D || B || F {
		t.Fatalf("allow-both resolution = %v %v %v %v, want true true false false", U, D, B, F)
	}

	// Any unrecognised method falls through to neutral resolution, which denies
	// both sides of every conflicting pair.
	sys.cfg.Input.SOCDResolution = 4
	ir = &InputReader{SocdAllow: [4]bool{true, true, true, true}}
	U, D, B, F = ir.SocdResolution(true, true, true, true)
	if U || D || B || F {
		t.Fatalf("neutral resolution = %v %v %v %v, want all false", U, D, B, F)
	}
}

func TestButtonAssistCheckResetsOutOfMatchAndBuffersInMatch(t *testing.T) {
	oldSys := sys
	t.Cleanup(func() { sys = oldSys })
	sys = System{}

	ir := &InputReader{ButtonAssistBuffer: [9]bool{true, true, true, true, true, true, true, true, true}}

	// Paused: the current frame passes through untouched and the buffer is
	// dropped so held buttons do not leak back in on resume.
	sys.paused = true
	current := [9]bool{true, false, true, false, true, false, true, false, true}
	if got := ir.ButtonAssistCheck(current); got != current {
		t.Fatalf("ButtonAssistCheck() paused path = %#v, want %#v", got, current)
	}
	if ir.ButtonAssistBuffer != [9]bool{} {
		t.Fatalf("ButtonAssistCheck() paused path should clear buffer, got %#v", ir.ButtonAssistBuffer)
	}

	// middleOfMatch() is the in-match gate: it needs a non-zero matchTime with
	// neither the fight loop ended nor the post-match flag set.
	sys.paused = false
	sys.matchTime = 1
	sys.fightLoopEnd = false
	sys.postMatchFlg = false

	// With a button held last frame, this frame's presses are merged in too.
	ir.ButtonAssistBuffer = [9]bool{true, false, false, false, false, false, false, false, false}
	current = [9]bool{true, true, false, false, false, false, false, false, false}
	got := ir.ButtonAssistCheck(current)
	if want := [9]bool{true, true, false, false, false, false, false, false, false}; got != want {
		t.Fatalf("ButtonAssistCheck() in-match path = %#v, want %#v", got, want)
	}
	if ir.ButtonAssistBuffer != current {
		t.Fatalf("ButtonAssistCheck() should store current frame, got %#v", ir.ButtonAssistBuffer)
	}

	// With an empty buffer only the previous frame counts, so nothing registers.
	ir.ButtonAssistBuffer = [9]bool{}
	got = ir.ButtonAssistCheck([9]bool{true, true, false, false, false, false, false, false, false})
	if got != ([9]bool{}) {
		t.Fatalf("ButtonAssistCheck() with empty buffer = %#v, want all false", got)
	}
}
