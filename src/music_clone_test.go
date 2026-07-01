package main

import (
	"arena"
	"testing"
)

func TestCloneMusicMapShallow(t *testing.T) {
	t.Parallel()

	a := arena.NewArena()
	defer a.Free()

	first := &bgMusic{bgmusic: "a", selected: true}
	second := &bgMusic{bgmusic: "b"}
	src := Music{
		"battle": []*bgMusic{first, second},
		"empty":  nil,
	}

	dst := cloneMusicMapShallow(a, src)
	if len(dst) != len(src) {
		t.Fatalf("cloneMusicMapShallow len = %d, want %d", len(dst), len(src))
	}
	if dst["battle"] == nil || len(dst["battle"]) != 2 {
		t.Fatalf("cloneMusicMapShallow battle = %#v", dst["battle"])
	}
	if dst["battle"][0] != first || dst["battle"][1] != second {
		t.Fatal("expected shallow clone to preserve bgMusic pointers")
	}
	if &dst["battle"][0] == &src["battle"][0] {
		t.Fatal("expected battle slice to be copied, not aliased")
	}
	if dst["empty"] != nil {
		t.Fatalf("expected nil slice to stay nil, got %#v", dst["empty"])
	}

	dst["battle"][0] = &bgMusic{bgmusic: "changed"}
	if src["battle"][0] != first {
		t.Fatal("mutating cloned slice should not affect source slice")
	}
}
