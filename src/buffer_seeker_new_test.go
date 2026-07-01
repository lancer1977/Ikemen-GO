package main

import (
	"testing"

	"github.com/faiface/beep"
)

func TestNewBufferSeeker(t *testing.T) {
	t.Parallel()

	buf := beep.NewBuffer(beep.Format{SampleRate: 44100, NumChannels: 2, Precision: 2})
	buf.Append(beep.Silence(4))

	bs := newBufferSeeker(buf)
	if bs == nil {
		t.Fatal("newBufferSeeker returned nil")
	}
	if bs.buf != buf {
		t.Fatalf("buffer mismatch: got %#v want %#v", bs.buf, buf)
	}
	if got := bs.Position(); got != 0 {
		t.Fatalf("initial position = %d, want 0", got)
	}
	if got := bs.Len(); got != 4 {
		t.Fatalf("Len = %d, want 4", got)
	}
	if err := bs.Seek(10); err != nil {
		t.Fatalf("Seek(10) = %v", err)
	}
	if got := bs.Position(); got != 4 {
		t.Fatalf("clamped position = %d, want 4", got)
	}
	if err := bs.Seek(-5); err != nil {
		t.Fatalf("Seek(-5) = %v", err)
	}
	if got := bs.Position(); got != 0 {
		t.Fatalf("negative seek clamp = %d, want 0", got)
	}
}
