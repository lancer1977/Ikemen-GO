package main

import "testing"

func TestNewStateBlockAndBytecode(t *testing.T) {
	sb := newStateBlock()
	if sb == nil {
		t.Fatal("expected newStateBlock to allocate")
	}
	if sb.persistent != 1 || sb.persistentIndex != -1 || sb.ignorehitpause != -2 {
		t.Fatalf("newStateBlock defaults = %#v", sb)
	}

	bc := newStateBytecode(7)
	if bc == nil {
		t.Fatal("expected newStateBytecode to allocate")
	}
	if bc.stateType != ST_S || bc.moveType != MT_I || bc.physics != ST_N || bc.playerNo != 7 {
		t.Fatalf("newStateBytecode defaults = %#v", bc)
	}
	if bc.block.persistent != 1 || bc.block.persistentIndex != -1 || bc.block.ignorehitpause != -2 {
		t.Fatalf("newStateBytecode block defaults = %#v", bc.block)
	}
}
