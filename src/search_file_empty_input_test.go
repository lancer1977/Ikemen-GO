package main

import "testing"

func TestSearchFile_ReturnsEmptyForBlankInputAndIgnoresBlankDefaultDirs(t *testing.T) {
	if got := SearchFile("   ", []string{"  "}, " ", "\t"); got != "" {
		t.Fatalf("SearchFile(blank) = %q, want empty", got)
	}
}
