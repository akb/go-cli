package cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
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
	var out string
	_, err := fmt.Fscan(ts.In, &out)
	if err != nil {
		return "", err
	}
	return out, nil
}

func (ts *TestSystem) Close() {
	ts.Console.Tty().Close()
	ts.Console.Close()
}
