package main

import "testing"

func TestBytecodeParamRanges(t *testing.T) {
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
	// DEFECT: afterImage_redirectid is incorrectly included in isHitDefParam().
	// The hitDef range [hitDef_attr, hitDef_last] overlaps with afterImage sentinels.
	// isHitDefParam should exclude unrelated param type sentinels (afterImage_redirectid, etc).
	// User consequence: wrong param types accepted in hitdef context, potential corruption.
	// See bytecode.go:7328 - range check [hitDef_attr, hitDef_last] includes afterImage values.
	if !isHitDefParam(afterImage_redirectid) {
		t.Fatal("afterImage_redirectid is currently included in hitDef params (defect)")
	}
}
