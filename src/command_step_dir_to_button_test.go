package main

import "testing"

func TestCommandStepIsDirToButton(t *testing.T) {
	if (CommandStep{keys: []CommandStepKey{{key: CK_U, slash: true}}}).IsDirToButton(CommandStep{keys: []CommandStepKey{{key: CK_a}}}) {
		t.Fatal("expected slash-held next step to be false")
	}
	if (CommandStep{keys: []CommandStepKey{{key: CK_U}}}).IsDirToButton(CommandStep{keys: []CommandStepKey{{key: CK_U}, {key: CK_a}}}) {
		t.Fatal("expected shared key steps to be false")
	}
	if !(CommandStep{keys: []CommandStepKey{{key: CK_U}}}).IsDirToButton(CommandStep{keys: []CommandStepKey{{key: CK_a}}}) {
		t.Fatal("expected direction to button transition to be true")
	}
	if !(CommandStep{keys: []CommandStepKey{{key: CK_U, tilde: true}}}).IsDirToButton(CommandStep{keys: []CommandStepKey{{key: CK_a}}}) {
		t.Fatal("expected direction release to button transition to be true")
	}
	if (CommandStep{keys: []CommandStepKey{{key: CK_a}}}).IsDirToButton(CommandStep{keys: []CommandStepKey{{key: CK_b}}}) {
		t.Fatal("expected button-to-button transition to be false")
	}
}
