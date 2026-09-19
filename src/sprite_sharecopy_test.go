package main

import (
	"testing"
	"time"
)

func TestSpriteShareCopy(t *testing.T) {
	orig := sys.mainThreadTask
	defer func() { sys.mainThreadTask = orig }()

	sys.mainThreadTask = make(chan func(), 1)

	src := &Sprite{
		Pal:      []uint32{1, 2, 3},
		Size:     [2]uint16{4, 5},
		coldepth: 16,
		Tex:      &fakeTexture{},
	}

	// shareCopy only adopts the source palidx when the destination has none
	// (a negative index). An index already set is left alone.
	dst := &Sprite{palidx: 42}
	dst.shareCopy(src)

	if len(dst.Pal) != len(src.Pal) || (len(dst.Pal) > 0 && &dst.Pal[0] != &src.Pal[0]) {
		t.Fatalf("shareCopy should share palette slice")
	}
	if dst.Size != src.Size {
		t.Fatalf("shareCopy should copy size, got %#v want %#v", dst.Size, src.Size)
	}
	if dst.palidx != 42 {
		t.Fatalf("shareCopy should preserve existing palidx when already set, got %d want 42", dst.palidx)
	}
	if dst.coldepth != src.coldepth {
		t.Fatalf("shareCopy should copy coldepth, got %d want %d", dst.coldepth, src.coldepth)
	}

	select {
	case fn := <-sys.mainThreadTask:
		fn()
	case <-time.After(time.Second):
		t.Fatal("shareCopy did not enqueue texture copy")
	}
	if dst.Tex != src.Tex {
		t.Fatalf("shareCopy should copy texture on main thread")
	}

	dst2 := &Sprite{palidx: -1}
	dst2.shareCopy(src)
	<-sys.mainThreadTask
	if dst2.palidx != src.palidx {
		t.Fatalf("shareCopy should copy palidx when unset")
	}
}

// The other half of the palidx guard: an unset destination adopts the source's
// index. Covering only the preserve branch would let the copy branch rot.
func TestSpriteShareCopy_AdoptsPalidxWhenUnset(t *testing.T) {
	src := &Sprite{
		Pal:      []uint32{1, 2, 3},
		Size:     [2]uint16{4, 5},
		palidx:   7,
		coldepth: 8,
	}

	dst := &Sprite{palidx: -1}
	dst.shareCopy(src)

	if dst.palidx != src.palidx {
		t.Fatalf("shareCopy should adopt palidx when unset, got %d want %d", dst.palidx, src.palidx)
	}
}
