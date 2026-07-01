package main

import "testing"

func TestFloatToS16(t *testing.T) {
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
}
