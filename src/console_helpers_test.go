package main

import (
	"io"
	"os"
	"testing"
)

func TestConsoleHelpersAppendAndPrint(t *testing.T) {
	oldSys := sys
	oldStdout := os.Stdout
	defer func() {
		sys = oldSys
		os.Stdout = oldStdout
	}()

	sys = oldSys
	sys.cfg.Debug.ConsoleRows = 2
	sys.consoleText = nil

	sys.appendToConsole("first")
	sys.appendToConsole("second")
	sys.appendToConsole("third")
	if got := len(sys.consoleText); got != 2 || sys.consoleText[0] != "second" || sys.consoleText[1] != "third" {
		t.Fatalf("appendToConsole() truncation = %#v, want [second third]", sys.consoleText)
	}

	sys.stringPool[0].List = []string{"hello %s\nworld"}
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe(): %v", err)
	}
	os.Stdout = w
	sys.printToConsole(0, 0, "there")
	w.Close()
	out, err := io.ReadAll(r)
	r.Close()
	if err != nil {
		t.Fatalf("ReadAll(stdout): %v", err)
	}
	if string(out) != "hello there\nworld\n" {
		t.Fatalf("printToConsole() stdout = %q, want %q", string(out), "hello there\nworld\n")
	}
	if len(sys.consoleText) != 2 || sys.consoleText[0] != "world" || sys.consoleText[1] != "third" {
		t.Fatalf("printToConsole() consoleText = %#v, want tail append", sys.consoleText)
	}
}
