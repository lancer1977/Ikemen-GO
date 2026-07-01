package main

import "testing"

func TestFlagLStrAndAttrLStrEncodeExpectedSymbols(t *testing.T) {
	if got := flagLStr(int32(HF_H | HF_L | HF_PLS)); got != "HL+" {
		t.Fatalf("flagLStr = %q, want HL+", got)
	}

	if got := attrLStr(0); got != "" {
		t.Fatalf("attrLStr(0) = %q, want empty", got)
	}

	// ST_S + AT_AN + AT_AA encodes as "S, NA"
	attr := int32(ST_S | AT_AN | AT_AA)
	if got := attrLStr(attr); got != "S, NA" {
		t.Fatalf("attrLStr = %q, want S, NA", got)
	}
}
