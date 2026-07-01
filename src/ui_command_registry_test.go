package main

import "testing"

func TestUIRegisterCommandAddsOnceAndSkipsWhitespace(t *testing.T) {
	s := &System{}

	if err := s.uiRegisterCommand("   ", CommandSpec{Cmd: "a"}); err != nil {
		t.Fatalf("uiRegisterCommand() whitespace name returned error: %v", err)
	}
	if err := s.uiRegisterCommand("jump", CommandSpec{Cmd: "   "}); err != nil {
		t.Fatalf("uiRegisterCommand() whitespace command returned error: %v", err)
	}
	if len(s.uiCommandRegistry) != 0 {
		t.Fatalf("uiRegisterCommand() should ignore blank entries, got %#v", s.uiCommandRegistry)
	}

	spec := CommandSpec{Cmd: "a", Time: 12, BufTime: 3}
	if err := s.uiRegisterCommand("jump", spec); err != nil {
		t.Fatalf("uiRegisterCommand() = %v", err)
	}
	if got := s.uiCommandRegistry["jump"]; got.Cmd != "a" || got.Time != 12 || got.BufTime != 3 {
		t.Fatalf("uiCommandRegistry entry = %#v", got)
	}
	if err := s.uiRegisterCommand("jump", CommandSpec{Cmd: "b"}); err != nil {
		t.Fatalf("uiRegisterCommand() duplicate should be ignored, got %v", err)
	}
	if got := s.uiCommandRegistry["jump"]; got.Cmd != "a" {
		t.Fatalf("uiRegisterCommand() should preserve first spec, got %#v", got)
	}
}

func TestUIApplyCommandRegistryAddsMissingCommandsOnly(t *testing.T) {
	s := &System{}
	s.uiCommandRegistry = map[string]CommandSpec{
		"jump": {Cmd: "a", Time: 12},
		"kick": {Cmd: "b", Time: 8},
	}
	cl := NewCommandList(NewInputBuffer())
	if err := cl.AddCommand("jump", CommandSpec{Cmd: "a", Time: 12}); err != nil {
		t.Fatalf("AddCommand() setup failed: %v", err)
	}

	if err := s.uiApplyCommandRegistry(cl); err != nil {
		t.Fatalf("uiApplyCommandRegistry() = %v", err)
	}
	if _, ok := cl.Names["jump"]; !ok {
		t.Fatal("uiApplyCommandRegistry() should keep existing command names")
	}
	if _, ok := cl.Names["kick"]; !ok {
		t.Fatal("uiApplyCommandRegistry() should add missing commands")
	}
	if len(cl.Commands) != 2 {
		t.Fatalf("uiApplyCommandRegistry() command count = %d, want 2", len(cl.Commands))
	}
}
