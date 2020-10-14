package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// System is passed to commands as an argument when the command is run. It
// provides an IO interface for the command to use that can be easily attached
// to STDIN/STDOUT or to bytes.Buffer for testing
type System interface {
	Environ() []string
	Getenv(string) string
	Args() []string

	Print(...interface{}) (int, error)
	Printf(string, ...interface{}) (int, error)
	Println(...interface{}) (int, error)

	Scan(...interface{}) (int, error)
	Scanf(string, ...interface{}) (int, error)

	Log(...interface{})
	Logf(string, ...interface{})

	ReadPassword() (string, error)
}

type BaseSystem struct {
	In          io.Reader
	Out         io.Writer
	Logger      *log.Logger
	Environment map[string]string
	Arguments   []string
}

func (s *BaseSystem) Environ() []string {
	environ := make([]string, len(s.Environment))

	i := 0
	for k, v := range s.Environment {
		environ[i] = fmt.Sprintf("%s=%s", k, v)
		i += 1
	}

	sort.Strings(environ)

	return environ
}

func (s *BaseSystem) Getenv(k string) string {
	if v, ok := s.Environment[k]; ok {
		return v
	}
	return ""
}

func (s *BaseSystem) Args() []string {
	return s.Arguments
}

func (s *BaseSystem) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(s.Out, a...)
}

func (s *BaseSystem) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(s.Out, format, a...)
}

func (s *BaseSystem) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(s.Out, a...)
}

func (s *BaseSystem) Scan(a ...interface{}) (int, error) {
	return fmt.Fscan(s.In, a...)
}

func (s *BaseSystem) Scanf(format string, a ...interface{}) (int, error) {
	return fmt.Fscanf(s.In, format, a...)
}

func (s *BaseSystem) Scanln(a ...interface{}) (int, error) {
	return fmt.Fscanln(s.In, a...)
}

func (s *BaseSystem) Log(a ...interface{}) {
	s.Logger.Println(a...)
}

func (s *BaseSystem) Logf(format string, a ...interface{}) {
	s.Logger.Printf(format, a...)
}

type UnixSystem struct {
	*BaseSystem
}

func NewUnixSystem() *UnixSystem {
	env := os.Environ()
	environment := make(map[string]string, len(env))
	for _, e := range env {
		split := strings.Split(e, "=")
		environment[split[0]] = split[1]
	}

	arguments := os.Args

	return &UnixSystem{&BaseSystem{
		In:          os.Stdin,
		Out:         os.Stdout,
		Logger:      log.New(os.Stderr, "", log.LstdFlags),
		Environment: environment,
		Arguments:   arguments,
	}}
}

func (s *UnixSystem) ReadPassword() (string, error) {
	cloaked, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", nil
	}
	return string(cloaked), nil
}
