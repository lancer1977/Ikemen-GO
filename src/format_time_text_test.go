package main

import "testing"

func TestFormatTimeText(t *testing.T) {
	t.Parallel()

	got := formatTimeText("%h:%m:%s.%x", 3661.234)
	if got != "1:01:01.23" {
		t.Fatalf("formatTimeText = %q, want %q", got, "1:01:01.23")
	}

	got = formatTimeText("time %h", 59.9)
	if got != "time 0" {
		t.Fatalf("formatTimeText hour-only = %q, want %q", got, "time 0")
	}
}
