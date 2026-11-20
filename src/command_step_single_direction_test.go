package main

import "testing"

func TestCommandStepIsSingleDirection(t *testing.T) {
	if !(&CommandStep{keys: []CommandStepKey{{key: CK_U}}}).IsSingleDirection() {
		t.Fatal("expected single direction command step")
	}
	if (&CommandStep{keys: []CommandStepKey{{key: CK_U}, {key: CK_D}}}).IsSingleDirection() {
		t.Fatal("expected multi-key step to be false")
	}
	if (&CommandStep{keys: []CommandStepKey{{key: CK_a}}}).IsSingleDirection() {
		t.Fatal("expected button step to be false")
	}
	if (&CommandStep{keys: nil}).IsSingleDirection() {
		t.Fatal("expected empty step to be false")
	}
}
