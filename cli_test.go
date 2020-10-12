package cli

import (
	"context"
	"flag"
	"testing"
)

type testCommand struct {
	helpDidRun        bool
	flagsDidRun       bool
	commandDidRun     bool
	subcommandsDidRun bool
}

type testMainCommand struct {
	*testCommand

	subc Command
}

func (c *testMainCommand) Help() {
	c.helpDidRun = true
}

func (c *testMainCommand) Flags(f *flag.FlagSet) {
	c.flagsDidRun = true
}

func (c *testMainCommand) Command(ctx context.Context, args []string, s System) {
	c.commandDidRun = true
}

func (c *testMainCommand) Subcommands() CLI {
	c.subcommandsDidRun = true
	return CLI{"testsub": c.subc}
}

type testSubcommand struct {
	*testCommand
}

func (c *testSubcommand) Help() {
	c.helpDidRun = true
}

func (c *testSubcommand) Flags(f *flag.FlagSet) {
	c.flagsDidRun = true
}

func (c *testSubcommand) Command(ctx context.Context, args []string, s System) {
	c.commandDidRun = true
}

func (c *testSubcommand) Subcommands() CLI {
	c.subcommandsDidRun = true
	return nil
}

func TestMainCommand(t *testing.T) {
	subc := &testSubcommand{&testCommand{}}
	cmd := &testMainCommand{&testCommand{}, subc}
	system, _ := NewTestSystem(t, []string{"testmain"}, nil)
	result := Main(context.Background(), cmd, system)

	if result != 0 {
		t.Errorf("command did not return a 0 status\n")
	}

	if cmd.helpDidRun {
		t.Errorf("cmd.Help ran but should not have\n")
	}

	if !cmd.flagsDidRun {
		t.Errorf("cmd.Flags method did not run\n")
	}

	if !cmd.commandDidRun {
		t.Errorf("cmd.Command method did not run\n")
	}

	if !cmd.subcommandsDidRun {
		t.Errorf("cmd.Subcommands method did not run\n")
	}

	if subc.helpDidRun {
		t.Errorf("subc.Help method ran but should not have\n")
	}

	if subc.flagsDidRun {
		t.Errorf("subc.Flags method ran but should not have\n")
	}

	if subc.commandDidRun {
		t.Errorf("subc.Command method ran but should not have\n")
	}

	if subc.subcommandsDidRun {
		t.Errorf("subc.Subcommands method ran but should not have\n")
	}
}

func TestSubcommand(t *testing.T) {
	subc := &testSubcommand{&testCommand{}}
	cmd := &testMainCommand{&testCommand{}, subc}
	system, _ := NewTestSystem(t, []string{"testmain", "testsub"}, nil)
	result := Main(context.Background(), cmd, system)

	if result != 0 {
		t.Errorf("command did not return a 0 status\n")
	}

	if cmd.helpDidRun {
		t.Errorf("cmd.Help ran but should not have\n")
	}

	if cmd.flagsDidRun {
		t.Errorf("cmd.Flags ran but should not have\n")
	}

	if cmd.commandDidRun {
		t.Errorf("cmd.Command ran but should not have\n")
	}

	if !cmd.subcommandsDidRun {
		t.Errorf("cmd.Subcommands method did not run\n")
	}

	if subc.helpDidRun {
		t.Errorf("subc.Help method ran but should not have\n")
	}

	if !subc.flagsDidRun {
		t.Errorf("subc.Flags did not run\n")
	}

	if !subc.commandDidRun {
		t.Errorf("subc.Command did not run\n")
	}

	if subc.subcommandsDidRun {
		t.Errorf("subc.Subcommands method ran but should not have\n")
	}
}

type testFatalCmd struct {
	msg string
}

func (testFatalCmd) Help() {}

func (testFatalCmd) Command(c context.Context, args []string, s System) {
	s.Fatalln("FUBRRRR")
}

func TestFatal(t *testing.T) {
	cmd := &testFatalCmd{}

	system, _ := NewTestSystem(t, []string{"testmain"}, nil)
	status := Main(context.Background(), cmd, system)

	if status != 1 {
		t.Errorf("Expected command to return status 1, instead received %d.\n", status)
	}

	received, err := system.Console.ExpectString("FUBRRRR")
	if err != nil {
		t.Errorf("Error while expecting string: %s\n", err.Error())
		t.Logf("Received: \"%s\"\n", received)
	}
}
