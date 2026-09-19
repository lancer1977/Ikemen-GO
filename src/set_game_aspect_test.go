package main

import "testing"

func TestSetGameAspect_PropagatesStageAndFighterScaling(t *testing.T) {
	prevSys := sys
	defer func() { sys = prevSys }()

	stage := &Stage{
		stageCamera: stageCamera{localcoord: [2]int32{400, 300}},
	}
	fighter := &Char{playerNo: 0} // Explicitly set playerNo

	sys = System{
		SystemStateVars: SystemStateVars{
			gameWidth:  400,
			gameHeight: 240,
		},
		stage: stage,
	}
	// Set up screen rect so that applyFightAspect can calculate aspect ratio.
	// setGameAspect calls applyFightAspect which uses getMotifAspect() which calculates
	// from scrrect. We set it to 400x240 (1.667 aspect) to get gameWidth = 240 * 1.667 = 400.
	sys.scrrect = [4]int32{0, 0, 400, 240}
	sys.fightScreen.localcoord = [2]int32{320, 240}
	sys.chars[0] = []*Char{fighter}
	// Set fighter's local coords after sys is set up (so it uses the new sys.cgi)
	fighter.gi().localcoord[0] = 320
	fighter.gi().localcoord[1] = 240

	// Call setGameAspect to test that it properly propagates scaling
	sys.setGameAspect()

	// Debug: verify the inputs were set correctly
	if fighter.gi().localcoord[0] != 320 {
		t.Fatalf("fighter gi localcoord not set correctly: %v (expected 320)", fighter.gi().localcoord[0])
	}
	if sys.gameWidth != 400 {
		t.Fatalf("gameWidth after setGameAspect: %v (expected 400)", sys.gameWidth)
	}

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
