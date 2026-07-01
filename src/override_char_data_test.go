package main

import "testing"

func TestNewOverrideCharDataInitializesDefaults(t *testing.T) {
	ocd := newOverrideCharData()
	if ocd == nil {
		t.Fatal("newOverrideCharData returned nil")
	}
	if ocd.life != -1 || ocd.lifeMax != -1 || ocd.power != -1 || ocd.dizzyPoints != -1 || ocd.guardPoints != -1 {
		t.Fatalf("unexpected OverrideCharData defaults: %#v", ocd)
	}
	if ocd.maps == nil {
		t.Fatal("newOverrideCharData should initialize maps")
	}
}
