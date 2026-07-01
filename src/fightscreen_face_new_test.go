package main

import "testing"

func TestNewFightScreenFace(t *testing.T) {
	t.Parallel()

	fa := newFightScreenFace()
	if fa == nil {
		t.Fatal("newFightScreenFace returned nil")
	}
	if fa.face_spr != [2]int32{-1, 0} || fa.teammate_face_spr != [2]int32{-1, 0} {
		t.Fatalf("unexpected sprite defaults: %#v", fa)
	}
	if !fa.face_palshare || !fa.teammate_face_palshare {
		t.Fatalf("unexpected palette-share defaults: %#v", fa)
	}
	if fa.face_pfx == nil {
		t.Fatal("newFightScreenFace should allocate face palfx")
	}
}
