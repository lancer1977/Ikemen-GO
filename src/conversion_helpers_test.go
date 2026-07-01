package main

import "testing"

func TestBtoiAndNarrowIntegerConversions_HandleBounds(t *testing.T) {
	if got := Btoi(true); got != 1 {
		t.Fatalf("unexpected true conversion: %v", got)
	}
	if got := Btoi(false); got != 0 {
		t.Fatalf("unexpected false conversion: %v", got)
	}

	if got := I32ToI16(12345); got != 12345 {
		t.Fatalf("unexpected passthrough i16 conversion: %v", got)
	}
	if got := I32ToI16(40000); got != 32767 {
		t.Fatalf("unexpected i16 upper clamp: %v", got)
	}
	if got := I32ToI16(-40000); got != -32768 {
		t.Fatalf("unexpected i16 lower clamp: %v", got)
	}

	if got := I32ToU16(12345); got != 12345 {
		t.Fatalf("unexpected passthrough u16 conversion: %v", got)
	}
	if got := I32ToU16(70000); got != 65535 {
		t.Fatalf("unexpected u16 upper clamp: %v", got)
	}
	if got := I32ToU16(-1); got != 0 {
		t.Fatalf("unexpected u16 lower clamp: %v", got)
	}
}

func TestRoundFloat_RoundsToRequestedPrecision(t *testing.T) {
	if got := RoundFloat(12.3456, 2); got != 12.35 {
		t.Fatalf("unexpected rounded value: %v", got)
	}
	if got := RoundFloat(12.344, 2); got != 12.34 {
		t.Fatalf("unexpected rounded value: %v", got)
	}
}
