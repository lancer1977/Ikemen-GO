package main

import "testing"

func TestNewSelectInitializesDefaultState(t *testing.T) {
	sel := newSelect()
	if sel == nil {
		t.Fatal("newSelect returned nil")
	}
	if sel.selectedStageNo != -1 {
		t.Fatalf("newSelect selectedStageNo = %d, want -1", sel.selectedStageNo)
	}
	if sel.charAnimPreload == nil || sel.stageAnimPreload == nil || sel.charSpritePreload == nil || sel.stageSpritePreload == nil || sel.palOverwrite == nil || sel.cdefOverwrite == nil || sel.music == nil || sel.gameParams == nil {
		t.Fatalf("newSelect should initialize internal maps and params: %#v", sel)
	}
	if !sel.charSpritePreload[[2]uint16{9000, 0}] || !sel.charSpritePreload[[2]uint16{9000, 1}] {
		t.Fatalf("newSelect should seed char sprite preload sentinels: %#v", sel.charSpritePreload)
	}
}
