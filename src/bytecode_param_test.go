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
	// DEFECT: afterImage_redirectid and hitDef_attr are the same byte (26), so
	// isHitDefParam cannot tell them apart. afterImage_redirectid is declared
	// after afterImage_last inside the same const block, taking the value the
	// next namespace computes as its own first parameter. palFX_redirectid and
	// bgPalFX_id collide at 11 the same way. The Projectile controller orders
	// isHitDefParam before isAfterImageParam, so a colliding id routes to the
	// hitdef handler as its attack attribute.
	// Tracked as lancer1977/Ikemen-GO#20.
	if afterImage_redirectid != hitDef_attr {
		t.Fatalf("afterImage_redirectid (%d) and hitDef_attr (%d) no longer collide; "+
			"#20 is fixed, so this should assert exclusion again",
			afterImage_redirectid, hitDef_attr)
	}
	if !isHitDefParam(afterImage_redirectid) {
		t.Fatal("expected the colliding id to be misclassified as a hitDef param")
	}
}
