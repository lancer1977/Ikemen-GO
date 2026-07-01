package main

import (
	"testing"

	"github.com/faiface/beep"
)

type fakeStreamSeeker struct {
	pos int
	ln  int
}

func (f *fakeStreamSeeker) Stream(out [][2]float64) (int, bool) { return 0, false }
func (f *fakeStreamSeeker) Err() error                          { return nil }
func (f *fakeStreamSeeker) Len() int                            { return f.ln }
func (f *fakeStreamSeeker) Position() int                       { return f.pos }
func (f *fakeStreamSeeker) Seek(p int) error                    { f.pos = p; return nil }

var _ beep.StreamSeeker = (*fakeStreamSeeker)(nil)

func TestNewSwapSeeker(t *testing.T) {
	t.Parallel()

	ss := &fakeStreamSeeker{pos: 7, ln: 42}
	sw := newSwapSeeker(ss)
	if sw == nil {
		t.Fatal("newSwapSeeker returned nil")
	}
	if sw.ss != ss {
		t.Fatalf("newSwapSeeker stored %#v, want %#v", sw.ss, ss)
	}
}
