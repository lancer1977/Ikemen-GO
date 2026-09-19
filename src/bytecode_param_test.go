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

	// Verify parameter ID collisions are fixed (lancer1977/Ikemen-GO#20)
	if palFX_redirectid == bgPalFX_id {
		t.Fatalf("palFX_redirectid (%d) and bgPalFX_id (%d) collide; #20 is not fixed",
			palFX_redirectid, bgPalFX_id)
	}
	if afterImage_redirectid == hitDef_attr {
		t.Fatalf("afterImage_redirectid (%d) and hitDef_attr (%d) collide; #20 is not fixed",
			afterImage_redirectid, hitDef_attr)
	}
}
