package main

import (
	"strings"
	"testing"
)

// A format string ending in a bare "%" is left in place: OldSprintf's scanner
// breaks out of the loop when the percent has no verb after it, and the
// unchanged format is then handed to fmt.Sprintf, which renders the dangling
// verb as its own error text rather than a literal percent sign.
//
// The test name this replaced ("LeavesTruncatedFormatPrefixIntact") described
// the behaviour someone intended, not the behaviour that exists. Pinning the
// real output here keeps the gap visible; see the engine issue tracking it.
func TestOldSprintf_TruncatedFormatLeaksSprintfErrorText(t *testing.T) {
	got := OldSprintf("value=%")
	if !strings.HasPrefix(got, "value=") {
		t.Fatalf("OldSprintf(truncated) = %q, want it to keep the literal prefix", got)
	}
	if !strings.Contains(got, "NOVERB") {
		t.Fatalf("OldSprintf(truncated) = %q, want fmt's dangling-verb error text; "+
			"if this now returns %q the engine bug is fixed and this test should "+
			"assert that instead", got, "value=%")
	}

	// Same shape with a real verb in front: the argument still substitutes
	// correctly, and only the trailing percent degrades.
	got = OldSprintf("x=%d%", 42)
	if !strings.HasPrefix(got, "x=42") {
		t.Fatalf("OldSprintf(trailing-percent) = %q, want the verb to substitute", got)
	}
	if !strings.Contains(got, "NOVERB") {
		t.Fatalf("OldSprintf(trailing-percent) = %q, want fmt's dangling-verb error text", got)
	}

	// The escaped form is the supported way to emit a literal percent, and it
	// is unaffected.
	if got := OldSprintf("rate=100%%"); got != "rate=100%" {
		t.Fatalf("OldSprintf(escaped) = %q, want %q", got, "rate=100%")
	}
}
