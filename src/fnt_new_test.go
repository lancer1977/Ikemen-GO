package main

import "testing"

func TestNewFnt(t *testing.T) {
	t.Parallel()

	fnt := newFnt()
	if fnt == nil {
		t.Fatal("newFnt returned nil")
	}
	if fnt.images == nil || fnt.paltexCache == nil {
		t.Fatalf("newFnt should allocate maps: %#v", fnt)
	}
	if fnt.BankType != "palette" || fnt.lastPalBank != -1 {
		t.Fatalf("unexpected defaults: %#v", fnt)
	}
}
