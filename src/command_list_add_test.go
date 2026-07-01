package main

import "testing"

func TestCommandListAddCommandHandlesNilAndBlankInputs(t *testing.T) {
	var cl *CommandList
	if err := cl.AddCommand("jump", CommandSpec{Cmd: "a"}); err == nil {
		t.Fatal("AddCommand() on nil CommandList should return an error")
	}

	c := NewCommandList(NewInputBuffer())
	if err := c.AddCommand("jump", CommandSpec{Cmd: "   "}); err != nil {
		t.Fatalf("AddCommand() blank command returned error: %v", err)
	}
	if len(c.Commands) != 0 || len(c.Names) != 0 {
		t.Fatalf("AddCommand() blank command should be a no-op, got commands=%d names=%d", len(c.Commands), len(c.Names))
	}
}

func TestCommandListAddCommandStoresCompiledCommand(t *testing.T) {
	c := NewCommandList(NewInputBuffer())
	spec := CommandSpec{
		Cmd:            "a",
		Time:           9,
		BufTime:        4,
		BufferHitpause: false,
		BufferPauseend: false,
		StepTime:       2,
	}
	if err := c.AddCommand("punch", spec); err != nil {
		t.Fatalf("AddCommand() = %v", err)
	}
	if len(c.Commands) != 1 {
		t.Fatalf("AddCommand() command count = %d, want 1", len(c.Commands))
	}
	if got := c.Names["punch"]; got != 0 {
		t.Fatalf("AddCommand() name index = %d, want 0", got)
	}
	cmd := c.Commands[0][0]
	if cmd.name != "punch" || cmd.maxtime != 9 || cmd.maxbuftime != 4 || cmd.buffer_hitpause || cmd.buffer_pauseend || cmd.maxsteptime != 2 {
		t.Fatalf("AddCommand() stored command = %#v", cmd)
	}
}
