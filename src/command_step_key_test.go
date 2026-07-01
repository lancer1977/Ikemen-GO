package main

import "testing"

func TestCommandStepKeyClassifiers(t *testing.T) {
	t.Parallel()

	if !((CommandStepKey{key: CK_U}).IsDirectionPress()) {
		t.Fatal("expected CK_U to be a direction press")
	}
	if (CommandStepKey{key: CK_U, tilde: true}).IsDirectionPress() {
		t.Fatal("expected tilde direction to not be a press")
	}
	if !((CommandStepKey{key: CK_U, tilde: true}).IsDirectionRelease()) {
		t.Fatal("expected tilde direction to be a release")
	}
	if (CommandStepKey{key: CK_a}).IsDirectionPress() {
		t.Fatal("expected button key to not be a direction press")
	}
	if !((CommandStepKey{key: CK_a}).IsButtonPress()) {
		t.Fatal("expected CK_a to be a button press")
	}
	if (CommandStepKey{key: CK_a, tilde: true}).IsButtonPress() {
		t.Fatal("expected tilde button to not be a press")
	}
	if !((CommandStepKey{key: CK_a, tilde: true}).IsButtonRelease()) {
		t.Fatal("expected tilde button to be a release")
	}
}
