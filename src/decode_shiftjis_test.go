package main

import (
	"testing"

	"golang.org/x/text/encoding/japanese"
)

func TestDecodeShiftJIS_PreservesUtf8AndDecodesShiftJISBytes(t *testing.T) {
	if got := decodeShiftJIS("hello"); got != "hello" {
		t.Fatalf("decodeShiftJIS(utf8) = %q, want %q", got, "hello")
	}

	encoded, err := japanese.ShiftJIS.NewEncoder().String("あ")
	if err != nil {
		t.Fatalf("encode shift-jis: %v", err)
	}
	if got := decodeShiftJIS(encoded); got != "あ" {
		t.Fatalf("decodeShiftJIS(shift-jis) = %q, want %q", got, "あ")
	}
}

// decodeShiftJIS now correctly detects when the ShiftJIS decoder substitutes
// U+FFFD for bytes it cannot map, and falls back to the original input.
func TestDecodeShiftJIS_SubstitutesReplacementCharForUndecodableBytes(t *testing.T) {
	// 0x81 is a lead byte with no following trail byte, so it cannot form a
	// valid Shift_JIS sequence, and it is not valid UTF-8 either.
	input := string([]byte{0x81})
	got := decodeShiftJIS(input)

	// The function should now return the original input unchanged when the
	// decoder would have substituted replacement characters.
	if got != input {
		t.Fatalf("decodeShiftJIS(0x81) = %q, want the original input %q", got, input)
	}
}
