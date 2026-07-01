package main

import "testing"

func TestNewFSMsg(t *testing.T) {
	prev := sys.fightScreen
	defer func() { sys.fightScreen = prev }()

	msg := newFSMsg(0)
	if msg == nil {
		t.Fatal("expected newFSMsg to allocate")
	}
	if msg.resttime != -1 || msg.counterX != 0 || msg.fontNo != -1 || msg.fontBank != -1 || msg.fontAlign != IErr {
		t.Fatalf("newFSMsg defaults = %#v", msg)
	}
	if msg.fontColor != [4]int32{255, 255, 255, 255} || msg.snd != [2]int32{-1, -1} || msg.spr != [2]int32{-1, -1} || msg.animNo != -1 || !msg.top {
		t.Fatalf("newFSMsg defaults = %#v", msg)
	}

	sys.fightScreen.actions = []*FSAction{{start_x: 7}}
	msg = newFSMsg(0)
	if msg.counterX != 14 {
		t.Fatalf("newFSMsg counterX = %v, want 14", msg.counterX)
	}
}
