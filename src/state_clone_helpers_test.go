package main

import "testing"

func TestCopySlice(t *testing.T) {
	src := []int{1, 2, 3}
	dst := []int{9, 9}

	CopySlice(&src, &dst)

	want := []int{1, 2, 3}
	if len(dst) != len(want) {
		t.Fatalf("CopySlice len = %d, want %d", len(dst), len(want))
	}
	for i := range want {
		if dst[i] != want[i] {
			t.Fatalf("CopySlice[%d] = %d, want %d", i, dst[i], want[i])
		}
	}
	if &dst[0] == &src[0] {
		t.Fatal("CopySlice should copy values into destination storage")
	}
}

func TestCopyMap(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2}
	dst := map[string]int{"b": 9, "c": 3}

	CopyMap(&src, &dst)

	if len(dst) != 2 || dst["a"] != 1 || dst["b"] != 2 {
		t.Fatalf("CopyMap = %#v", dst)
	}
	if _, ok := dst["c"]; ok {
		t.Fatal("CopyMap should remove keys absent from source")
	}
}
