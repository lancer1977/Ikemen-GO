package main

import "testing"

func TestNormalizeLiveCommandKind_FallsBackToCommandAndTrimsWhitespace(t *testing.T) {
	got := normalizeLiveCommandKind(LiveCommandRequest{CommandKind: "   ", Command: " power_adjust "})
	if got != "power-adjust" {
		t.Fatalf("normalizeLiveCommandKind fallback = %q, want power-adjust", got)
	}

	got = normalizeLiveCommandKind(LiveCommandRequest{CommandKind: "next_match"})
	if got != "next-match" {
		t.Fatalf("normalizeLiveCommandKind underscore = %q, want next-match", got)
	}
}
