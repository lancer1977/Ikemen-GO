package main

import "testing"

func TestNewAfterImage(t *testing.T) {
	t.Parallel()

	ai := newAfterImage()
	if ai == nil {
		t.Fatal("newAfterImage returned nil")
	}
	if ai.time != 0 || ai.length != 20 || ai.timegap != 1 || ai.framegap != 4 {
		t.Fatalf("unexpected afterimage defaults: %#v", ai)
	}
	if ai.trans != TT_default || ai.alpha != [2]int32{-1, 0} || !ai.ignorehitpause {
		t.Fatalf("unexpected afterimage defaults: %#v", ai)
	}
	if len(ai.palfx) != 1 || len(ai.imgs) != 1 {
		t.Fatalf("unexpected slice lengths: palfx=%d imgs=%d", len(ai.palfx), len(ai.imgs))
	}
	if ai.palfx[0] == nil || !ai.palfx[0].enable || !ai.palfx[0].allowNeg {
		t.Fatalf("unexpected palfx defaults: %#v", ai.palfx[0])
	}
}
