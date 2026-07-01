package main

import (
	"encoding/binary"
	"reflect"
	"testing"
)

type fakeTexture struct {
	data [][]byte
}

func (t *fakeTexture) SetData(data []byte) {
	t.data = append(t.data, append([]byte{}, data...))
}

func (t *fakeTexture) SetSubData(data []byte, x, y, width, height, stride int32)   {}
func (t *fakeTexture) SetDataG(data []byte, mag, min, ws, wt TextureSamplingParam) {}
func (t *fakeTexture) SetPixelData(data []float32)                                 {}
func (t *fakeTexture) IsValid() bool                                               { return true }
func (t *fakeTexture) GetWidth() int32                                             { return 0 }
func (t *fakeTexture) GetHeight() int32                                            { return 0 }
func (t *fakeTexture) CopyData(src *Texture)                                       {}

func TestSpriteDefaultsAndCachePalTex(t *testing.T) {
	t.Run("newSprite", func(t *testing.T) {
		s := newSprite()
		if s.palidx != -1 {
			t.Fatalf("expected palidx -1, got %d", s.palidx)
		}
		if !s.isBlank() {
			t.Fatalf("expected new sprite to be blank")
		}
		s.Size = [2]uint16{2, 1}
		if !s.isBlank() {
			t.Fatalf("expected sprite with nil texture to be blank")
		}
		s.Tex = &fakeTexture{}
		if s.isBlank() {
			t.Fatalf("expected sprite with size and texture to be non-blank")
		}
	})

	t.Run("cache hit and miss", func(t *testing.T) {
		s := &Sprite{
			PalTex:  &fakeTexture{},
			paltemp: []uint32{1, 2, 3},
		}
		paltex := s.PalTex

		got := s.CachePalTex([]uint32{1, 2, 3})
		if got != paltex {
			t.Fatalf("expected cache hit to reuse texture")
		}
		if ft := s.PalTex.(*fakeTexture); len(ft.data) != 0 {
			t.Fatalf("expected no texture writes on cache hit, got %d", len(ft.data))
		}

		got = s.CachePalTex([]uint32{4, 5, 6})
		if got != paltex {
			t.Fatalf("expected cache miss to update existing texture")
		}
		if ft := s.PalTex.(*fakeTexture); len(ft.data) != 1 || len(ft.data[0]) != 1024 {
			t.Fatalf("unexpected texture writes: %v", ft.data)
		} else {
			want := make([]byte, 12)
			binary.LittleEndian.PutUint32(want[0:], 4)
			binary.LittleEndian.PutUint32(want[4:], 5)
			binary.LittleEndian.PutUint32(want[8:], 6)
			if !reflect.DeepEqual(ft.data[0][:12], want) {
				t.Fatalf("unexpected texture prefix: %v", ft.data[0][:12])
			}
		}
		if !reflect.DeepEqual(s.paltemp, []uint32{4, 5, 6}) {
			t.Fatalf("expected paltemp to update, got %v", s.paltemp)
		}
	})
}
