package main

import "testing"

func TestNewFightScreenRound(t *testing.T) {
	t.Parallel()

	snd := &Snd{}
	ro := newFightScreenRound(snd)
	if ro == nil {
		t.Fatal("newFightScreenRound returned nil")
	}
	if ro.snd != snd {
		t.Fatal("newFightScreenRound should retain provided sound pointer")
	}
	if ro.start_waittime != 30 || ro.ctrl_time != 30 || ro.slow_time != 60 || ro.slow_fadetime != 45 {
		t.Fatalf("unexpected round timing defaults: %#v", ro)
	}
	if ro.slow_speed != 0.25 || ro.over_waittime != 45 || ro.over_hittime != 10 || ro.over_wintime != 45 {
		t.Fatalf("unexpected round timing defaults: %#v", ro)
	}
	if ro.over_forcewintime != 900 || ro.over_time != 210 || ro.shutter_time != 15 || ro.callfight_time != 60 || ro.clutch_threshold != 10 {
		t.Fatalf("unexpected round timing defaults: %#v", ro)
	}
}
