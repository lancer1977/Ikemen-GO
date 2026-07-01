package main

import "testing"

func TestNewFightScreenMatch(t *testing.T) {
	t.Parallel()

	ma := newFightScreenMatch()
	if ma == nil {
		t.Fatal("newFightScreenMatch returned nil")
	}
	if ma.enabled == nil {
		t.Fatal("newFightScreenMatch should allocate enabled map")
	}
	if ma.active {
		t.Fatal("newFightScreenMatch should start inactive")
	}
}
