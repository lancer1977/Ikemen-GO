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

func TestDecodeShiftJIS_FallsBackToOriginalOnDecodeFailure(t *testing.T) {
	input := string([]byte{0x81})
	if got := decodeShiftJIS(input); got != input {
		t.Fatalf("decodeShiftJIS(fallback) = %q, want original input %q", got, input)
	}
}
