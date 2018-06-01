// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package errors

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ErrorList is a list of errors; a more generic version of Go's parser ErrorList
// The zero value for an ErrorList is an empty ErrorList ready to use.
type ErrorList []error

//NewErrorList returns an empty ErrorList; not needed but a good practice in case we change the implementation
func NewErrorList() ErrorList {
	return ErrorList{}
}

// Add adds an error to an ErrorList.
func (p *ErrorList) Add(e error) {
	*p = append(*p, e)
}

// Count returns the number of errors in the List; preferred to len(errorList)
func (p *ErrorList) Count() int {
	return len(*p)
}

// Reset resets an ErrorList to no errors.
func (p *ErrorList) Reset() { *p = (*p)[0:0] }

// ErrorList implements the sort Interface; override Less for any special sorting needs
func (p ErrorList) Len() int      { return len(p) }
func (p ErrorList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p ErrorList) Less(i, j int) bool {
	return p[i].Error() < p[j].Error()
}

// Sort sorts an ErrorList. The default sorting is by error message (asc)
func (p ErrorList) Sort() {
	sort.Sort(p)
}

// RemoveMultiples sorts an ErrorList and removes all but the first error.
func (p *ErrorList) RemoveMultiples() {
	sort.Sort(p)
	var last string
	i := 0
	for _, e := range *p {
		if e.Error() != last {
			last = e.Error()
			(*p)[i] = e
			i++
		}
	}
	(*p) = (*p)[0:i]
}

// An ErrorList implements the error interface.
func (p ErrorList) Error() string {
	switch len(p) {
	case 0:
		return "no errors"
	case 1:
		return p[0].Error()
	}
	//TODO: change to default:
	return fmt.Sprintf("%s (and %d more errors)", p[0], len(p)-1)
}

// Err returns an error equivalent to this error list.
// If the list is empty, Err returns nil.
func (p ErrorList) Err() error {
	if len(p) == 0 {
		return nil
	}
	return p
}

// PrintError is a utility function that prints a list of errors to w,
// one error per line, if the err parameter is an ErrorList. Otherwise
// it prints the err string.
//
func PrintError(w io.Writer, err error) {
	if list, ok := err.(ErrorList); ok {
		for _, e := range list {
			fmt.Fprintf(w, "%s\n", e)
		}
	} else if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	}
}

//ErrorsToError converts list of error into one error
func ErrorsToError(err error) error {
	switch e := err.(type) {
	case ErrorList:
		var w strings.Builder
		PrintError(&w, e)
		return fmt.Errorf("%s", w.String())
	default:
		return e
	}
}

//ErrorsToString converts list of error into one error
func ErrorsToString(err error) string {
	switch e := err.(type) {
	case ErrorList:
		var w strings.Builder
		PrintError(&w, e)
		return w.String()
	default:
		return e.Error()
	}
}
