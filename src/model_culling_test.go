package main

import (
	"testing"

	mgl "github.com/go-gl/mathgl/mgl32"
)

func TestIsCulled(t *testing.T) {
	t.Parallel()

	id := mgl.Ident4()
	if got := isCulled(id, BoundingBox{min: [3]float32{-1, -1, -1}, max: [3]float32{1, 1, 1}}); got {
		t.Fatal("identity matrix should not cull an in-frustum box")
	}

	if got := isCulled(id, BoundingBox{min: [3]float32{10, 10, 10}, max: [3]float32{11, 11, 11}}); !got {
		t.Fatal("identity matrix should cull a far-away box")
	}
}
