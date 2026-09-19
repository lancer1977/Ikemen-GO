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

	pl := PaletteList{}
	pl.init()
	pl.SetSource(3, []uint32{9, 8, 7})

	// Sprite.GetPal has two branches: a sprite carrying its own Pal (or a
	// colour depth above 8) uses that palette directly, and everything else
	// falls through to the global PaletteList. cloneSpriteWithPal deep-copies
	// whichever one GetPal returns, so both paths need covering -- with an own
	// palette in play the deep-copy assertion has something to compare against.
	t.Run("own_palette", func(t *testing.T) {
		src := &Sprite{Pal: []uint32{1, 2}, palidx: 3}
		sff.sprites[[2]uint16{1, 2}] = src

		got := sff.cloneSpriteWithPal(1, 2, &pl)
		if got == nil {
			t.Fatal("cloneSpriteWithPal should return a sprite")
		}
		if got == src {
			t.Fatal("cloneSpriteWithPal should copy the sprite value")
		}
		// GetPal short-circuits on the sprite's own Pal, so the PaletteList
		// entry at palidx 3 is deliberately not consulted here.
		if len(got.Pal) != 2 || got.Pal[0] != 1 || got.Pal[1] != 2 {
			t.Fatalf("cloneSpriteWithPal copied wrong palette: %#v", got.Pal)
		}
		if &got.Pal[0] == &src.Pal[0] {
			t.Fatal("cloneSpriteWithPal should deep-copy the palette, not alias it")
		}
		if got.palidx != src.palidx {
			t.Fatalf("cloneSpriteWithPal should preserve palidx, got %d want %d", got.palidx, src.palidx)
		}
	})

	t.Run("palette_list_lookup", func(t *testing.T) {
		src := &Sprite{palidx: 3}
		sff.sprites[[2]uint16{4, 5}] = src

		got := sff.cloneSpriteWithPal(4, 5, &pl)
		if got == nil {
			t.Fatal("cloneSpriteWithPal should return a sprite")
		}
		if len(got.Pal) != 3 || got.Pal[0] != 9 || got.Pal[2] != 7 {
			t.Fatalf("cloneSpriteWithPal did not fetch palidx 3 from the list: %#v", got.Pal)
		}
		if listPal := pl.Get(3); len(listPal) > 0 && &got.Pal[0] == &listPal[0] {
			t.Fatal("cloneSpriteWithPal should deep-copy the list palette, not alias it")
		}
		if got.palidx != src.palidx {
			t.Fatalf("cloneSpriteWithPal should preserve palidx, got %d want %d", got.palidx, src.palidx)
		}
	})

	if got := sff.cloneSpriteWithPal(9, 9, &pl); got != nil {
		t.Fatalf("cloneSpriteWithPal missing sprite = %#v, want nil", got)
	}
}
