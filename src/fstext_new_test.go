package main

import "testing"

func TestNewFSText(t *testing.T) {
	ft := newFSText(2)
	if ft == nil {
		t.Fatal("expected newFSText to allocate")
	}
	if ft.font != [8]int32{-1, 0, 2, 255, 255, 255, 255, -1} {
		t.Fatalf("newFSText font = %#v", ft.font)
	}
	if ft.palfx == nil {
		t.Fatal("newFSText should initialize palfx")
	}
	if ft.frgba != [4]float32{1, 1, 1, 1} {
		t.Fatalf("newFSText frgba = %#v", ft.frgba)
	}
}
