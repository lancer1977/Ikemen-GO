package main

import (
	"reflect"
	"testing"
)

func TestSpriteRle8Decode(t *testing.T) {
	s := &Sprite{Size: [2]uint16{3, 1}}

	t.Run("empty", func(t *testing.T) {
		if got := s.Rle8Decode(nil); got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
		if got := s.Rle8Decode([]byte{}); len(got) != 0 {
			t.Fatalf("expected empty slice, got %v", got)
		}
	})

	t.Run("literal and repeat", func(t *testing.T) {
		got := s.Rle8Decode([]byte{0x41, 0x07, 0x02})
		want := []byte{0x07, 0x02, 0x02}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("expected %v, got %v", want, got)
		}
	})
}
