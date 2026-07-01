package main

import "testing"

func TestNewAnimationTable_InitializesAnimMap(t *testing.T) {
	at := NewAnimationTable()
	if at.anims == nil || len(at.anims) != 0 {
		t.Fatalf("NewAnimationTable() = %#v, want empty anim map", at)
	}
	if at.filename != "" {
		t.Fatalf("NewAnimationTable() filename = %q, want empty", at.filename)
	}

	at.anims[1] = &Animation{}
	if at.anims[1] == nil {
		t.Fatal("NewAnimationTable anim map should be writable")
	}
}
