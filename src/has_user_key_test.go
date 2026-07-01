package main

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestHasUserKey(t *testing.T) {
	t.Parallel()

	f := ini.Empty()
	sec, err := f.NewSection("select")
	if err != nil {
		t.Fatalf("NewSection: %v", err)
	}
	sec.NewKey("title", "yes")

	if !hasUserKey(f, "select", "title") {
		t.Fatal("expected existing key to be found")
	}
	if hasUserKey(f, "select", "missing") {
		t.Fatal("expected missing key to be absent")
	}
	if hasUserKey(f, "other", "title") {
		t.Fatal("expected missing section to be absent")
	}
}
