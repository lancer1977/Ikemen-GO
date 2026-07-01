package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"
)

func TestSffReadPaletteClampsDepthAndForcesAlpha(t *testing.T) {
	var buf bytes.Buffer
	for _, rgba := range [][4]byte{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
	} {
		if err := binary.Write(&buf, binary.LittleEndian, rgba); err != nil {
			t.Fatal(err)
		}
	}

	s := &Sff{header: SffHeader{Version: [4]byte{2, 0, 0, 0}}}
	pal, err := s.ReadPalette(bytes.NewReader(buf.Bytes()), 0, uint32(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	if len(pal) != 16 {
		t.Fatalf("expected depth clamp to 16, got %d", len(pal))
	}
	if pal[0] != 0x00030201 {
		t.Fatalf("expected first color to force transparent alpha, got %#08x", pal[0])
	}
	if pal[1] != 0xff080706 {
		t.Fatalf("expected second color to force opaque alpha, got %#08x", pal[1])
	}
	if pal[2] != 0xff0c0b0a {
		t.Fatalf("expected third color to force opaque alpha, got %#08x", pal[2])
	}
	if pal[15] != 0xff000000 {
		t.Fatalf("expected padded colors to remain opaque black, got %#08x", pal[15])
	}

	buf.Reset()
	for i := 0; i < 257; i++ {
		rgba := [4]byte{byte(i), byte(i >> 1), byte(i >> 2), byte(i >> 3)}
		if err := binary.Write(&buf, binary.LittleEndian, rgba); err != nil {
			t.Fatal(err)
		}
	}

	s.header.Version = [4]byte{2, 0, 1, 0}
	pal, err = s.ReadPalette(bytes.NewReader(buf.Bytes()), 0, uint32(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	if len(pal) != 256 {
		t.Fatalf("expected depth clamp to 256, got %d", len(pal))
	}
	if pal[255] == 0 {
		t.Fatalf("expected last palette entry to be populated, got %#08x", pal[255])
	}
}

func TestSffReadPaletteRejectsBadSeek(t *testing.T) {
	s := &Sff{}
	if _, err := s.ReadPalette(errSeekReadSeeker{}, 0, 4); err == nil {
		t.Fatal("expected seek failure")
	}
}

type errSeekReadSeeker struct{}

func (errSeekReadSeeker) Read([]byte) (int, error) { return 0, nil }
func (errSeekReadSeeker) Seek(int64, int) (int64, error) {
	return 0, errors.New("seek failed")
}
