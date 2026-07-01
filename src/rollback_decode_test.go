package main

import "testing"

func TestReadI16AndI32(t *testing.T) {
	t.Parallel()

	if got := readI16([]byte{0x34, 0x12}); got != 0x1234 {
		t.Fatalf("readI16 = %#x, want %#x", got, 0x1234)
	}
	if got := readI16([]byte{0x34}); got != 0 {
		t.Fatalf("readI16 short input = %#x, want 0", got)
	}

	if got := readI32([]byte{0x78, 0x56, 0x34, 0x12}); got != 0x12345678 {
		t.Fatalf("readI32 = %#x, want %#x", got, 0x12345678)
	}
	if got := readI32([]byte{0x78, 0x56, 0x34}); got != 0 {
		t.Fatalf("readI32 short input = %#x, want 0", got)
	}
}
