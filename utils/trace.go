package utils

import (
	"fmt"
	"io"
	"os"
)

//Tracer defines interface for outputing debug and tracing info
type Tracer interface {
	Printf(string, ...interface{}) (int, error)
	Print(...interface{}) (int, error)
	Println(...interface{}) (int, error)
	On() error
	Off() error
}

type xtrace struct {
	on     bool
	writer io.Writer
}

func NewTrace(on bool, writer io.Writer) Tracer {
	tr := &xtrace{on, writer}
	if tr.writer == nil {
		tr.writer = os.Stdout
	}
	return tr
}

func (tr *xtrace) Printf(format string, a ...interface{}) (n int, err error) {
	if tr.on {
		return fmt.Fprintf(tr.writer, format, a...)
	}
	return 0, nil
}

func (tr *xtrace) Print(a ...interface{}) (n int, err error) {
	if tr.on {
		return fmt.Fprint(tr.writer, a...)
	}
	return 0, nil
}

func (tr *xtrace) Println(a ...interface{}) (n int, err error) {
	if tr.on {
		return fmt.Fprintln(tr.writer, a...)
	}
	return 0, nil
}

func (tr *xtrace) On() error {
	tr.on = true
	return nil
}

func (tr *xtrace) Off() error {
	tr.on = false
	return nil
}
