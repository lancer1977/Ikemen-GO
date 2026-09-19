package main

import (
	"math"
	"testing"
)

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
	// StageProps currently holds a single bool that newStageProps sets to false,
	// so a correctly initialised value is indistinguishable from the zero value.
	// Asserting equality with newStageProps() rather than with a hand-written
	// zero value keeps this honest: if a field with a non-zero default is added,
	// this starts checking something real instead of silently staying true.
	if s.stageprops != newStageProps() {
		t.Fatalf("newStage stageprops = %#v, want %#v", s.stageprops, newStageProps())
	}
}
