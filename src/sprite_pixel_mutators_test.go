package main

import "testing"

func TestSpriteSetPxlAndSetRaw_ValidateBufferShapeAndQueueWork(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.mainThreadTask = make(chan func(), 2)

	s := &Sprite{Size: [2]uint16{2, 2}}
	s.SetPxl(nil)
	if len(sys.mainThreadTask) != 0 {
		t.Fatal("SetPxl(nil) should not queue work")
	}

	s.SetPxl([]byte{1, 2, 3})
	if len(sys.mainThreadTask) != 0 {
		t.Fatal("SetPxl(wrong size) should not queue work")
	}

	s.SetPxl([]byte{1, 2, 3, 4})
	if len(sys.mainThreadTask) != 1 {
		t.Fatalf("SetPxl(valid) queued %d tasks, want 1", len(sys.mainThreadTask))
	}
	<-sys.mainThreadTask

	s.SetRaw([]byte{9, 8, 7}, 1, 1, 8)
	if len(sys.mainThreadTask) != 1 {
		t.Fatalf("SetRaw should queue one task, got %d", len(sys.mainThreadTask))
	}
}

func TestSpriteGetPalAndGetPalTexFallbacks(t *testing.T) {
	pl := PaletteList{}
	pl.init()
	pl.SetSource(0, []uint32{11})
	pl.SetSource(1, []uint32{22})
	tex0 := &fakeTexture{}
	tex1 := &fakeTexture{}
	pl.PalTex[0] = tex0
	pl.PalTex[1] = tex1

	s := &Sprite{palidx: 1}
	if got := s.GetPal(&pl); len(got) != 1 || got[0] != 22 {
		t.Fatalf("GetPal() = %#v, want palette slot 1", got)
	}

	s.Pal = []uint32{33}
	if got := s.GetPal(&pl); len(got) != 1 || got[0] != 33 {
		t.Fatalf("GetPal() should prefer sprite palette, got %#v", got)
	}

	if got := s.GetPalTex(&pl); got != tex1 {
		t.Fatalf("GetPalTex() should return palette texture for slot 1, got %#v", got)
	}

	s = &Sprite{palidx: 9}
	if got := s.GetPalTex(&pl); got != tex0 {
		t.Fatalf("GetPalTex() should fall back to slot 0 when index is out of range, got %#v", got)
	}

	s = &Sprite{coldepth: 16}
	if got := s.GetPalTex(&pl); got != nil {
		t.Fatalf("GetPalTex() should return nil for true-color sprites, got %#v", got)
	}

	pl.PalTex = nil
	s = &Sprite{palidx: 1}
	if got := s.GetPalTex(&pl); got != nil {
		t.Fatalf("GetPalTex() should return nil when texture slots are missing, got %#v", got)
	}

	var empty PaletteList
	if got := s.GetPalTex(&empty); got != nil {
		t.Fatalf("GetPalTex() should return nil for empty palette maps, got %#v", got)
	}
}
