package cli

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ExitError can be used to expecify an exit status when returning from Command
type ExitError struct {
	Status  int
	Message string
}

// Error returns a string representation of the ExitError
func (e *ExitError) Error() string {
	return e.Message
}

// Command is an interface used to represent a CLI component. Both primary
// commands and subcommands implement Command
type Command interface {
	// Help is called for a command if the command line fails to parse. It may
	// also be manually called in the `Command` method if appropriate.
	Help()
}

// Action is an interface for commands that do things other than display
// information
type Action interface {
	// Command is the method that actually performs the command.
	Command(context.Context, []string, System) error
}

// HasFlags is an interface for commands that use flags
type HasFlags interface {
	// Flags is called before `Command` and is passed a pointer to a flag.FlagSet
	// where the Command may define flags to be automatically parsed
	Flags(*flag.FlagSet)
}

// HasSubcommands is an interface for commands that have subcommands
type HasSubcommands interface {
	// Subcommands should return a CLI if the command has subcommands
	Subcommands() CLI
}

// NoOpCommand is a command that does nothing.
type NoOpCommand struct{}

func (NoOpCommand) Help() {}

// Command performs no operation and returns successfully
func (NoOpCommand) Command(c context.Context, args []string, s *System) error {
	return nil
}

// CLI is a map of names to Command implementations. It is used to represent a
// set of subcommands for a given Command
type CLI map[string]Command

// ListSubcommands returns a slice of names of the subcommands within a CLI
func (c CLI) ListSubcommands(prefix string) []string {
	var subcommands []string
	for name, cmd := range c {
		if len(prefix) > 0 {
			name = fmt.Sprintf("%s %s", prefix, name)
		}

		if _, ok := (interface{})(cmd).(Action); ok {
			subcommands = append(subcommands, name)
		}

		if sc, ok := cmd.(HasSubcommands); ok {
			for _, sck := range sc.Subcommands().ListSubcommands(name) {
				if _, ok := (interface{})(cmd).(Action); ok {
					subcommands = append(subcommands, sck)
				}
			}
		}
	}
	return subcommands
}

// Main should be called from a CLI application's `main` function. It should be
// passed the Command that represents the root of the subcommand tree. Main
// will parse the command line, determine which subcommand is the intended
// target, create a FlagSet then execute that subcommand. If no suitable
// subcommand is found, or if flag parsing fails, it will call the Help method
// from the most-recently visited subcommand. Main returns the Unix status code
// which should be returned to the underlying OS
func Main(ctx context.Context, mainCmd Command, sys System) (status int) {
	var cmd Command = mainCmd
	var args, flags []string
	var head, name string
	var tail []string = sys.Args()
	for {
		var subcommands CLI
		if b, ok := (interface{})(cmd).(HasSubcommands); ok {
			subcommands = b.Subcommands()
		}

		head = tail[0]
		if head[0] == '-' {
			flags = append(flags, head)
		} else if subcommands == nil {
			args = append(args, head)
		} else {
			if subcommand, ok := subcommands[head]; ok {
				cmd = subcommand

				if len(name) == 0 {
					name = head
				} else {
					name = strings.Join([]string{name, head}, " ")
				}
			} else if head != sys.Args()[0] {
				args = append(args, head)
			}
		}

		if len(tail) == 1 {
			break
		}

		tail = tail[1:]
	}

	if b, ok := (interface{})(cmd).(HasFlags); ok {
		f := flag.NewFlagSet(name, flag.ExitOnError)
		f.Usage = cmd.Help
		b.Flags(f)
		if err := f.Parse(flags); err != nil {
			sys.Logf("Failed to parse command-line arguments:\n%s\n", err)
			return 1
		}
	}

	if b, ok := (interface{})(cmd).(Action); ok {
		ctx = context.WithValue(ctx, "origin", name)
		ctx = context.WithValue(ctx, "trace-id", traceID())
		if err := b.Command(ctx, args, sys); err != nil {
			sys.Log(err.Error())
			switch err := errors.Cause(err).(type) {
			case *ExitError:
				return err.Status
			default:
				return 1
			}
		}
	}

	return 0
}

func traceID() string {
	stamp := []byte(string(time.Now().UnixNano()))
	return fmt.Sprintf("%x", sha256.Sum256(stamp))[:45]
}
