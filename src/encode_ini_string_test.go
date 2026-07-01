package main

import "testing"

func TestEncodeIniStringQuotesWhenNeeded(t *testing.T) {
	if got := encodeIniString("plain"); got != "plain" {
		t.Fatalf("encodeIniString(plain) = %q, want plain", got)
	}
	if got := encodeIniString("a,b"); got != `"a,b"` {
		t.Fatalf("encodeIniString(comma) = %q, want quoted", got)
	}
	if got := encodeIniString(" spaced"); got != `" spaced"` {
		t.Fatalf("encodeIniString(spaced) = %q, want quoted", got)
	}
}
