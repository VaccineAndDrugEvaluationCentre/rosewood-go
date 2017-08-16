package rosewood

import (
	"bytes"
	"fmt"
	"text/scanner"
)

type rwError int

const (
	ErrGeneric rwError = iota
	ErrSyntaxError
	ErrEmpty
)

// A EmError is a generic error returned for parsing errors.
// The first line is 1.  The first column is 0.
type EmError struct {
	Type rwError
	scanner.Position
	Message string
}

func (e EmError) Error() string {
	var msg string
	switch e.Type {
	case ErrSyntaxError:
		msg = fmt.Sprintf("%sline %d:%d: %s", "syntax error:", e.Line, e.Column, e.Message)
	case ErrEmpty:
		msg = fmt.Sprintf("%sline %d:%d: %s", "", e.Line, e.Column, "nothing to parse")
	default:
		msg = fmt.Sprintf("%sline %d:%d: %s", "", e.Line, e.Column, e.Message)
	}
	return msg
}

func NewError(etype rwError, pos scanner.Position, msg string) *EmError {
	return &EmError{etype, pos, msg}
}

//ErrorManager stores and prints errors
type ErrorManager struct {
	Errors []error
}

func NewErrorManager() *ErrorManager {
	return &ErrorManager{Errors: make([]error, 0, 10)} //start with an initial capacity of 10 errors
}

func (em *ErrorManager) Add(e error) {
	em.Errors = append(em.Errors, e)
}

func (em *ErrorManager) Count() int {
	return len(em.Errors)
}

func (em *ErrorManager) Reset() {
	em.Errors = nil //for clarity
	em.Errors = make([]error, 0, 10)
}

func (em *ErrorManager) String() string {
	var b bytes.Buffer
	for i := 0; i < len(em.Errors); i++ {
		b.WriteString(em.Errors[i].Error())
		b.WriteString("\n")
	}
	return b.String()
}
