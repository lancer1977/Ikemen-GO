package main

import "testing"

func TestNewSnd(t *testing.T) {
	t.Parallel()

	snd := newSnd()
	if snd == nil {
		t.Fatal("newSnd returned nil")
	}
	if snd.table == nil {
		t.Fatal("newSnd should allocate table")
	}
	if len(snd.table) != 0 {
		t.Fatalf("newSnd table len = %d, want 0", len(snd.table))
	}
	if snd.ver != 0 || snd.ver2 != 0 || snd.filename != "" {
		t.Fatalf("unexpected snd defaults: %#v", snd)
	}
}
