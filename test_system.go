package cli

import (
	"bytes"
	"io"
	"log"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"golang.org/x/crypto/ssh/terminal"
)

type TestSystem struct {
	*BaseSystem
	Console *expect.Console
}

func NewTestSystem(t *testing.T, arguments []string, environment map[string]string) *TestSystem {
	var stdout bytes.Buffer

	timeout := 1 * time.Second
	console, err := expect.NewTestConsole(t, func(opts *expect.ConsoleOpts) error {
		opts.ReadTimeout = &timeout
		opts.Stdouts = []io.Writer{&stdout}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if environment == nil {
		environment = map[string]string{}
	}

	return &TestSystem{&BaseSystem{
		In:          console.Tty(),
		Out:         console.Tty(),
		Logger:      log.New(console.Tty(), "", log.LstdFlags),
		Environment: environment,
		Arguments:   arguments,
	}, console}
}

func (ts *TestSystem) ScanSilent() (string, error) {
	cloaked, err := terminal.ReadPassword(int(ts.Console.Tty().Fd()))
	if err != nil {
		return "", nil
	}
	return string(cloaked), nil
}

func (ts *TestSystem) Close() {
	ts.Console.Tty().Close()
	ts.Console.Close()
}
