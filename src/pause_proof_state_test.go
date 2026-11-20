package main

import "testing"

func TestUpdatePauseProof_AdvancesThroughCapturePhases(t *testing.T) {
	prevSeconds := pauseProofSeconds
	prevSys := sys
	t.Cleanup(func() {
		pauseProofSeconds = prevSeconds
		sys = prevSys
	})

	ensureGlobalRound(t)

	pauseProofSeconds = 1
	s := &System{
		bgm:       newBgm(),
		sel:       newSelect(),
		loader:    newLoader(),
		selMutex:  newTestRWMutex(),
		loadMutex: newTestMutex(),
		statePool: NewGameStatePool(),
		savePool:  NewGameStatePool(),
		loadPool:  NewGameStatePool(),
		SystemStateVars: SystemStateVars{
			matchTime:    60,
			curRoundTime: 120,
		},
	}
	s.fightScreen.active = true
	s.fightScreen.round = &FightScreenRound{} // Ensure round is set before copying to sys
	s.intro = 0
	s.maxRoundTime = 180
	// Set the -pauseproof flag to enable the pause proof system
	s.cmdFlags = map[string]string{"-pauseproof": ""}
	sys = *s

	s.updatePauseProof()
	if !s.pauseProofStarted || s.pauseProofCapturePhase != 1 || !s.isTakingScreenshot {
		t.Fatalf("expected pause proof to start capture, got started=%v phase=%d screenshot=%v", s.pauseProofStarted, s.pauseProofCapturePhase, s.isTakingScreenshot)
	}

	s.isTakingScreenshot = false
	s.pausetime = 120
	s.updatePauseProof()
	if s.pauseProofCapturePhase != 2 || !s.isTakingScreenshot {
		t.Fatalf("expected countdown capture phase, got phase=%d screenshot=%v", s.pauseProofCapturePhase, s.isTakingScreenshot)
	}

	s.isTakingScreenshot = false
	s.pausetime = 1
	s.updatePauseProof()
	if s.pauseProofCapturePhase != 3 || !s.isTakingScreenshot {
		t.Fatalf("expected final pause capture phase, got phase=%d screenshot=%v", s.pauseProofCapturePhase, s.isTakingScreenshot)
	}

	s.isTakingScreenshot = false
	s.pausetime = 0
	s.updatePauseProof()
	if s.pauseProofCapturePhase != 4 || !s.isTakingScreenshot {
		t.Fatalf("expected resumed capture phase, got phase=%d screenshot=%v", s.pauseProofCapturePhase, s.isTakingScreenshot)
	}
}

func TestUpdatePauseProof_IsNoOpWithoutFlag(t *testing.T) {
	s := newTestSystem()
	s.updatePauseProof()
	if s.pauseProofStarted || s.pauseProofCapturePhase != 0 || s.isTakingScreenshot {
		t.Fatalf("expected no-op without pauseproof flag, got started=%v phase=%d screenshot=%v", s.pauseProofStarted, s.pauseProofCapturePhase, s.isTakingScreenshot)
	}
}
