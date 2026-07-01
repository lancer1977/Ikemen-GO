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

	dst := &Sprite{palidx: -1}
	dst.shareCopy(src)

	if dst.Pal != src.Pal {
		t.Fatalf("shareCopy should share palette slice")
	}
	if dst.Size != src.Size {
		t.Fatalf("shareCopy should copy size, got %#v want %#v", dst.Size, src.Size)
	}
	if dst.palidx != -1 {
		t.Fatalf("shareCopy should preserve existing palidx when already set")
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
