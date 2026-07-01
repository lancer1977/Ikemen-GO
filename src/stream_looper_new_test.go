package main

import (
	"io"
	"testing"
)

type fakeStreamSeeker struct {
	length   int
	position int
}

func (f *fakeStreamSeeker) Stream(samples [][2]float64) (int, bool) { return 0, false }
func (f *fakeStreamSeeker) Err() error                              { return nil }
func (f *fakeStreamSeeker) Len() int                                { return f.length }
func (f *fakeStreamSeeker) Position() int                           { return f.position }
func (f *fakeStreamSeeker) Seek(p int) error                        { f.position = p; return nil }

func TestNewStreamLooper(t *testing.T) {
	base := &fakeStreamSeeker{length: 100}

	sl := newStreamLooper(base, 2, -5, 200).(*StreamLooper)
	if sl.loopstart != 0 || sl.loopend != 100 {
		t.Fatalf("newStreamLooper should clamp loop points: %#v", sl)
	}

	sl = newStreamLooper(base, 2, 10, 5).(*StreamLooper)
	if sl.loopstart != 10 || sl.loopend != 100 {
		t.Fatalf("newStreamLooper should widen invalid loop end: %#v", sl)
	}

	if _, ok := newStreamLooper(base, 1, 0, 0).(io.Seeker); ok {
		t.Fatal("unexpected interface assertion")
	}
}
