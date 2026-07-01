package main

import "testing"

func TestNewBgMusic(t *testing.T) {
	bg := newBgMusic()
	if bg == nil {
		t.Fatal("expected newBgMusic to allocate")
	}
	if bg.bgmloop != 1 || bg.bgmvolume != 100 || bg.bgmfreqmul != 1 || bg.bgmloopcount != -1 {
		t.Fatalf("unexpected defaults: %#v", bg)
	}
}
