package main

import "testing"

func TestStringPool_AddAndClear_DeduplicatesAndResets(t *testing.T) {
	sp := NewStringPool()
	if sp == nil || sp.Map == nil {
		t.Fatalf("NewStringPool() = %#v", sp)
	}

	if got := sp.Add("alpha"); got != 0 {
		t.Fatalf("first Add() = %d, want 0", got)
	}
	if got := sp.Add("alpha"); got != 0 {
		t.Fatalf("dedup Add() = %d, want 0", got)
	}
	if got := sp.Add("beta"); got != 1 {
		t.Fatalf("second unique Add() = %d, want 1", got)
	}
	if len(sp.List) != 2 {
		t.Fatalf("List len = %d, want 2", len(sp.List))
	}

	sp.Clear()
	if len(sp.List) != 0 || len(sp.Map) != 0 {
		t.Fatalf("Clear() = %#v", sp)
	}
}

func TestBytecodeExp_ReadHelpersAdvanceCursorAndFollowStringPool(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.workingState = &StateBytecode{playerNo: 0}
	sys.stringPool[0] = StringPool{List: []string{"zero", "one"}, Map: map[string]int{"zero": 0, "one": 1}}

	raw := []byte{0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x03, 0, 0, 0, 0x09, 0, 0, 0}
	be := BytecodeExp(raw)
	i := 0

	if got := be.ReadIntAt(&i); got != 2 || i != 4 {
		t.Fatalf("ReadIntAt() = %d, cursor %d", got, i)
	}
	if got := be.PeekLength(0); got != 2 {
		t.Fatalf("PeekLength(0) = %d, want 2", got)
	}
	if got := be.ReadPoolStringAt(&i); got != "one" || i != 8 {
		t.Fatalf("ReadPoolStringAt() = %q, cursor %d", got, i)
	}

	j := 8
	be.JumpToNext(&j)
	if j != 15 {
		t.Fatalf("JumpToNext() cursor = %d, want 15", j)
	}
}

func TestBytecodeExp_ReadHelpersHandleZeroLengthAndZeroIndex(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys.workingState = &StateBytecode{playerNo: 0}
	sys.stringPool[0] = StringPool{List: []string{"zero"}, Map: map[string]int{"zero": 0}}

	be := BytecodeExp([]byte{0, 0, 0, 0, 0, 0, 0, 0})
	if got := be.PeekLength(0); got != 0 {
		t.Fatalf("PeekLength(0) = %d, want 0", got)
	}

	i := 0
	if got := be.ReadPoolStringAt(&i); got != "zero" || i != 4 {
		t.Fatalf("ReadPoolStringAt() = %q, cursor %d", got, i)
	}

	j := 0
	be.JumpToNext(&j)
	if j != 4 {
		t.Fatalf("JumpToNext() cursor = %d, want 4", j)
	}
}
