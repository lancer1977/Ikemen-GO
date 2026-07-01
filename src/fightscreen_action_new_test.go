package main

import "testing"

func TestNewFightScreenAction(t *testing.T) {
	t.Parallel()

	ac := newFightScreenAction()
	if ac == nil {
		t.Fatal("newFightScreenAction returned nil")
	}
	if ac.displaytime != 90 || ac.showspeed != 8 || ac.hidespeed != 4 || ac.max != 8 {
		t.Fatalf("unexpected action defaults: %#v", ac)
	}
	if ac.is == nil {
		t.Fatal("newFightScreenAction should allocate backing map")
	}
}
