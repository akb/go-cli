package cli

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
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
}

// Action is an interface for commands that do things other than display
// information
type Action interface {
	// Command is the method that actually performs the command.
	Command(context.Context, []string, System) int
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

func (NoOpCommand) Command(c context.Context, args []string, s System) int {
	return 0
}

// CLI is a map of names to Command implementations. It is used to represent a
// set of subcommands for a given Command
type CLI map[string]Command

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

// System is passed to commands as an argument when the command is run. It
// provides an IO interface for the command to use that can be easily attached
// to STDIN/STDOUT or to bytes.Buffer for testing
type System struct {
	In     io.Reader
	Out    io.Writer
	Logger *log.Logger
	Env    map[string]string
}

func (s System) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(s.Out, a...)
}

func (s System) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(s.Out, format, a...)
}

func (s System) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(s.Out, a...)
}

func (s System) Scan(a ...interface{}) (int, error) {
	return fmt.Fscan(s.In, a...)
}

func (s System) Scanf(format string, a ...interface{}) (int, error) {
	return fmt.Fscanf(s.In, format, a...)
}

func (s System) Scanln(a ...interface{}) (int, error) {
	return fmt.Fscanln(s.In, a...)
}

func (s System) Log(a ...interface{}) {
	s.Logger.Println(a...)
}

func (s System) Logf(format string, a ...interface{}) {
	s.Logger.Printf(format, a...)
}

func (s System) Fatal(v ...interface{}) {
	s.Logger.Fatal(v...)
}

func (s System) Fatalf(format string, v ...interface{}) {
	s.Logger.Fatalf(format, v...)
}

func (s System) Fatalln(v ...interface{}) {
	s.Logger.Fatalln(v...)
}

func (s System) Getenv(key string) string {
	if v, ok := s.Env[key]; ok {
		return v
	} else {
		return ""
	}
}

// Main should be called from a CLI application's `main` function. It should be
// passed the Command that represents the root of the subcommand tree. Main
// will parse the command line, determine which subcommand is the intended
// target, create a FlagSet then execute that subcommand. If no suitable
// subcommand is found, or if flag parsing fails, it will call the Help method
// from the most-recently visited subcommand. Main returns the Unix status code
// which should be returned to the underlying OS
func Main(mainCmd Command, in io.Reader, out io.Writer, logger *log.Logger) int {
	var cmd Command = mainCmd
	var args, flags []string
	var head, name string
	var tail []string = os.Args
	for {
		head = tail[0]

		var subcommands CLI
		if b, ok := (interface{})(cmd).(HasSubcommands); ok {
			subcommands = b.Subcommands()
		}

		if head[0] == '-' {
			flags = append(flags, head)
		} else if subcommands == nil {
			args = append(args, head)
		} else {
			if c, ok := subcommands[head]; ok {
				cmd = c

				if len(name) == 0 {
					name = head
				} else {
					name = strings.Join([]string{name, head}, " ")
				}
			} else if head != os.Args[0] {
				args = append(args, head)
			}
		}

		if len(tail) == 1 {
			break
		} else {
			tail = tail[1:]
			continue
		}
	}

	if b, ok := (interface{})(cmd).(HasFlags); ok {
		f := flag.NewFlagSet(name, flag.ExitOnError)
		f.Usage = cmd.Help
		b.Flags(f)
		if err := f.Parse(flags); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse command-line arguments:\n%s\n", err)
			return 1
		}
	}

	if b, ok := (interface{})(cmd).(Action); ok {
		env := map[string]string{}
		for _, e := range os.Environ() {
			split := strings.Split(e, "=")
			env[split[0]] = split[1]
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, "origin", name)
		stamp := time.Now().UnixNano()
		traceID := fmt.Sprintf("%x", sha256.Sum256([]byte(string(stamp))))[:45]
		ctx = context.WithValue(ctx, "trace-id", traceID)

		if status := b.Command(ctx, args, System{in, out, logger, env}); status != 0 {
			return status
		}
	}

	return 0
}
