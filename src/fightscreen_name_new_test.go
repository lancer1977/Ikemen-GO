package main

import "testing"

func TestNewFightScreenName(t *testing.T) {
	t.Parallel()

	nm := newFightScreenName()
	if nm == nil {
		t.Fatal("newFightScreenName returned nil")
	}
}
