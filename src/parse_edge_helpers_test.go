package main

import "testing"

func TestParseHelpers_HandleEmptyAndSignedOnlyInputs(t *testing.T) {
	if got := Atof(""); got != 0 {
		t.Fatalf("Atof(empty) = %v, want 0", got)
	}
	if got := Atof("7.5xyz"); got != 7.5 {
		t.Fatalf("Atof(trailing junk) = %v, want 7.5", got)
	}
	if IsInt("+") || IsInt("-") {
		t.Fatal("IsInt should reject signed-only strings")
	}
	if !IsInt("  +12 ") {
		t.Fatal("IsInt should accept signed numeric strings")
	}
}
