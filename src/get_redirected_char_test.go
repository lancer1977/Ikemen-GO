package main

import ("testing"
	"unsafe")

func TestGetRedirectedChar(t *testing.T) {
	t.Parallel()

	oldSys := sys
	defer func() { sys = oldSys }()
	sys.charList.idMap = make(map[int32]*Char)
	sys.bcStack = nil

	src := &Char{}
	dst := &Char{}
	sys.charList.idMap[7] = dst

	exp := BytecodeExp{}
	exp.appendValue(BytecodeInt(7))
	block := StateControllerBase{
		3, 1, // param id, exp count
	}
	block = append(block, []byte{2, 0, 0, 0}...) // exp length
	block = append(block, unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(exp))), len(exp))...)

	if got := getRedirectedChar(src, block, 3, "redirectid"); got != dst {
		t.Fatalf("getRedirectedChar returned %#v, want %#v", got, dst)
	}
}
