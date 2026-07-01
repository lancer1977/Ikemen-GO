package main

import "testing"

func TestLayoutDrawFaceSpriteAndAnim_NoOpOnMismatchedLayerOrEmptyAnimation(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.scrrect = [4]int32{0, 0, 640, 480}
	sys.fightScreen.scale = 1
	sys.fightScreen.fnt_scale = 1

	l := newLayout(3)
	l.DrawFaceSprite(10, 20, 4, nil, nil, 1, nil)
	l.DrawAnim(nil, 10, 20, 1, 1, 1, 4, nil, nil)

	// Match the layer but keep the animation empty so DrawAnim remains a no-op.
	l.DrawFaceSprite(10, 20, 3, nil, nil, 1, nil)
	l.DrawAnim(nil, 10, 20, 1, 1, 1, 3, &Animation{}, nil)
}
