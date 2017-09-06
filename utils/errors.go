package utils

import (
	"bytes"
	"fmt"
	"text/scanner"
)

type RwErrType int

const (
	ErrGeneric RwErrType = iota
	ErrSyntaxError
	ErrEmpty
	ErrUnknown
)

// A EmError is a generic error returned for parsing errors.
// The first line is 1.  The first column is 0.
type EmError struct {
	Type RwErrType
	scanner.Position
	Message string
}

func (e EmError) Error() string {
	var msg string
	switch e.Type {
	case ErrSyntaxError:
		msg = fmt.Sprintf("%s:line %d:%d: %s", "syntax error", e.Line+e.Offset, e.Column, e.Message)
	case ErrEmpty:
		msg = fmt.Sprintf("%sline %d:%d: %s", "", e.Line+e.Offset, e.Column, "nothing to parse")
	default:
		msg = fmt.Sprintf("%sline %d:%d: %s", "", e.Line+e.Offset, e.Column, e.Message)
	}
	return msg
}

func NewError(etype RwErrType, pos scanner.Position, msg string) *EmError {
	//	fmt.Println(pos.Offset, " ", pos.Line, " ", msg)
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
	em.Errors = make([]error, 0, 10)
}

func (em *ErrorManager) String() string {
	var b bytes.Buffer
	for i := 0; i < len(em.Errors); i++ {
		fmt.Fprintln(&b, em.Errors[i].Error())
	}
	return b.String()
}
