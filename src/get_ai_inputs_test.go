package main

import "testing"

func TestGetAIInputs(t *testing.T) {
	// Not parallel: this test mutates the package-level sys, which every test in
	// this package shares.

	oldSys := sys
	defer func() { sys = oldSys }()

	sys.aiInput[0] = AiInput{
		dir:  7,
		dirt: 1,
		at:   1,
		zt:   1,
		wt:   1,
	}

	got := getAIInputs(0)
	if len(got) != 4 {
		t.Fatalf("getAIInputs length = %d, want 4", len(got))
	}
	// SetInputAI packs bits as: U(0)|D(1)|L(2)|R(3)|a(4)|b(5)|c(6)|x(7)|y(8)|z(9)|s(10)|d(11)|w(12)|m(13)
	// With dir=7 and dirt=1: U() returns true (dir 7,0,1 means up component), L() also returns true
	// (dir 5,6,7 means left component). With at=1, zt=1, wt=1: a, z, w return true.
	// Expected bits: U(0), L(2), a(4), z(9), w(12) = 0x1215
	const expectedBits = 1<<0 | 1<<2 | 1<<4 | 1<<9 | 1<<12
	if val := readI32(got); val != expectedBits {
		t.Fatalf("getAIInputs packed bits = %#x, want %#x", val, expectedBits)
	}
}
