package main

import "testing"

func TestLoadingCanceled(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys.loader = newLoader()
	sys.gameEnd = false
	if loadingCanceled() {
		t.Fatal("expected loadingCanceled to be false with default state")
	}

	sys.gameEnd = true
	if !loadingCanceled() {
		t.Fatal("expected loadingCanceled to be true when gameEnd is set")
	}

	sys.gameEnd = false
	sys.loader.state = LS_Cancel
	if !loadingCanceled() {
		t.Fatal("expected loadingCanceled to be true when loader is canceled")
	}

	sys.loader.state = LS_NotYet
	sys.loader.cancelCh = make(chan struct{})
	sys.loader.requestCancel()
	if !loadingCanceled() {
		t.Fatal("expected loadingCanceled to be true when cancel is requested")
	}
}
