package main

import (
	"testing"
	"time"
)

func TestNewLoopTimer(t *testing.T) {
	t.Parallel()

	lt := NewLoopTimer(60, 3)
	if lt.framesToSpreadWait != 3 {
		t.Fatalf("framesToSpreadWait = %d, want 3", lt.framesToSpreadWait)
	}
	if lt.usPergameLoop != time.Second/60 {
		t.Fatalf("usPergameLoop = %v, want %v", lt.usPergameLoop, time.Second/60)
	}
	if lt.lastAdvantage != 0 || lt.waitCount != 0 || lt.timeWait != 0 || lt.waitTotal != 0 {
		t.Fatalf("expected zeroed runtime fields, got %#v", lt)
	}
}
