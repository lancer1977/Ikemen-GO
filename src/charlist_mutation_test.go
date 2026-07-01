package main

import "testing"

func TestCharListDeleteRemovesPointerFromAllIndexes(t *testing.T) {
	cl := CharList{idMap: make(map[int32]*Char)}
	a := &Char{id: 1}
	b := &Char{id: 2}
	cl.creationOrder = []*Char{a, b}
	cl.runOrder = []*Char{b, a}
	cl.idMap[1] = a
	cl.idMap[2] = b

	cl.delete(a)
	if _, ok := cl.idMap[1]; ok {
		t.Fatal("delete() should remove the deleted pointer from idMap")
	}
	if len(cl.creationOrder) != 1 || cl.creationOrder[0] != b {
		t.Fatalf("delete() should remove from creationOrder, got %#v", cl.creationOrder)
	}
	if len(cl.runOrder) != 1 || cl.runOrder[0] != b {
		t.Fatalf("delete() should remove from runOrder, got %#v", cl.runOrder)
	}
}

func TestCharListReplaceSwapsMatchingSlotOnly(t *testing.T) {
	cl := CharList{idMap: make(map[int32]*Char)}
	oldChar := &Char{id: 1, playerNo: 0, helperIndex: 0}
	newChar := &Char{id: 2, playerNo: 0, helperIndex: 0}
	other := &Char{id: 3, playerNo: 1, helperIndex: 0}
	cl.creationOrder = []*Char{other, oldChar}
	cl.runOrder = []*Char{other, oldChar}
	cl.idMap[1] = oldChar
	cl.idMap[3] = other

	if !cl.replace(newChar, 0, 0) {
		t.Fatal("replace() should return true when it finds a matching slot occupant")
	}
	if _, ok := cl.idMap[1]; ok {
		t.Fatal("replace() should delete the old occupant from idMap")
	}
	if cl.idMap[2] != newChar {
		t.Fatalf("replace() should add the new char to idMap, got %#v", cl.idMap[2])
	}
	if len(cl.creationOrder) != 2 || cl.creationOrder[1] != newChar {
		t.Fatalf("replace() should preserve unrelated chars and append new char, got %#v", cl.creationOrder)
	}
	if len(cl.runOrder) != 2 || cl.runOrder[1] != newChar {
		t.Fatalf("replace() should preserve runOrder and append new char, got %#v", cl.runOrder)
	}

	if cl.replace(&Char{id: 4, playerNo: 9, helperIndex: 9}, 9, 9) {
		t.Fatal("replace() should return false when no matching slot occupant exists")
	}
}
