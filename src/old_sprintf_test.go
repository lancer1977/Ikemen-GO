package main

import "testing"

func TestOldSprintf_RewritesLegacyLengthModifiers(t *testing.T) {
	got := OldSprintf("x=%li y=%hu z=%Lf", int32(7), uint16(8), float64(9.5))
	if got != "x=7 y=8 z=9.500000" {
		t.Fatalf("OldSprintf() = %q, want %q", got, "x=7 y=8 z=9.500000")
	}
}

func TestOldSprintf_IgnoresEscapedPercentSigns(t *testing.T) {
	got := OldSprintf("rate=100%% value=%d", 42)
	if got != "rate=100% value=42" {
		t.Fatalf("OldSprintf() = %q, want %q", got, "rate=100% value=42")
	}
}
