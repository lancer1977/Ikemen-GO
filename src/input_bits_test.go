package main

import "testing"

func TestInputBitsKeysToBitsAndBack(t *testing.T) {
	var bits InputBits
	keys := [14]bool{true, false, true, false, true, false, true, false, true, false, true, false, true, false}
	bits.KeysToBits(keys)

	if bits != InputBits(0x1555) {
		t.Fatalf("KeysToBits = %#x, want 0x1555", int16(bits))
	}

	got := bits.BitsToKeys()
	if got != keys {
		t.Fatalf("BitsToKeys = %#v, want %#v", got, keys)
	}
}
