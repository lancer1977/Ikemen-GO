package main

import "testing"

func TestOldSprintf_LeavesTruncatedFormatPrefixIntact(t *testing.T) {
	got := OldSprintf("value=%")
	if got != "value=%" {
		t.Fatalf("OldSprintf(truncated) = %q, want %q", got, "value=%")
	}

	got = OldSprintf("x=%d%", 42)
	if got != "x=42%" {
		t.Fatalf("OldSprintf(trailing-percent) = %q, want %q", got, "x=42%")
	}
}
