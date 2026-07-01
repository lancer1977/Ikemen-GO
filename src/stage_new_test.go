package main

import "testing"

func TestNewStage(t *testing.T) {
	s := newStage("stages/test.def")
	if s == nil {
		t.Fatal("expected newStage to allocate")
	}
	if s.def != "stages/test.def" {
		t.Fatalf("newStage def = %q", s.def)
	}
	if s.leftbound != -1000 || s.rightbound != 1000 || s.screenleft != 15 || s.screenright != 15 {
		t.Fatalf("newStage bounds = %#v %#v %#v %#v", s.leftbound, s.rightbound, s.screenleft, s.screenright)
	}
	if !s.autoturn || !s.resetbg || s.localscl != 1 || s.bgmratio != 0.3 || s.partnerspacing != 25 {
		t.Fatalf("newStage defaults = %#v", s)
	}
	if !math.IsNaN(float64(s.scale[0])) || !math.IsNaN(float64(s.scale[1])) {
		t.Fatalf("newStage scale should be NaN, got %#v", s.scale)
	}
	if s.stageprops == (StageProps{}) {
		t.Fatalf("newStage stageprops should be initialized")
	}
}
