package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestLogMessage_WritesFormattedLineToStderr(t *testing.T) {
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stderr = w
	defer func() {
		os.Stderr = oldStderr
	}()

	LogMessage("hello %s %d", "world", 7)
	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read stderr: %v", err)
	}
	if !bytes.Equal(got, []byte("hello world 7\n")) {
		t.Fatalf("LogMessage stderr = %q, want %q", string(got), "hello world 7\n")
	}
}
