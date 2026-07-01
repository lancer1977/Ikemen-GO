package main

import "testing"

func TestNormalizeNewlines_ConvertsCRLFAndCRToLFLegacy(t *testing.T) {
	got := NormalizeNewlines("a\r\nb\rc\n")
	if got != "a\nb\nc\n" {
		t.Fatalf("NormalizeNewlines() = %q, want %q", got, "a\nb\nc\n")
	}
}

func TestNormalizeNewlines_LeavesLFOnlyInputUnchanged(t *testing.T) {
	got := NormalizeNewlines("line1\nline2\n")
	if got != "line1\nline2\n" {
		t.Fatalf("NormalizeNewlines() = %q, want %q", got, "line1\nline2\n")
	}
}
