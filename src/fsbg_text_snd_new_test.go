package main

import "testing"

func TestNewFSBgTextSnd(t *testing.T) {
	t.Parallel()

	bts := newFSBgTextSnd()
	if bts.snd != [2]int32{-1, 0} {
		t.Fatalf("snd = %#v, want [-1 0]", bts.snd)
	}
}
