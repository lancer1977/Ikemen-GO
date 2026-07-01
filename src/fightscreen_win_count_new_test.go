package main

import "testing"

func TestNewFightScreenWinCount(t *testing.T) {
	t.Parallel()

	wc := newFightScreenWinCount()
	if wc == nil {
		t.Fatal("newFightScreenWinCount returned nil")
	}
	if wc.enabled == nil {
		t.Fatal("newFightScreenWinCount should allocate enabled map")
	}
	if wc.active {
		t.Fatal("newFightScreenWinCount should start inactive")
	}
}
