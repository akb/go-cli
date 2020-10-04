package cli

import (
	"fmt"
	"io"
	"log"
)

// System is passed to commands as an argument when the command is run. It
// provides an IO interface for the command to use that can be easily attached
// to STDIN/STDOUT or to bytes.Buffer for testing
type System struct {
	In          io.Reader
	Out         io.Writer
	Logger      *log.Logger
	Environment map[string]string
	Arguments   []string
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
