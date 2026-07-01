package main

import "testing"

func TestFloatAndSampleConversions(t *testing.T) {
	t.Parallel()

	if got := floatToS16(1.5); got != 32767 {
		t.Fatalf("floatToS16(1.5) = %d, want 32767", got)
	}
	if got := floatToS16(-1.5); got != -32767 {
		t.Fatalf("floatToS16(-1.5) = %d, want -32767", got)
	}
	if got := floatToS16(0.5); got != 16383 {
		t.Fatalf("floatToS16(0.5) = %d, want 16383", got)
	}

	if got := convertI16toI8(32767); got != 127 {
		t.Fatalf("convertI16toI8(32767) = %d, want 127", got)
	}
	if got := convertI16toI8(-32768); got != -128 {
		t.Fatalf("convertI16toI8(-32768) = %d, want -128", got)
	}
	if got := convertI16toI8(0); got != 0 {
		t.Fatalf("convertI16toI8(0) = %d, want 0", got)
	}
}
