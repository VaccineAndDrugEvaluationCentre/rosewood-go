// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package parser

import (
	"fmt"
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
	Position
	Message string
}

// EmError implements the error interface
func (e EmError) Error() string {
	formatPos := func(pos Position) string {
		//fmt.Printf("%d %d\n", pos.Offset, pos.Line)
		if pos.Offset > 0 && pos.Line > -1 {
			return fmt.Sprintf("line #%d col #%d", e.Line+e.Offset, e.Column)
		}
		return ""
	}
	switch e.Type {
	case ErrSyntaxError:
		return fmt.Sprintf("%s %s: %s", "syntax error", formatPos(e.Position), e.Message)
	case ErrEmpty:
		return fmt.Sprintf("%s: %s", formatPos(e.Position), "nothing to parse")
	default:
		return fmt.Sprintf("%s: %s", formatPos(e.Position), e.Message)
	}
}

//NewError returns a pointer to a new EmError
func NewError(etype int, pos Position, msg string) *EmError {
	return &EmError{etype, pos, msg}
}
