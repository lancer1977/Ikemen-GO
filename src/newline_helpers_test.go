package main

import "testing"

func TestNormalizeNewlines_ConvertsCRLFAndCRToLF(t *testing.T) {
	if got := NormalizeNewlines("a\r\nb\rc"); got != "a\nb\nc" {
		t.Fatalf("unexpected normalized output: %q", got)
	}
	if got := NormalizeNewlines("already\nnormalized"); got != "already\nnormalized" {
		t.Fatalf("unexpected unchanged output: %q", got)
	}
}
