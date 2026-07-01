package main

import (
	"path/filepath"
	"testing"
)

func TestLowercaseNoExtension_StripsExtensionAndLowercasesBaseNameAgain(t *testing.T) {
	if got := LowercaseNoExtension(filepath.Join("data", "Chars", "RYU.CNS")); got != "ryu" {
		t.Fatalf("LowercaseNoExtension() = %q, want %q", got, "ryu")
	}
	if got := LowercaseNoExtension(filepath.Join("data", "Chars", "Stage")); got != "stage" {
		t.Fatalf("LowercaseNoExtension() = %q, want %q", got, "stage")
	}
}
