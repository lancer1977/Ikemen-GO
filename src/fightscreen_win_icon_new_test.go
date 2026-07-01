package main

import "testing"

func TestNewFightScreenWinIcon(t *testing.T) {
	t.Parallel()

	wi := newFightScreenWinIcon()
	if wi == nil {
		t.Fatal("newFightScreenWinIcon returned nil")
	}
	if wi.useiconupto != 4 {
		t.Fatalf("useiconupto = %d, want 4", wi.useiconupto)
	}
}
