package main

import "testing"

func TestNewChar(t *testing.T) {
	prevAI := sys.aiLevel
	defer func() { sys.aiLevel = prevAI }()

	c := newChar(2, 0)
	if c == nil {
		t.Fatal("expected newChar to allocate")
	}
	if c.playerNo != 2 || c.helperIndex != 0 || c.controller != 2 {
		t.Fatalf("newChar identity defaults = %#v", c)
	}
	if c.id != -1 || c.parentId != -1 || c.hoverIdx != -1 || c.winquote != -1 {
		t.Fatalf("newChar sentinel defaults = %#v", c)
	}
	if c.facing != 1 || c.minus != 3 || c.zScale != 1 || !c.ownpal {
		t.Fatalf("newChar gameplay defaults = %#v", c)
	}
	if !c.kovelocity || c.helperType != 0 {
		t.Fatalf("newChar player defaults = %#v", c)
	}
	if c.keyctrl != [4]bool{true, true, true, true} {
		t.Fatalf("newChar keyctrl = %#v", c.keyctrl)
	}
	if c.controller != c.playerNo {
		t.Fatalf("newChar should keep controller on human input when AI level is zero")
	}
}
