package main

import "testing"

func TestBytecodeParamClassifiers(t *testing.T) {
	t.Parallel()

	if !isPalFXParam(0) || !isPalFXParam(byte(palFX_last)) {
		t.Fatal("expected palfx params to accept low and high valid ids")
	}
	if isPalFXParam(byte(afterImage_time)) {
		t.Fatal("expected palfx params to reject afterimage ids")
	}

	if !isAfterImageParam(byte(afterImage_time)) || !isAfterImageParam(byte(afterImage_last)) {
		t.Fatal("expected afterimage params to accept low and high valid ids")
	}
	if isAfterImageParam(0) {
		t.Fatal("expected afterimage params to reject palfx ids")
	}

	if !isHitDefParam(0) || !isHitDefParam(byte(hitDef_last)) {
		t.Fatal("expected hitdef params to accept palFX ids and hitdef ids")
	}
	if isHitDefParam(byte(afterImage_time)) {
		t.Fatal("expected hitdef params to reject afterimage-only ids")
	}
}
