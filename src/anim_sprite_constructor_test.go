package main

import "testing"

func TestAnimSpriteConstructors(t *testing.T) {
	t.Parallel()

	sd := newSpriteData()
	if sd == nil {
		t.Fatal("newSpriteData returned nil")
	}
	if sd.trans != TT_default || sd.alpha != [2]int32{255, 0} || sd.airOffsetFix != [2]float32{1, 1} {
		t.Fatalf("unexpected sprite data defaults: %#v", sd)
	}

	ss := newShadowSprite()
	if ss == nil {
		t.Fatal("newShadowSprite returned nil")
	}
	if ss.shadowColor != -1 || ss.shadowAlpha != 255 || ss.shadowIntensity != -1 || !ss.shadowKeeptransform || ss.shadowProjection != -1 {
		t.Fatalf("unexpected shadow sprite defaults: %#v", ss)
	}

	rs := newReflectionSprite()
	if rs == nil {
		t.Fatal("newReflectionSprite returned nil")
	}
	if rs.reflectColor != -1 || rs.reflectIntensity != -1 || !rs.reflectKeeptransform || rs.reflectProjection != -1 {
		t.Fatalf("unexpected reflection sprite defaults: %#v", rs)
	}
}
