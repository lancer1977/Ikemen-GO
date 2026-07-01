package main

import "testing"

func TestGetDefaultDefPathInZip(t *testing.T) {
	t.Parallel()

	path1, path2 := getDefaultDefPathInZip(`/tmp/My Pack.ZIP`)
	if path1 != "my pack.def" {
		t.Fatalf("path1 = %q, want %q", path1, "my pack.def")
	}
	if path2 != "my pack/my pack.def" {
		t.Fatalf("path2 = %q, want %q", path2, "my pack/my pack.def")
	}

	path1, path2 = getDefaultDefPathInZip(`/tmp/bonus.zip`)
	if path1 != "bonus.def" || path2 != "bonus/bonus.def" {
		t.Fatalf("unexpected default def paths: %q %q", path1, path2)
	}
}
