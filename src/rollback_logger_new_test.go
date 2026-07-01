package main

import "testing"

func TestNewRollbackLogger(t *testing.T) {
	t.Parallel()

	gl := NewRollbackLogger("2024-01-01T00:00:00Z")
	if gl.filename != "save/logs/Rollback-State-2024-01-01T00:00:00Z.log" {
		t.Fatalf("filename = %q, want timestamped path", gl.filename)
	}
	if gl.currentLog.Len() != 0 {
		t.Fatalf("expected empty current log, got %d", gl.currentLog.Len())
	}
}
