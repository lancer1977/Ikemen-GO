package main

import (
	"testing"
)

type fakeStreamSeekerLooper struct {
	length   int
	position int
}

func (f *fakeStreamSeekerLooper) Stream(samples [][2]float64) (int, bool) { return 0, false }
func (f *fakeStreamSeekerLooper) Err() error                              { return nil }
func (f *fakeStreamSeekerLooper) Len() int                                { return f.length }
func (f *fakeStreamSeekerLooper) Position() int                           { return f.position }
func (f *fakeStreamSeekerLooper) Seek(p int) error                        { f.position = p; return nil }

func TestNewStreamLooper(t *testing.T) {
	base := &fakeStreamSeekerLooper{length: 100}

	sl := newStreamLooper(base, 2, -5, 200).(*StreamLooper)
	if sl.loopstart != 0 || sl.loopend != 100 {
		t.Fatalf("newStreamLooper should clamp loop points: %#v", sl)
	}

	sl = newStreamLooper(base, 2, 10, 5).(*StreamLooper)
	if sl.loopstart != 10 || sl.loopend != 100 {
		t.Fatalf("newStreamLooper should widen invalid loop end: %#v", sl)
	}

	// A zero loopend is invalid and must widen to the full stream length.
	sl = newStreamLooper(base, 1, 0, 0).(*StreamLooper)
	if sl.loopcount != 1 || sl.loopstart != 0 || sl.loopend != 100 {
		t.Fatalf("newStreamLooper should widen a zero loop end: %#v", sl)
	}
}
