package parser

import (
	"fmt"
	"text/scanner"
)

const (
	ErrGeneric int = iota
	ErrSyntaxError
	ErrEmpty
	ErrUnknown
)

// A EmError is a generic error returned for parsing errors.
type EmError struct {
	Type int
	scanner.Position
	Message string
}

// EmError implements the error interface
func (e EmError) Error() string {
	var msg string
	formatPos := func(pos scanner.Position) string {
		if pos.Offset > 0 && pos.Line > 0 {
			return fmt.Sprintf("line:%d col %d:", e.Line+e.Offset, e.Column)
		}
		return ""
	}
	switch e.Type {
	case ErrSyntaxError:
		msg = fmt.Sprintf("%s %s: %s", "syntax error", formatPos(e.Position), e.Message)
	case ErrEmpty:
		msg = fmt.Sprintf("%s: %s", formatPos(e.Position), "nothing to parse")
	default:
		msg = fmt.Sprintf("%s: %s", formatPos(e.Position), e.Message)
	}
	return msg
}

//NewError returns a pointer to a new EmError
func NewError(etype int, pos scanner.Position, msg string) *EmError {
	return &EmError{etype, pos, msg}
}
