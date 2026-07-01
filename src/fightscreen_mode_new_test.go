package main

import "testing"

func TestNewFightScreenMode(t *testing.T) {
	t.Parallel()

	mo := newFightScreenMode()
	if mo == nil {
		t.Fatal("newFightScreenMode returned nil")
	}
}
