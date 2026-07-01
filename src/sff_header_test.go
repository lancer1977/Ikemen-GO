package main

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestSffHeaderRead(t *testing.T) {
	t.Run("invalid signature", func(t *testing.T) {
		sh := &SffHeader{}
		if err := sh.Read(bytes.NewReader([]byte("bad header")), new(uint32), new(uint32)); err == nil {
			t.Fatal("expected invalid signature error")
		}
	})

	t.Run("unsupported version", func(t *testing.T) {
		var buf bytes.Buffer
		buf.WriteString("ElecbyteSpr\x00")
		for _, v := range []any{
			byte(3), byte(0), byte(0), byte(0),
			uint32(0),
		} {
			if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
				t.Fatal(err)
			}
		}

		sh := &SffHeader{}
		if err := sh.Read(bytes.NewReader(buf.Bytes()), new(uint32), new(uint32)); err == nil {
			t.Fatal("expected unsupported version error")
		}
	})

	t.Run("version 1", func(t *testing.T) {
		var buf bytes.Buffer
		buf.WriteString("ElecbyteSpr\x00")
		for _, v := range []any{
			byte(1), byte(0), byte(0), byte(0),
			uint32(0),
			uint32(7),
			uint32(8),
			uint32(9),
		} {
			if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
				t.Fatal(err)
			}
		}

		var sh SffHeader
		var lofs, tofs uint32
		if err := sh.Read(bytes.NewReader(buf.Bytes()), &lofs, &tofs); err != nil {
			t.Fatal(err)
		}
		if sh.Version != [4]byte{1, 0, 0, 0} || sh.FirstPaletteHeaderOffset != 0 || sh.NumberOfPalettes != 0 {
			t.Fatalf("unexpected version 1 header: %#v", sh)
		}
		if sh.NumberOfSprites != 7 || sh.FirstSpriteHeaderOffset != 8 {
			t.Fatalf("unexpected version 1 fields: %#v", sh)
		}
	})

	t.Run("version 2", func(t *testing.T) {
		var buf bytes.Buffer
		buf.WriteString("ElecbyteSpr\x00")
		for _, v := range []any{
			byte(2), byte(0), byte(0), byte(0),
			uint32(0),
			uint32(0), uint32(0), uint32(0), uint32(0),
			uint32(11),
			uint32(12),
			uint32(13),
			uint32(14),
			uint32(15),
			uint32(16),
		} {
			if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
				t.Fatal(err)
			}
		}

		var sh SffHeader
		var lofs, tofs uint32
		if err := sh.Read(bytes.NewReader(buf.Bytes()), &lofs, &tofs); err != nil {
			t.Fatal(err)
		}
		if sh.Version != [4]byte{2, 0, 0, 0} {
			t.Fatalf("unexpected version: %#v", sh.Version)
		}
		if sh.FirstSpriteHeaderOffset != 11 || sh.NumberOfSprites != 12 || sh.FirstPaletteHeaderOffset != 13 || sh.NumberOfPalettes != 14 {
			t.Fatalf("unexpected version 2 fields: %#v", sh)
		}
		if lofs != 15 || tofs != 16 {
			t.Fatalf("unexpected offsets: lofs=%d tofs=%d", lofs, tofs)
		}
	})
}
