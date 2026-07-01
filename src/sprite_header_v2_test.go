package main

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestSpriteReadHeaderV2(t *testing.T) {
	makeHeader := func(flag uint16) []byte {
		var buf bytes.Buffer
		for _, v := range []any{
			uint16(10),
			uint16(20),
			[2]uint16{30, 40},
			[2]int16{50, 60},
			uint16(70),
			byte(4),
			byte(8),
			uint32(100),
			uint32(200),
			uint16(300),
			uint16(flag),
		} {
			if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
				t.Fatal(err)
			}
		}
		return buf.Bytes()
	}

	t.Run("uses_lofs_when_flag_is_even", func(t *testing.T) {
		s := &Sprite{}
		var ofs, size uint32
		var link uint16
		if err := s.readHeaderV2(bytes.NewReader(makeHeader(0)), &ofs, &size, 11, 22, &link); err != nil {
			t.Fatal(err)
		}
		if ofs != 111 || size != 200 || link != 70 || s.palidx != 300 {
			t.Fatalf("unexpected header values: ofs=%d size=%d link=%d sprite=%#v", ofs, size, link, s)
		}
	})

	t.Run("uses_tofs_when_flag_is_odd", func(t *testing.T) {
		s := &Sprite{}
		var ofs, size uint32
		var link uint16
		if err := s.readHeaderV2(bytes.NewReader(makeHeader(1)), &ofs, &size, 11, 22, &link); err != nil {
			t.Fatal(err)
		}
		if ofs != 122 || size != 200 || link != 70 || s.palidx != 300 {
			t.Fatalf("unexpected header values: ofs=%d size=%d link=%d sprite=%#v", ofs, size, link, s)
		}
	})
}
