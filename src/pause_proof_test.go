package main

import (
	"testing"
)

func TestParsePauseProofSeconds_ClampsAndRejectsInvalidInput(t *testing.T) {
	t.Setenv("IKEMEN_PAUSE_PROOF_SECONDS", "")
	if got := parsePauseProofSeconds(); got != 0 {
		t.Fatalf("expected empty env to disable pause proof, got %d", got)
	}

	t.Setenv("IKEMEN_PAUSE_PROOF_SECONDS", "7")
	if got := parsePauseProofSeconds(); got != 7 {
		t.Fatalf("expected explicit value to be parsed, got %d", got)
	}

	t.Setenv("IKEMEN_PAUSE_PROOF_SECONDS", "45")
	if got := parsePauseProofSeconds(); got != 30 {
		t.Fatalf("expected pause proof seconds to clamp to 30, got %d", got)
	}

	t.Setenv("IKEMEN_PAUSE_PROOF_SECONDS", "not-a-number")
	if got := parsePauseProofSeconds(); got != 0 {
		t.Fatalf("expected invalid input to disable pause proof, got %d", got)
	}
}
