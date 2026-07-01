package main

import "testing"

func TestAnimTextSnd_ResetRestoresCountersAndPaletteTiming(t *testing.T) {
	ats := newAnimTextSnd(&Sff{}, 0)
	ats.cnt = 5
	ats.text.pfxinit = 42
	ats.text.palfx.time = 7

	ats.Reset()

	if ats.cnt != 0 {
		t.Fatalf("Reset cnt = %d, want 0", ats.cnt)
	}
	if ats.text.palfx.time != 42 {
		t.Fatalf("Reset palfx time = %d, want 42", ats.text.palfx.time)
	}
}

func TestAnimTextSnd_ActionAdvancesCounterAndTextStep(t *testing.T) {
	ats := newAnimTextSnd(&Sff{}, 0)
	ats.cnt = 3
	ats.text.pfxinit = 11
	ats.text.palfx.time = 11
	ats.text.forcecolor = true

	ats.Action()

	if ats.cnt != 4 {
		t.Fatalf("Action cnt = %d, want 4", ats.cnt)
	}
	if ats.text.palfx.time != 11 {
		t.Fatalf("Action should not force palette timing change, got %d", ats.text.palfx.time)
	}
}
