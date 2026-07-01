package main

import "testing"

func TestNewPowerBar(t *testing.T) {
	t.Parallel()

	pb := newPowerBar()
	if pb == nil {
		t.Fatal("newPowerBar returned nil")
	}
	if pb.front == nil || pb.bg0 == nil || pb.counter == nil || pb.value == nil {
		t.Fatalf("newPowerBar should allocate maps: %#v", pb)
	}
	if pb.counter_rounding != 1000 || pb.value_rounding != 1 || pb.levelmax_snd != [2]int32{-1, -1} {
		t.Fatalf("unexpected numeric defaults: %#v", pb)
	}
	for i, snd := range pb.level_snd {
		if snd != [2]int32{-1, -1} {
			t.Fatalf("level_snd[%d] = %#v, want [-1 -1]", i, snd)
		}
	}
}
