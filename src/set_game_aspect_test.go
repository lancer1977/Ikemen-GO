package main

import "testing"

func TestSetGameAspect_PropagatesStageAndFighterScaling(t *testing.T) {
	prevSys := sys
	defer func() { sys = prevSys }()

	stage := &Stage{
		stageCamera: stageCamera{localcoord: [2]int32{400, 300}},
	}
	fighter := &Char{}
	fighter.gi().localcoord = [2]int32{320, 240}

	sys = System{
		SystemStateVars: SystemStateVars{
			gameWidth:  400,
			gameHeight: 240,
		},
		stage: stage,
	}
	sys.fightScreen.localcoord = [2]int32{320, 240}
	sys.chars[0] = []*Char{fighter}

	sys.setGameAspect()

	if stage.localscl != 1 {
		t.Fatalf("expected stage localscl to be updated, got %v", stage.localscl)
	}
	if stage.stageCamera.localscl != 1 {
		t.Fatalf("expected stage camera localscl to be updated, got %v", stage.stageCamera.localscl)
	}
	if fighter.localcoord != 256 {
		t.Fatalf("unexpected fighter localcoord: %v", fighter.localcoord)
	}
	if fighter.localscl != 1.25 {
		t.Fatalf("unexpected fighter localscl: %v", fighter.localscl)
	}
}
