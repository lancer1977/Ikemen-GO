package main

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestSpriteReadHeader(t *testing.T) {
	var buf bytes.Buffer
	for _, v := range []any{
		uint32(11),
		uint32(22),
		uint16(33),
		uint16(44),
		uint16(55),
		uint16(66),
	} {
		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			t.Fatal(err)
		}
	}

	s := &Sprite{}
	var ofs, size uint32
	var link uint16
	if err := s.readHeader(&buf, &ofs, &size, &link); err != nil {
		t.Fatal(err)
	}
	if ofs != 11 || size != 22 || s.Offset != [2]int16{33, 44} || s.Group != 55 || s.Number != 66 || link != 66 {
		t.Fatalf("unexpected readHeader result: ofs=%d size=%d sprite=%#v link=%d", ofs, size, s, link)
	}
}

func TestSpriteReadHeaderShortInput(t *testing.T) {
	s := &Sprite{}
	var ofs, size uint32
	var link uint16
	if err := s.readHeader(bytes.NewReader([]byte{1, 2, 3}), &ofs, &size, &link); err == nil {
		t.Fatal("expected readHeader to fail on short input")
	}
}

func TestSpriteReadPcxHeader(t *testing.T) {
	t.Run("invalid depth", func(t *testing.T) {
		var buf bytes.Buffer
		var header [128]byte
		binary.LittleEndian.PutUint16(header[0:], 0x0A01)
		header[2] = 1
		header[3] = 4
		binary.LittleEndian.PutUint16(header[4:], 0)
		binary.LittleEndian.PutUint16(header[6:], 0)
		binary.LittleEndian.PutUint16(header[8:], 7)
		binary.LittleEndian.PutUint16(header[10:], 9)
		binary.LittleEndian.PutUint16(header[66:], 4)
		buf.Write(header[:])

		s := &Sprite{}
		if err := s.readPcxHeader(bytes.NewReader(buf.Bytes()), 0); err == nil {
			t.Fatal("expected invalid depth error")
		}
	})

	t.Run("valid 8-bit encoding", func(t *testing.T) {
		var buf bytes.Buffer
		var header [128]byte
		binary.LittleEndian.PutUint16(header[0:], 0x0A01)
		header[2] = 1
		header[3] = 8
		binary.LittleEndian.PutUint16(header[4:], 0)
		binary.LittleEndian.PutUint16(header[6:], 0)
		binary.LittleEndian.PutUint16(header[8:], 3)
		binary.LittleEndian.PutUint16(header[10:], 4)
		binary.LittleEndian.PutUint16(header[66:], 2)
		buf.Write(header[:])

		s := &Sprite{}
		if err := s.readPcxHeader(bytes.NewReader(buf.Bytes()), 0); err != nil {
			t.Fatal(err)
		}
		if s.Size != [2]uint16{4, 5} || s.rle != 2 {
			t.Fatalf("unexpected PCX header parse: %#v", s)
		}
	})
}
