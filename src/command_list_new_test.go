package main

import "testing"

func TestNewCommandList(t *testing.T) {
	t.Parallel()

	cb := NewInputBuffer()
	cl := NewCommandList(cb)
	if cl == nil {
		t.Fatal("NewCommandList returned nil")
	}
	if cl.Buffer != cb {
		t.Fatal("expected buffer to be wired through")
	}
	if cl.Names == nil {
		t.Fatal("expected Names map to be initialized")
	}
	if cl.DefaultTime != 15 || cl.DefaultStepTime != -1 || !cl.DefaultAutoGreater || cl.DefaultBufferTime != 1 || !cl.DefaultBufferHitpause || !cl.DefaultBufferPauseEnd || !cl.DefaultBufferShared {
		t.Fatalf("unexpected defaults: %#v", cl)
	}
}
