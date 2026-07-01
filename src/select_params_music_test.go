package main

import "testing"

func TestSplitMusicParamValue(t *testing.T) {
	t.Parallel()

	path, extras := splitMusicParamValue(" bgm/title.ogg  80  100  200 ")
	if path != "bgm/title.ogg" {
		t.Fatalf("unexpected path: %q", path)
	}
	want := []string{"80", "100", "200"}
	if len(extras) != len(want) {
		t.Fatalf("unexpected extras length: %v", extras)
	}
	for i, v := range want {
		if extras[i] != v {
			t.Fatalf("extras[%d] = %q, want %q", i, extras[i], v)
		}
	}

	path, extras = splitMusicParamValue("loopless_track")
	if path != "loopless_track" || extras != nil {
		t.Fatalf("unexpected fallback parse: path=%q extras=%v", path, extras)
	}

	path, extras = splitMusicParamValue("   ")
	if path != "" || extras != nil {
		t.Fatalf("unexpected blank parse: path=%q extras=%v", path, extras)
	}
}

func TestExpandMusicKV(t *testing.T) {
	t.Parallel()

	got := expandMusicKV("music", " intro.ogg 75 12 34 ")
	want := []string{
		"music=intro.ogg",
		"bgmvolume=75",
		"bgmloopstart=12",
		"bgmloopend=34",
	}
	if len(got) != len(want) {
		t.Fatalf("expandMusicKV length = %d, want %d (%v)", len(got), len(want), got)
	}
	for i, v := range want {
		if got[i] != v {
			t.Fatalf("expandMusicKV[%d] = %q, want %q", i, got[i], v)
		}
	}

	if got := expandMusicKV("", "intro.ogg"); got != nil {
		t.Fatalf("expected empty key to return nil, got %v", got)
	}
}
