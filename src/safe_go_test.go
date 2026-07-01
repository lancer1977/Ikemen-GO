package main

import (
	"testing"
	"time"
)

func TestSafeGo_RunsFunctionAndForwardsPanicsToMainThreadTask(t *testing.T) {
	oldTask := sys.mainThreadTask
	defer func() { sys.mainThreadTask = oldTask }()
	sys.mainThreadTask = make(chan func(), 1)

	done := make(chan struct{}, 1)
	SafeGo(func() {
		done <- struct{}{}
	})

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("SafeGo did not run the function")
	}

	SafeGo(func() {
		panic("boom")
	})

	select {
	case task := <-sys.mainThreadTask:
		if task == nil {
			t.Fatal("expected panic task to be forwarded")
		}
	case <-time.After(time.Second):
		t.Fatal("SafeGo did not forward panic to mainThreadTask")
	}
}
