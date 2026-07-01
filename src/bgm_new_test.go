package main

import "testing"

func TestNewBgm(t *testing.T) {
	t.Parallel()

	bgm := newBgm()
	if bgm == nil {
		t.Fatal("newBgm returned nil")
	}
	if bgm.filename != "" || bgm.loop != 0 || bgm.bgmVolume != 0 {
		t.Fatalf("unexpected value fields: %#v", bgm)
	}
	if bgm.ctrl != nil || bgm.volctrl != nil || bgm.streamer != nil || bgm.cancel != nil {
		t.Fatalf("unexpected initialized pointers: %#v", bgm)
	}
	if bgm.pauseVolumeApplied || bgm.volRestore != 0 || bgm.startPos != 0 {
		t.Fatalf("unexpected flags/counters: %#v", bgm)
	}
}
