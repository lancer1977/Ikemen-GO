package main

import "testing"

func TestNewFightScreenTimer(t *testing.T) {
	t.Parallel()

	tr := newFightScreenTimer()
	if tr == nil {
		t.Fatal("newFightScreenTimer returned nil")
	}
	if tr.enabled == nil {
		t.Fatal("newFightScreenTimer should allocate enabled map")
	}
	if tr.active {
		t.Fatal("newFightScreenTimer should start inactive")
	}
}
