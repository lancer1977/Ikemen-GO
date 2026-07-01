package main

import "testing"

func TestIniSectionReadF32_PreservesUnusedOutputsAndSkipsBlankFields(t *testing.T) {
	is := IniSection{
		"coords": "1.5, , -3.25",
	}
	a, b, c, d := float32(1), float32(2), float32(3), float32(4)
	if !is.ReadF32("coords", &a, &b, &c, &d) {
		t.Fatal("expected ReadF32 to report success")
	}
	if a != 1.5 || b != 2 || c != -3.25 || d != 4 {
		t.Fatalf("ReadF32 outputs = %v %v %v %v, want 1.5 2 -3.25 4", a, b, c, d)
	}
}
