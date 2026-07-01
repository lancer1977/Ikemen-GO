package main

import "testing"

func TestNewFightScreenAiLevel(t *testing.T) {
	t.Parallel()

	ai := newFightScreenAiLevel()
	if ai == nil {
		t.Fatal("newFightScreenAiLevel returned nil")
	}
	if ai.separator != "." {
		t.Fatalf("separator = %q, want \".\"", ai.separator)
	}
	if ai.enabled == nil {
		t.Fatal("newFightScreenAiLevel should allocate enabled map")
	}
	if ai.active {
		t.Fatal("newFightScreenAiLevel should start inactive")
	}
}
