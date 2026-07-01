package main

import "testing"

func TestLoadText_ReturnsErrorForMissingFile(t *testing.T) {
	got, err := LoadText("definitely_missing_file_ikemen.txt")
	if err == nil {
		t.Fatal("expected error")
	}
	if got != "" {
		t.Fatalf("LoadText() text = %q, want empty on error", got)
	}
}
