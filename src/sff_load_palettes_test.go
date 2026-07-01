package main

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestSffLoadPalettes(t *testing.T) {
	t.Run("version1_noop", func(t *testing.T) {
		s := &Sff{header: SffHeader{Version: [4]byte{1, 0, 0, 0}}}
		s.palList.init()
		if err := s.loadPalettes(bytes.NewReader(nil), 0); err != nil {
			t.Fatal(err)
		}
		if len(s.palList.palettes) != 0 {
			t.Fatalf("version 1 should not load palettes, got %#v", s.palList.palettes)
		}
	})

	t.Run("version2_unique_and_linked", func(t *testing.T) {
		var buf bytes.Buffer
		writeHeader := func(gn0, gn1, gn2 uint16, link uint16, ofs, size uint32) {
			for _, v := range []any{gn0, gn1, gn2, link, ofs, size} {
				if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
					t.Fatal(err)
				}
			}
		}

		writeHeader(1, 1, 2, 0, 0, 8)
		writeHeader(1, 2, 2, 0, 0, 0)

		paletteData := []byte{
			1, 2, 3, 4,
			5, 6, 7, 8,
		}
		for buf.Len() < 64 {
			buf.WriteByte(0)
		}
		buf.Write(paletteData)

		s := &Sff{
			header: SffHeader{
				Version:                  [4]byte{2, 0, 0, 0},
				FirstPaletteHeaderOffset: 0,
				NumberOfPalettes:         2,
			},
		}
		s.palList.init()
		if err := s.loadPalettes(bytes.NewReader(buf.Bytes()), 64); err != nil {
			t.Fatal(err)
		}
		if len(s.palList.palettes) != 2 {
			t.Fatalf("expected two loaded palettes, got %d", len(s.palList.palettes))
		}
		if s.palList.PalTable[[2]uint16{1, 1}] != 0 || s.palList.PalTable[[2]uint16{1, 2}] != 1 {
			t.Fatalf("unexpected palette table: %#v", s.palList.PalTable)
		}
		if len(s.palList.palettes[1]) == 0 || s.palList.palettes[1][0] != s.palList.palettes[0][0] {
			t.Fatalf("linked palette should reuse palette 0, got %#v %#v", s.palList.palettes[0], s.palList.palettes[1])
		}
	})

	t.Run("version2_duplicate_key_reuses_existing_palette", func(t *testing.T) {
		var buf bytes.Buffer
		writeHeader := func(gn0, gn1, gn2 uint16, link uint16, ofs, size uint32) {
			for _, v := range []any{gn0, gn1, gn2, link, ofs, size} {
				if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
					t.Fatal(err)
				}
			}
		}

		writeHeader(1, 1, 2, 0, 0, 8)
		writeHeader(1, 1, 2, 0, 0, 0)

		for buf.Len() < 64 {
			buf.WriteByte(0)
		}
		buf.Write([]byte{
			1, 2, 3, 4,
			5, 6, 7, 8,
		})

		s := &Sff{
			header: SffHeader{
				Version:                  [4]byte{2, 0, 0, 0},
				FirstPaletteHeaderOffset: 0,
				NumberOfPalettes:         2,
			},
		}
		s.palList.init()
		if err := s.loadPalettes(bytes.NewReader(buf.Bytes()), 64); err != nil {
			t.Fatal(err)
		}
		if len(s.palList.palettes) != 2 {
			t.Fatalf("expected two palettes, got %d", len(s.palList.palettes))
		}
		if s.palList.PalTable[[2]uint16{1, 1}] != 0 {
			t.Fatalf("duplicate key should reuse first palette index, got %#v", s.palList.PalTable)
		}
		if s.palList.palettes[1] != s.palList.palettes[0] {
			t.Fatalf("duplicate key should reuse existing palette slice")
		}
	})
}
