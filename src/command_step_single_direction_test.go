package main

import "testing"

func TestCommandStepIsSingleDirection(t *testing.T) {
	if !(CommandStep{key: []CommandKey{CK_U}}).IsSingleDirection() {
		t.Fatal("expected single direction command step")
	}
	if (CommandStep{key: []CommandKey{CK_U, CK_D}}).IsSingleDirection() {
		t.Fatal("expected multi-key step to be false")
	}
	if (CommandStep{key: []CommandKey{CK_a}}).IsSingleDirection() {
		t.Fatal("expected button step to be false")
	}
	if (CommandStep{key: nil}).IsSingleDirection() {
		t.Fatal("expected empty step to be false")
	}
}
