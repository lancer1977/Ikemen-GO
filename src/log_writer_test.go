package main

import (
	"os"
	"testing"
)

func TestNewLogWriter(t *testing.T) {
	t.Parallel()

	if got := NewLogWriter(); got != os.Stderr {
		t.Fatalf("NewLogWriter() = %#v, want os.Stderr", got)
	}
}
