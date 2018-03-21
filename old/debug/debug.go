package debug

import (
	"fmt"
	"io"
	"os"
)

// Stringer is implemented by any value that has a DebugString method, which defines the debug output of the value. The DebugString method is used to print values passed as an operand to any debug print function
type Stringer interface {
	DebugString() string
}

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func Fprint(w io.Writer, title string, s Stringer) (int, error) {
	dstr := "<nil>"
	switch st := s.(type) {
	case Stringer:
		if st != nil {
			dstr = s.DebugString()
		}
	}
	n, err := fmt.Fprintf(w, "%s: %s\n", title, dstr)
	return n, err
}

func Print(title string, s Stringer) (int, error) {
	return Fprint(os.Stdout, title, s)
}
