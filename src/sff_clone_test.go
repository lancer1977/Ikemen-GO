package main

import "testing"

func TestSffGetSpriteAndCloneSpriteWithPal(t *testing.T) {
	orig := sys.mainThreadTask
	defer func() { sys.mainThreadTask = orig }()
	sys.mainThreadTask = make(chan func(), 1)

	sff := &Sff{sprites: make(map[[2]uint16]*Sprite)}
	if got := sff.GetSprite(0xFFFF, 0); got != nil {
		t.Fatalf("GetSprite(0xFFFF, 0) = %#v, want nil", got)
	}

	src := &Sprite{Pal: []uint32{1, 2}, palidx: 3}
	sff.sprites[[2]uint16{1, 2}] = src
	pl := PaletteList{}
	pl.init()
	pl.SetSource(3, []uint32{9, 8, 7})

	got := sff.cloneSpriteWithPal(1, 2, &pl)
	if got == nil {
		t.Fatal("cloneSpriteWithPal should return a sprite")
	}
	if got == src {
		t.Fatal("cloneSpriteWithPal should copy the sprite value")
	}
	if got.Pal == src.Pal {
		t.Fatal("cloneSpriteWithPal should deep-copy the palette")
	}
	if len(got.Pal) != 3 || got.Pal[0] != 9 || got.Pal[2] != 7 {
		t.Fatalf("cloneSpriteWithPal copied wrong palette: %#v", got.Pal)
	}
	if got.palidx != src.palidx {
		t.Fatalf("cloneSpriteWithPal should preserve palidx, got %d want %d", got.palidx, src.palidx)
	}

	if got := sff.cloneSpriteWithPal(9, 9, &pl); got != nil {
		t.Fatalf("cloneSpriteWithPal missing sprite = %#v, want nil", got)
	}
}
