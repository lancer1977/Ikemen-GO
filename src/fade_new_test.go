package main

import "testing"

func TestNewFade(t *testing.T) {
	t.Parallel()

	fa := newFade()
	if fa == nil {
		t.Fatal("newFade returned nil")
	}
	if fa.time != 30 {
		t.Fatalf("time = %d, want 30", fa.time)
	}
	if fa.snd != [2]int32{-1, 0} {
		t.Fatalf("snd = %#v, want [-1 0]", fa.snd)
	}
	if fa.active || fa.totalTime != 0 || fa.overlayDelay != 0 || fa.isFadeIn {
		t.Fatalf("unexpected fade defaults: %#v", fa)
	}
}
