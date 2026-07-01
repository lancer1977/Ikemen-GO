package main

import "testing"

func TestLiveFightStageContext_UsesCurrentStageThenSelectedStageFallback(t *testing.T) {
	s := &System{}
	s.sel.selectedStageNo = 2
	s.sel.stagelist = []SelectStage{
		{def: "stages/first.def", name: "first", localcoord: [2]int32{100, 200}},
		{def: "stages/second.def", name: "second", localcoord: [2]int32{300, 400}},
	}

	got := s.liveFightStageContext()
	if got.SelectedStageNo != 2 || got.Def != "stages/second.def" || got.Name != "second" {
		t.Fatalf("unexpected selected-stage fallback: %#v", got)
	}
	if got.LocalCoord != [2]int32{300, 400} {
		t.Fatalf("unexpected selected-stage localcoord: %#v", got.LocalCoord)
	}

	s.stage = &Stage{
		def:         "stages/current.def",
		name:        "current",
		displayname: "Current Stage",
		stageCamera: stageCamera{localcoord: [2]int32{640, 480}},
		scale:       [2]float32{1.5, 1.25},
	}
	got = s.liveFightStageContext()
	if got.Def != "stages/current.def" || got.Name != "current" || got.DisplayName != "Current Stage" {
		t.Fatalf("unexpected current-stage data: %#v", got)
	}
	if got.LocalCoord != [2]int32{640, 480} || got.CameraLocalCoord != [2]int32{640, 480} {
		t.Fatalf("unexpected current-stage coordinates: %#v", got)
	}
	if got.Scale != [2]float32{1.5, 1.25} {
		t.Fatalf("unexpected current-stage scale: %#v", got.Scale)
	}
}

func TestLiveSelectedTeams_UsesSelectionAndActiveCharacterOverrides(t *testing.T) {
	s := &System{}
	s.sel.charlist = []SelectChar{
		{def: "chars/ryu/ryu.def", name: "Ryu", lifebarname: "RYU", author: "Capcom", localcoord: [2]int32{320, 240}, cns_scale: [2]float32{1, 1}},
	}
	s.sel.selected[0] = [][2]int{{0, 3}}

	live := &Char{
		name:     "Ryu Custom",
		teamside: 0,
		memberNo: 0,
		selectNo: 0,
		lifeMax:  1000,
	}
	live.gi().displayname = "Ryu Live"
	live.gi().palno = 7
	s.chars[0] = []*Char{live}

	teams := s.liveSelectedTeams()
	got := teams[0][0]
	if got.Def != "chars/ryu/ryu.def" || got.PaletteNo != 7 {
		t.Fatalf("unexpected fighter selection data: %#v", got)
	}
	if got.Name != "Ryu Custom" || got.DisplayName != "Ryu Live" {
		t.Fatalf("unexpected live fighter overrides: %#v", got)
	}
	if got.Health == nil || got.Health.Status != "loaded" {
		t.Fatalf("expected loaded health marker for active fighter, got %#v", got.Health)
	}
}
