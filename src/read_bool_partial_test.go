package main

import "testing"

func TestIniSectionReadBool_PreservesUnusedOutputsAndSkipsBlankFields(t *testing.T) {
	is := IniSection{
		"flags": "1, , 0",
	}
	a, b, c, d := false, true, true, true
	if !is.ReadBool("flags", &a, &b, &c, &d) {
		t.Fatal("expected ReadBool to report success")
	}
	if !a || !b || c || !d {
		t.Fatalf("ReadBool outputs = %v %v %v %v, want true true false true", a, b, c, d)
	}
}
