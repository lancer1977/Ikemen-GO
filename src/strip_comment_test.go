package main

import "testing"

func TestStripComment_RemovesInlineCommentAndTrimsSpaceAgain(t *testing.T) {
	if got := StripComment("value ; trailing comment  "); got != "value" {
		t.Fatalf("StripComment() = %q, want %q", got, "value")
	}
	if got := StripComment("  keep me  "); got != "keep me" {
		t.Fatalf("StripComment() = %q, want %q", got, "keep me")
	}
	if got := StripComment("x;y"); got != "x" {
		t.Fatalf("StripComment() = %q, want %q", got, "x")
	}
}
