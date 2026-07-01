package main

import "testing"

func TestGetAIInputs(t *testing.T) {
	t.Parallel()

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
	if val := readI32(got); val != 1<<0|1<<4|1<<9|1<<12 {
		t.Fatalf("getAIInputs packed bits = %#x, want %#x", val, 1<<0|1<<4|1<<9|1<<12)
	}
}
