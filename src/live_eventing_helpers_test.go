package main

import "testing"

func TestFirstNonEmpty_ReturnsFirstTrimmedNonBlankValue(t *testing.T) {
	if got := firstNonEmpty("   ", "\t", "Ryu", "Ken"); got != "Ryu" {
		t.Fatalf("unexpected firstNonEmpty result: %q", got)
	}
	if got := firstNonEmpty("", "   "); got != "" {
		t.Fatalf("expected all-blank input to return blank, got %q", got)
	}
	if got := firstNonEmpty("  Ken  "); got != "Ken" {
		t.Fatalf("expected trimming to apply, got %q", got)
	}
}
