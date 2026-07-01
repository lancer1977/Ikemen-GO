package main

import "testing"

func TestEncodeInputs(t *testing.T) {
	t.Parallel()

	got := encodeInputs(InputBits(0x1234))
	want := writeI16(0x1234)
	if len(got) != len(want) {
		t.Fatalf("encodeInputs len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("encodeInputs[%d] = %#x, want %#x", i, got[i], want[i])
		}
	}
}
