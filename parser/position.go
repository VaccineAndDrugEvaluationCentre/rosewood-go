package parser

import (
	"fmt"
	"text/scanner"
)

// Position represents a source position
// encapsulates go's scanner.Position and overrides its methods
// A position is valid if Line > 0 and Offset > 0.
type Position struct {
	scanner.Position
}

// IsValid reports whether the position is valid.
func (pos *Position) IsValid() bool { return pos.Offset > 0 && pos.Line > 0 }

func (pos Position) String() string {
	s := pos.Filename
	if s == "" {
		s = "<input>"
	}
	if pos.IsValid() {
		s += fmt.Sprintf(":%d:%d", pos.AdjLine(), pos.Column)
	}
	return s
}

//AdjLine returns an adjusted line number
func (pos *Position) AdjLine() int {
	if pos.IsValid() {
		return pos.Offset + pos.Line
	}
	return 0
}
