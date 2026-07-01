package main

import "testing"

func TestNewStageProps(t *testing.T) {
	sp := newStageProps()
	if sp.roundpos {
		t.Fatal("expected newStageProps to leave roundpos false")
	}
}
