package main

import "testing"

func TestSettingsToMap(t *testing.T) {
	t.Parallel()

	out := settingsToMap([]SyncSetting{
		{Path: "A", Value: "1"},
		{Path: "B", Value: "2"},
		{Path: "A", Value: "3"},
	})

	if len(out) != 2 {
		t.Fatalf("settingsToMap len = %d, want 2", len(out))
	}
	if out["A"] != "3" {
		t.Fatalf("settingsToMap overwrite = %q, want 3", out["A"])
	}
	if out["B"] != "2" {
		t.Fatalf("settingsToMap B = %q, want 2", out["B"])
	}

	if empty := settingsToMap(nil); len(empty) != 0 {
		t.Fatalf("settingsToMap(nil) = %#v, want empty", empty)
	}
}
