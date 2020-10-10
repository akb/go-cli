package cli

import (
	"bytes"
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

type TestOutput struct {
	*bytes.Buffer
}

func NewTestSystem(t *testing.T, arguments []string, environment map[string]string) *TestSystem {
	stdout := TestOutput{&bytes.Buffer{}}
	stderr := TestOutput{&bytes.Buffer{}}

	console, err := expect.NewConsole(
		expect.WithDefaultTimeout(1*time.Second),
		expect.WithStdout(stdout),
	)
	if err != nil {
		t.Fatal(err)
	}

	if environment == nil {
		environment = map[string]string{}
	}

	tty := console.Tty()

	return &TestSystem{
		BaseSystem: &BaseSystem{
			In:          tty,
			Out:         tty,
			Logger:      log.New(stderr, "", log.LstdFlags),
			Environment: environment,
			Arguments:   arguments,
		},
		Console: console,
	}
}

func (ts *TestSystem) ReadPassword() (string, error) {
	cloaked, err := terminal.ReadPassword(int(ts.Console.Tty().Fd()))
	if err != nil {
		return "", err
	}
	return string(cloaked), nil
}
