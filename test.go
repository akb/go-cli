package cli

import (
	"bytes"
	"os/exec"
	"regexp"
	"testing"
)

func ExpectSuccess(t *testing.T, err error) {
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("Command failed to execute")
		}

		t.Errorf("Expected zero exit status; received %s", err)
	}
}

func ExpectError(t *testing.T, err error) {
	if err == nil {
		t.Error("Expected nonzero exit status")
	}
}

func ExpectOutput(t *testing.T, stdout bytes.Buffer) {
	if len(stdout.Bytes()) < 1 {
		t.Errorf("Expected command to write output to STDOUT")
	}
}

func ExpectMatch(t *testing.T, stdout bytes.Buffer, pattern string) {
	matched, err := regexp.Match(pattern, stdout.Bytes())
	if err != nil {
		t.Fatalf("Unable to parse, bad regular expression: %s", pattern)
	}

	if !matched {
		t.Errorf("Output does not match pattern\nPattern: %s\nOutput:\n%s",
			pattern, stdout.String())
	}
}

func ExpectErrorOutput(t *testing.T, stderr bytes.Buffer, err error) {
	if len(stderr.Bytes()) < 1 {
		t.Errorf("Expected log output on STDERR with nonzero exit status")
	}
}

func ExpectHelp(t *testing.T, stderr bytes.Buffer, cmd Command) {
	var subcommands CLI
	if b, ok := (interface{})(cmd).(HasSubcommands); ok {
		subcommands = b.Subcommands()
	}

	for subcommand, _ := range subcommands {
		matched, err := regexp.Match(subcommand, stderr.Bytes())
		if err != nil {
			t.Fatalf("Unable to parse, bad regular expression: %s", subcommand)
		}

		if !matched {
			t.Errorf("Help text doesn't include subcommand: %s", subcommand)
		}
	}
}
