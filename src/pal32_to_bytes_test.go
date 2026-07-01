package main

import (
	"encoding/binary"
	"testing"
)

func TestPal32ToBytes_HandlesEmptyShortAndFullPalettes(t *testing.T) {
	if got := Pal32ToBytes(nil); got != nil {
		t.Fatalf("Pal32ToBytes(nil) = %#v, want nil", got)
	}
	if got := Pal32ToBytes([]uint32{}); got != nil {
		t.Fatalf("Pal32ToBytes(empty) = %#v, want nil", got)
	}

	short := Pal32ToBytes([]uint32{0x11223344, 0x55667788})
	if len(short) != 1024 {
		t.Fatalf("Pal32ToBytes(short) length = %d, want 1024", len(short))
	}
	if binary.LittleEndian.Uint32(short[:4]) != 0x11223344 || binary.LittleEndian.Uint32(short[4:8]) != 0x55667788 {
		t.Fatalf("Pal32ToBytes(short) did not preserve leading entries")
	}

	full := make([]uint32, 256)
	full[0] = 0x01020304
	full[255] = 0xAABBCCDD
	got := Pal32ToBytes(full)
	if len(got) != 1024 {
		t.Fatalf("Pal32ToBytes(full) length = %d, want 1024", len(got))
	}
	if binary.LittleEndian.Uint32(got[:4]) != 0x01020304 || binary.LittleEndian.Uint32(got[1020:1024]) != 0xAABBCCDD {
		t.Fatalf("Pal32ToBytes(full) did not expose backing palette data")
	}
}
