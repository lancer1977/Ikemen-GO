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

// decodeShiftJIS has an error branch that logs a warning and returns the input
// unchanged, but it is unreachable: japanese.ShiftJIS's decoder substitutes
// U+FFFD for bytes it cannot map rather than returning an error, so
// transform.Bytes never reports one. This test pins the behaviour that actually
// happens today so the dead branch is visible rather than assumed live.
func TestDecodeShiftJIS_SubstitutesReplacementCharForUndecodableBytes(t *testing.T) {
	// 0x81 is a lead byte with no following trail byte, so it cannot form a
	// valid Shift_JIS sequence, and it is not valid UTF-8 either.
	input := string([]byte{0x81})
	got := decodeShiftJIS(input)

	if got == input {
		t.Fatalf("decodeShiftJIS returned the input unchanged, which would mean " +
			"the error branch became reachable; update the comment above and " +
			"the engine issue tracking it")
	}
	if got != "\uFFFD" {
		t.Fatalf("decodeShiftJIS(0x81) = %q, want the replacement character", got)
	}
}
