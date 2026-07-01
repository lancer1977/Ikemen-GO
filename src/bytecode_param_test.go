package main

import "testing"

func TestBytecodeParamClassifiers(t *testing.T) {
	t.Parallel()

	if !isPalFXParam(palFX_time) || !isPalFXParam(palFX_last) {
		t.Fatal("expected palFX range to be inclusive")
	}
	if isPalFXParam(palFX_redirectid) {
		t.Fatal("expected palFX redirect sentinel to be excluded")
	}

	if !isAfterImageParam(afterImage_time) || !isAfterImageParam(afterImage_last) {
		t.Fatal("expected afterImage range to be inclusive")
	}
	if isAfterImageParam(afterImage_redirectid) {
		t.Fatal("expected afterImage redirect sentinel to be excluded")
	}

	if !isHitDefParam(hitDef_attr) || !isHitDefParam(hitDef_last) {
		t.Fatal("expected hitDef range to be inclusive")
	}
	if !isHitDefParam(palFX_time) {
		t.Fatal("expected hitDef to include palFX params")
	}
	if isHitDefParam(afterImage_redirectid) {
		t.Fatal("expected unrelated sentinel to be excluded")
	}
}
