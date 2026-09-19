package main

import (
	"testing"
	"time"
)

// fakeTexture implements Texture interface for testing
type fakeTexture struct{}

func (f *fakeTexture) SetData(data []byte)                                                      {}
func (f *fakeTexture) SetSubData(data []byte, x, y, width, height, stride int32)               {}
func (f *fakeTexture) SetDataG(data []byte, mag, min, ws, wt TextureSamplingParam)             {}
func (f *fakeTexture) SetPixelData(data []float32)                                              {}
func (f *fakeTexture) IsValid() bool                                                            { return true }
func (f *fakeTexture) GetWidth() int32                                                          { return 0 }
func (f *fakeTexture) GetHeight() int32                                                         { return 0 }
func (f *fakeTexture) CopyData(src *Texture)                                                    {}

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

	if len(dst.Pal) != len(src.Pal) || (len(dst.Pal) > 0 && &dst.Pal[0] != &src.Pal[0]) {
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
