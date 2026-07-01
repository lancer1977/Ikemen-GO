package main

import "testing"

func TestRollbackEncodingHelpers(t *testing.T) {
	t.Parallel()

	if got := readI16(writeI16(0x1234)); got != 0x1234 {
		t.Fatalf("readI16(writeI16()) = %#x, want 0x1234", got)
	}
	if got := readI32(writeI32(0x12345678)); got != 0x12345678 {
		t.Fatalf("readI32(writeI32()) = %#x, want 0x12345678", got)
	}

	inputs, analog := decodeInputs([][]byte{
		{0x34, 0x12, 1, 2, 3, 4, 5, 6},
		{0x78, 0x56},
	})
	if len(inputs) != 2 || len(analog) != 2 {
		t.Fatalf("unexpected decode lengths: %d %d", len(inputs), len(analog))
	}
	if inputs[0] != InputBits(0x1234) || inputs[1] != InputBits(0x5678) {
		t.Fatalf("unexpected decoded inputs: %#v", inputs)
	}
	if analog[0] != [6]int8{1, 2, 3, 4, 5, 6} {
		t.Fatalf("unexpected analog decode: %#v", analog[0])
	}
	if analog[1] != [6]int8{} {
		t.Fatalf("expected short buffer analogs to zero-fill, got %#v", analog[1])
	}
}
