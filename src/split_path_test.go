package main

import "testing"

func TestSplitPath_NormalizesSeparatorsAndSplitsAtLastSlashAgain(t *testing.T) {
	dir, file := SplitPath(`C:\games\ikemen\data\chars\ryu.def`)
	if dir != "C:/games/ikemen/data/chars/" {
		t.Fatalf("dir = %q, want %q", dir, "C:/games/ikemen/data/chars/")
	}
	if file != "ryu.def" {
		t.Fatalf("file = %q, want %q", file, "ryu.def")
	}

	dir, file = SplitPath("select.def")
	if dir != "" || file != "select.def" {
		t.Fatalf("SplitPath(no dir) = %q, %q, want empty, select.def", dir, file)
	}

	dir, file = SplitPath("/tmp/")
	if dir != "/tmp/" || file != "" {
		t.Fatalf("SplitPath(trailing slash) = %q, %q, want /tmp/, empty", dir, file)
	}

	dir, file = SplitPath("")
	if dir != "" || file != "" {
		t.Fatalf("SplitPath(empty) = %q, %q, want empty, empty", dir, file)
	}

	dir, file = SplitPath(`folder\`)
	if dir != "folder/" || file != "" {
		t.Fatalf("SplitPath(backslash trailing) = %q, %q, want folder/, empty", dir, file)
	}
}
