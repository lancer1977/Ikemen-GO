package main

import (
	"bytes"
	"testing"
)

func TestReplayInputIO(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	wantBits := InputBits(0x1234)
	wantAxes := [6]int8{-1, 0, 1, 2, 3, 4}

	if err := writeReplayInput(&buf, wantBits, wantAxes); err != nil {
		t.Fatalf("writeReplayInput returned error: %v", err)
	}

	var gotBits InputBits
	var gotAxes [6]int8
	if err := readReplayInput(&buf, &gotBits, &gotAxes); err != nil {
		t.Fatalf("readReplayInput returned error: %v", err)
	}
	if gotBits != wantBits {
		t.Fatalf("got bits = %#x, want %#x", gotBits, wantBits)
	}
	if gotAxes != wantAxes {
		t.Fatalf("got axes = %#v, want %#v", gotAxes, wantAxes)
	}
}

func TestReadReplayInputTruncated(t *testing.T) {
	t.Parallel()

	var gotBits InputBits
	var gotAxes [6]int8
	err := readReplayInput(bytes.NewReader([]byte{0x34, 0x12, 1, 2, 3}), &gotBits, &gotAxes)
	if err == nil {
		t.Fatal("readReplayInput on truncated buffer returned nil error")
	}
}
