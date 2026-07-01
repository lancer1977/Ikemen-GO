package main

import (
	"reflect"
	"testing"
)

func TestSetRect(t *testing.T) {
	t.Parallel()

	oldSys := sys
	defer func() { sys = oldSys }()
	sys.gameWidth = 640
	sys.gameHeight = 480
	sys.widthScale = 2
	sys.heightScale = 3
	sys.scrrect = [4]int32{0, 0, 640, 480}

	type sample struct {
		Rect       *Rect
		Time       int32
		Layerno    int32
		Window     [4]int32
		Col        [3]int32
		Alpha      [2]int32
		Localcoord [2]int32
	}

	s := sample{
		Time:       7,
		Layerno:    3,
		Window:     [4]int32{10, 20, 110, 220},
		Col:        [3]int32{1, 2, 3},
		Alpha:      [2]int32{4, 5},
		Localcoord: [2]int32{640, 480},
	}

	v := reflect.ValueOf(&s).Elem()
	SetRect(&s, v.FieldByName("Rect"), v, reflect.Value{})

	if s.Rect == nil {
		t.Fatal("expected Rect to be initialized")
	}
	if s.Rect.time != 7 || s.Rect.layerno != 3 {
		t.Fatalf("unexpected time/layer: %#v", s.Rect)
	}
	if s.Rect.col != 0x010203 {
		t.Fatalf("col = %#x, want 0x010203", s.Rect.col)
	}
	if s.Rect.alpha != [2]int32{4, 5} {
		t.Fatalf("alpha = %#v, want [4 5]", s.Rect.alpha)
	}
	if s.Rect.windowInit != [4]float32{10, 20, 110, 220} {
		t.Fatalf("windowInit = %#v", s.Rect.windowInit)
	}
}
