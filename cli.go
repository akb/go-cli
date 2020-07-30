package cli

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Command is an interface used to represent a CLI component. Both primary
// commands and subcommands implement Command
type Command interface {
	// Help is called for a command if the command line fails to parse. It may
	// also be manually called in the `Command` method if appropriate.
	Help()

	// Flags is called before `Command` and is passed a pointer to a flag.FlagSet
	// where the Command may define flags to be automatically parsed
	Flags(*flag.FlagSet)

	// Command is the method that actually performs the command.
	Command(context.Context, []string) int

	// Subcommands should return nil, or a pointer to a CLI if the command has
	// subcommands
	Subcommands() CLI
}

// CLI is a map of names to Command implementations. It is used to represent a
// set of subcommands for a given Command
type CLI map[string]Command

// Main should be called from a CLI application's `main` function. It should be
// passed the Command that represents the root of the subcommand tree. Main
// will parse the command line, determine which subcommand is the intended
// target, create a FlagSet then execute that subcommand. If no suitable
// subcommand is found, or if flag parsing fails, it will call the Help method
// from the most-recently visited subcommand. Main returns the Unix status code
// which should be returned to the underlying OS
func Main(mainCmd Command) int {
	var cmd Command
	var name string
	var flags []string
	var positional []string
	var head string
	var args []string = os.Args
	for {
		head = args[0]
		args = args[1:]

		if head[0] == '-' {
			flags = append(flags, head)
			continue
		}

		subcommands := cmd.Subcommands()
		if subcommands == nil {
			positional = append(positional, head)
			continue
		}

		c, ok := subcommands[head]
		if !ok {
			positional = append(positional, head)
			continue
		}

		if len(name) == 0 {
			name = head
		} else {
			name = strings.Join([]string{name, head}, " ")
		}

		cmd = c

		if len(args) == 0 {
			break
		}
	}

	f := flag.NewFlagSet(name, flag.ExitOnError)
	f.Usage = cmd.Help
	cmd.Flags(f)
	err := f.Parse(flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse command-line arguments:\n%s\n", err)
		return 1
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "origin", name)
	stamp := time.Now().UnixNano()
	traceID := fmt.Sprintf("%x", sha256.Sum256([]byte(string(stamp))))[:45]
	ctx = context.WithValue(ctx, "trace-id", traceID)

	if status := cmd.Command(ctx, positional); status != 0 {
		return status
	}

	return 0
}
