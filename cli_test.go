package cli

import (
	"context"
	"flag"
	"log"
	"os"
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

func (c *testMainCommand) Command(ctx context.Context, args []string, s *System) int {
	c.commandDidRun = true
	return 0
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

func (c *testSubcommand) Command(ctx context.Context, args []string, s *System) int {
	c.commandDidRun = true
	return 0
}

func (c *testSubcommand) Subcommands() CLI {
	c.subcommandsDidRun = true
	return nil
}

func TestMainCommand(t *testing.T) {
	subc := &testSubcommand{&testCommand{}}
	cmd := &testMainCommand{&testCommand{}, subc}

	result := Main(cmd, &System{
		In:          os.Stdin,
		Out:         os.Stdout,
		Logger:      log.New(os.Stderr, "", log.LstdFlags),
		Arguments:   []string{"testmain"},
		Environment: map[string]string{},
	})

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

	result := Main(cmd, &System{
		In:          os.Stdin,
		Out:         os.Stdout,
		Logger:      log.New(os.Stderr, "", log.LstdFlags),
		Arguments:   []string{"testmain", "testsub"},
		Environment: map[string]string{},
	})

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
