package main

import "testing"

func TestNormalizeNewlines_ConvertsCRLFAndCRToLFAgain(t *testing.T) {
	if got := NormalizeNewlines("a\r\nb\rc\n"); got != "a\nb\nc\n" {
		t.Fatalf("NormalizeNewlines() = %q, want %q", got, "a\nb\nc\n")
	}
	if got := NormalizeNewlines("line1\nline2\n"); got != "line1\nline2\n" {
		t.Fatalf("NormalizeNewlines(no-op) = %q, want %q", got, "line1\nline2\n")
	}
}
