package main

import "testing"

func TestParseTriggerNumber(t *testing.T) {
	t.Parallel()

	tn, all, ok := parseTriggerNumber("triggerall")
	if !ok || !all || tn != 0 {
		t.Fatalf("triggerall parsed as tn=%d all=%v ok=%v", tn, all, ok)
	}

	tn, all, ok = parseTriggerNumber("trigger12")
	if !ok || all || tn != 12 {
		t.Fatalf("trigger12 parsed as tn=%d all=%v ok=%v", tn, all, ok)
	}

	for _, in := range []string{"trigger0", "trigger65537", "trig12", "triggerx"} {
		if tn, all, ok := parseTriggerNumber(in); ok || all || tn != 0 {
			t.Fatalf("%q unexpectedly parsed as tn=%d all=%v ok=%v", in, tn, all, ok)
		}
	}
}
