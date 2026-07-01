package main

import "testing"

func TestIniSectionReadI32_PreservesUnusedOutputsAndSkipsBlankFields(t *testing.T) {
	is := IniSection{
		"coords": "10, , -30",
	}
	a, b, c, d := int32(1), int32(2), int32(3), int32(4)
	if !is.ReadI32("coords", &a, &b, &c, &d) {
		t.Fatal("expected ReadI32 to report success")
	}
	if a != 10 || b != 2 || c != -30 || d != 4 {
		t.Fatalf("ReadI32 outputs = %v %v %v %v, want 10 2 -30 4", a, b, c, d)
	}
}
