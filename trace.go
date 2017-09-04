package rosewood

import (
	"fmt"
	"io"
	"os"
)

//TODO: move to own package and make race-proof

type tracer interface {
	Printf(string, ...interface{}) (int, error)
	Print(...interface{}) (int, error)
	Println(...interface{}) (int, error)
	On() error
	Off() error
}

type traceState int

const (
	off traceState = iota
	on
)

type xtrace struct {
	state  traceState
	writer io.Writer
}

func newTrace(state traceState, writer io.Writer) tracer {
	tr := &xtrace{state, writer}
	if tr.writer == nil {
		tr.writer = os.Stdout
	}
	return tr
}

func (tr *xtrace) Printf(format string, a ...interface{}) (n int, err error) {
	if tr.state == on {
		return fmt.Fprintf(tr.writer, format, a...)
	}
	return 0, nil
}

func (tr *xtrace) Print(a ...interface{}) (n int, err error) {
	if tr.state == on {
		return fmt.Fprint(tr.writer, a...)
	}
	return 0, nil
}

func (tr *xtrace) Println(a ...interface{}) (n int, err error) {
	if tr.state == on {
		return fmt.Fprintln(tr.writer, a...)
	}
	return 0, nil
}

func (tr *xtrace) On() error {
	tr.state = on
	return nil
}

func (tr *xtrace) Off() error {
	tr.state = off
	return nil
}
