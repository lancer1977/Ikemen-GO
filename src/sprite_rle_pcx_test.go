package main

import (
	"bytes"
	"testing"
)

func TestSpriteRlePcxDecode(t *testing.T) {
	t.Run("early_return_when_not_rle", func(t *testing.T) {
		s := &Sprite{rle: 0}
		if got := s.RlePcxDecode([]byte{1, 2, 3}); len(got) != 3 || got[0] != 1 {
			t.Fatalf("expected passthrough, got %#v", got)
		}
	})

	t.Run("decode_with_line_wrap", func(t *testing.T) {
		s := &Sprite{Size: [2]uint16{2, 2}, rle: 2}
		got := s.RlePcxDecode([]byte{1, 2, 3, 4})
		want := []byte{1, 2, 3, 4}
		if !bytes.Equal(got, want) {
			t.Fatalf("unexpected decoded output: %#v", got)
		}
		if s.rle != 0 {
			t.Fatalf("RlePcxDecode should clear rle after decode, got %d", s.rle)
		}
	})
}
