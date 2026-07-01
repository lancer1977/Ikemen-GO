package main

import "testing"

func TestCharacterLookupHelpersRespectMapsOrderAndFlags(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.charList.idMap = make(map[int32]*Char)
	sys.charList.creationOrder = nil
	sys.charList.runOrder = nil
	sys.chars = [MaxPlayerNo][]*Char{}

	activeID := &Char{}
	disabledID := &Char{}
	destroyedID := &Char{}
	activeIndex := &Char{}
	disabledIndex := &Char{}
	root := &Char{}
	disabledRoot := &Char{}

	disabledID.setSCF(SCF_disabled)
	destroyedID.setCSF(CSF_destroy)
	disabledIndex.setSCF(SCF_disabled)
	disabledRoot.setSCF(SCF_disabled)

	sys.charList.idMap[1] = activeID
	sys.charList.idMap[2] = disabledID
	sys.charList.idMap[3] = destroyedID
	sys.charList.creationOrder = []*Char{activeIndex, disabledIndex}
	sys.charList.runOrder = []*Char{activeIndex, disabledIndex}
	sys.chars[0] = []*Char{root}
	sys.chars[1] = []*Char{disabledRoot}

	if got := sys.playerID(1); got != activeID {
		t.Fatalf("playerID() returned %#v, want active char", got)
	}
	if got := sys.playerID(2); got != nil {
		t.Fatalf("playerID() should ignore disabled chars, got %#v", got)
	}
	if got := sys.playerID(3); got != nil {
		t.Fatalf("playerID() should ignore destroyed chars, got %#v", got)
	}

	if got := sys.playerIDExist(BytecodeInt(1)); got != BytecodeBool(true) {
		t.Fatalf("playerIDExist() for active char = %#v, want true", got)
	}
	if got := sys.playerIDExist(BytecodeInt(2)); got != BytecodeBool(false) {
		t.Fatalf("playerIDExist() for disabled char = %#v, want false", got)
	}
	if got := sys.playerIDExist(BytecodeUndefined()); !got.IsUndefined() {
		t.Fatalf("playerIDExist() should preserve undefined input, got %#v", got)
	}

	if got := sys.playerIndex(0); got != activeIndex {
		t.Fatalf("playerIndex() returned %#v, want first active char", got)
	}
	if got := sys.playerIndex(1); got != nil {
		t.Fatalf("playerIndex() should skip disabled chars, got %#v", got)
	}
	if got := sys.playerIndexExist(BytecodeInt(0)); got != BytecodeBool(true) {
		t.Fatalf("playerIndexExist() for active index = %#v, want true", got)
	}
	if got := sys.playerIndexExist(BytecodeInt(1)); got != BytecodeBool(false) {
		t.Fatalf("playerIndexExist() for skipped disabled index = %#v, want false", got)
	}

	if got := sys.getCharRoot(0); got != root {
		t.Fatalf("getCharRoot() returned %#v, want active root", got)
	}
	if got := sys.getCharRoot(1); got != nil {
		t.Fatalf("getCharRoot() should ignore disabled roots, got %#v", got)
	}
	if got := sys.playerNoExist(BytecodeInt(1)); got != BytecodeBool(true) {
		t.Fatalf("playerNoExist() for active root = %#v, want true", got)
	}
	if got := sys.playerNoExist(BytecodeInt(2)); got != BytecodeBool(false) {
		t.Fatalf("playerNoExist() for disabled root = %#v, want false", got)
	}
}
