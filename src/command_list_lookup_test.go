package main

import "testing"

func TestCommandListAddAtAndGetHandleExistingAndMissingNames(t *testing.T) {
	cl := NewCommandList(NewInputBuffer())

	cl.Add(Command{name: "jump"})
	if len(cl.Commands) != 1 || cl.Names["jump"] != 0 || len(cl.At(0)) != 1 {
		t.Fatalf("first Add() state = names=%#v commands=%#v", cl.Names, cl.Commands)
	}

	cl.Add(Command{name: "jump"})
	if len(cl.Commands) != 1 || len(cl.At(0)) != 2 {
		t.Fatalf("second Add() should append to existing bucket, got commands=%#v", cl.Commands)
	}

	cl.Add(Command{name: "kick"})
	if len(cl.Commands) != 2 || cl.Names["kick"] != 1 {
		t.Fatalf("new command name should create a new bucket, got names=%#v commands=%#v", cl.Names, cl.Commands)
	}

	if got := cl.Get("jump"); len(got) != 2 {
		t.Fatalf("Get(jump) = %d commands, want 2", len(got))
	}
	if got := cl.Get("missing"); got != nil {
		t.Fatalf("Get(missing) = %#v, want nil", got)
	}
	if got := cl.At(-1); got != nil {
		t.Fatalf("At(-1) = %#v, want nil", got)
	}
	if got := cl.At(99); got != nil {
		t.Fatalf("At(out of range) = %#v, want nil", got)
	}
}
