package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestSpriteReadV2(t *testing.T) {
	t.Run("rle_positive_noop", func(t *testing.T) {
		s := &Sprite{rle: 1}
		if err := s.readV2(bytes.NewReader(nil), 0, 0); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("eight_bit_sets_pixels", func(t *testing.T) {
		oldSys := sys
		defer func() { sys = oldSys }()
		sys.mainThreadTask = make(chan func(), 1)

		s := &Sprite{Size: [2]uint16{2, 1}, coldepth: 8}
		if err := s.readV2(bytes.NewReader([]byte{1, 2}), 0, 2); err != nil {
			t.Fatal(err)
		}
		if len(sys.mainThreadTask) != 1 {
			t.Fatalf("expected SetPxl to queue a task, got %d", len(sys.mainThreadTask))
		}
	})

	t.Run("unknown_depth_errors", func(t *testing.T) {
		s := &Sprite{Size: [2]uint16{1, 1}, coldepth: 16}
		if err := s.readV2(bytes.NewReader([]byte{1, 2, 3, 4}), 0, 4); err == nil {
			t.Fatal("expected unknown color depth error")
		}
	})

	t.Run("raw_32bit_sets_raw_texture", func(t *testing.T) {
		oldSys := sys
		defer func() { sys = oldSys }()
		sys.mainThreadTask = make(chan func(), 1)

		s := &Sprite{Size: [2]uint16{1, 1}, coldepth: 32, rle: 0}
		if err := s.readV2(bytes.NewReader([]byte{1, 2, 3, 4}), 0, 4); err != nil {
			t.Fatal(err)
		}
		if len(sys.mainThreadTask) != 1 {
			t.Fatalf("expected raw 32-bit path to queue SetRaw, got %d", len(sys.mainThreadTask))
		}
	})

	t.Run("paletted_png_sets_pixels", func(t *testing.T) {
		paletted := image.NewPaletted(image.Rect(0, 0, 1, 1), color.Palette{
			color.RGBA{0, 0, 0, 0},
			color.RGBA{255, 0, 0, 255},
		})
		paletted.Pix[0] = 1
		var buf bytes.Buffer
		if err := png.Encode(&buf, paletted); err != nil {
			t.Fatal(err)
		}

		oldSys := sys
		defer func() { sys = oldSys }()
		sys.mainThreadTask = make(chan func(), 1)

		s := &Sprite{Size: [2]uint16{1, 1}, rle: -10}
		if err := s.readV2(bytes.NewReader(buf.Bytes()), 0, uint32(buf.Len())); err != nil {
			t.Fatal(err)
		}
		if len(sys.mainThreadTask) != 1 {
			t.Fatalf("expected paletted PNG path to queue SetPxl, got %d", len(sys.mainThreadTask))
		}
	})

	t.Run("rgba_png_sets_raw_texture", func(t *testing.T) {
		rgba := image.NewRGBA(image.Rect(0, 0, 1, 1))
		rgba.Pix[0] = 9
		rgba.Pix[1] = 8
		rgba.Pix[2] = 7
		rgba.Pix[3] = 6
		var buf bytes.Buffer
		if err := png.Encode(&buf, rgba); err != nil {
			t.Fatal(err)
		}

		oldSys := sys
		defer func() { sys = oldSys }()
		sys.mainThreadTask = make(chan func(), 1)

		s := &Sprite{Size: [2]uint16{1, 1}, rle: -11}
		if err := s.readV2(bytes.NewReader(buf.Bytes()), 0, uint32(buf.Len())); err != nil {
			t.Fatal(err)
		}
		if len(sys.mainThreadTask) != 1 {
			t.Fatalf("expected RGBA PNG path to queue SetRaw, got %d", len(sys.mainThreadTask))
		}
	})
}
