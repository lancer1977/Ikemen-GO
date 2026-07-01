package main

import (
	"bytes"
	"os"
	"testing"
)

func TestPrintBytecodeErrorRoutesToConsoleDuringMatches(t *testing.T) {
	oldSys := sys
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		sys = oldSys
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	sys = oldSys
	sys.loader.state = LS_Complete
	sys.tickCount = 123
	sys.workingChar = &Char{name: "Ryu", id: 7}
	sys.workingChar.ss.no = 42
	sys.cfg.Debug.ConsoleRows = 8
	sys.consoleText = nil

	sys.printBytecodeError("boom")
	if len(sys.consoleText) != 1 {
		t.Fatalf("expected one console line, got %#v", sys.consoleText)
	}
	if got := sys.consoleText[0]; got != "123: WARNING: Ryu (7) in state 42: boom" {
		t.Fatalf("unexpected console message: %q", got)
	}
}

func TestPrintBytecodeErrorWritesToStderrOutsideMatches(t *testing.T) {
	oldSys := sys
	oldStderr := os.Stderr
	defer func() {
		sys = oldSys
		os.Stderr = oldStderr
	}()

	sys = oldSys
	sys.loader.state = 0
	sys.workingChar = nil
	sys.ignoreMostErrors = false

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe(): %v", err)
	}
	os.Stderr = w

	sys.printBytecodeError("compile boom")
	w.Close()
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("reading stderr: %v", err)
	}
	r.Close()

	if got := buf.String(); got != "compile boom\n" {
		t.Fatalf("stderr output = %q, want %q", got, "compile boom\n")
	}
}
