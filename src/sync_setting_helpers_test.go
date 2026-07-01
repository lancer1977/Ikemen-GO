package main

import "testing"

func TestCloneSyncSettings(t *testing.T) {
	in := []SyncSetting{{Path: "Video.GameWidth", Value: "1280"}}
	out := cloneSyncSettings(in)

	if len(out) != 1 || out[0] != in[0] {
		t.Fatalf("cloneSyncSettings = %#v, want %#v", out, in)
	}

	out[0].Value = "640"
	if in[0].Value != "1280" {
		t.Fatalf("cloneSyncSettings returned alias to input slice")
	}
}

func TestVerifySyncSettings(t *testing.T) {
	issues := verifySyncSettings(
		[]SyncSetting{
			{Path: "A", Value: "1"},
			{Path: "B", Value: "2"},
		},
		[]SyncSetting{
			{Path: "A", Value: "1"},
			{Path: "B", Value: "3"},
			{Path: "C", Value: "4"},
		},
	)

	want := []string{
		"B mismatch (expected=2 actual=3)",
		"C unexpected (actual=4)",
	}
	if len(issues) != len(want) {
		t.Fatalf("verifySyncSettings issues = %#v, want %#v", issues, want)
	}
	for i := range want {
		if issues[i] != want[i] {
			t.Fatalf("verifySyncSettings issues[%d] = %q, want %q", i, issues[i], want[i])
		}
	}
}
