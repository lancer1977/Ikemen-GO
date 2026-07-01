package main

import "testing"

func TestNewRollbackSession(t *testing.T) {
	t.Parallel()

	rs := NewRollbackSession(RollbackProperties{})
	if rs.saveStates == nil || rs.players == nil || rs.handles == nil {
		t.Fatal("expected constructor to allocate core collections")
	}
	if rs.loopTimer.framesToSpreadWait != 100 {
		t.Fatalf("loopTimer.framesToSpreadWait = %d, want 100", rs.loopTimer.framesToSpreadWait)
	}
	if rs.timestamp == "" {
		t.Fatal("expected timestamp to be set")
	}
	if rs.log.filename == "" {
		t.Fatal("expected rollback logger filename to be set")
	}
	if len(rs.replayInputs) != 0 || len(rs.replayAnalogInputs) != 0 {
		t.Fatalf("expected empty replay buffers, got %#v %#v", rs.replayInputs, rs.replayAnalogInputs)
	}
}
