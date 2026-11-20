package main

import (
	"testing"
)

// A format string ending in a bare "%" should render that percent as a
// literal character, not as fmt.Sprintf's dangling-verb error diagnostic.
// When OldSprintf's scanner encounters a truncated format spec, it now
// escapes the trailing percent by converting it to %% before passing to
// fmt.Sprintf, so the output is clean and safe for game text.
func TestOldSprintf_TruncatedFormatLeaksSprintfErrorText(t *testing.T) {
	got := OldSprintf("value=%")
	if got != "value=%" {
		t.Fatalf("OldSprintf(truncated) = %q, want %q", got, "value=%")
	}

	// Same shape with a real verb in front: the argument substitutes correctly,
	// and the trailing percent renders as a literal character.
	got = OldSprintf("x=%d%", 42)
	if got != "x=42%" {
		t.Fatalf("OldSprintf(trailing-percent) = %q, want %q", got, "x=42%")
	}

	// The escaped form is the supported way to emit a literal percent, and it
	// is unaffected.
	if got := OldSprintf("rate=100%%"); got != "rate=100%" {
		t.Fatalf("OldSprintf(escaped) = %q, want %q", got, "rate=100%")
	}
}
