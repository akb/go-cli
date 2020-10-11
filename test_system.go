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
	STDOUT *bytes.Buffer
	STDERR *bytes.Buffer
}

func NewTestSystem(
	t *testing.T, arguments []string, environment map[string]string,
) (*TestSystem, *TestOutput) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	console, err := expect.NewTestConsole(t,
		expect.WithDefaultTimeout(5*time.Second),
		expect.WithStdout(stdout),
	)
	if err != nil {
		t.Fatal(err)
	}

	if environment == nil {
		environment = map[string]string{}
	}

	return &TestSystem{
		BaseSystem: &BaseSystem{
			In:          console.Tty(),
			Out:         console.Tty(),
			Logger:      log.New(stderr, "", log.LstdFlags),
			Environment: environment,
			Arguments:   arguments,
		},
		Console: console,
	}, &TestOutput{stdout, stderr}
}

func (ts *TestSystem) ReadPassword() (string, error) {
	cloaked, err := terminal.ReadPassword(int(ts.Console.Tty().Fd()))
	if err != nil {
		return "", err
	}
	return string(cloaked), nil
}
