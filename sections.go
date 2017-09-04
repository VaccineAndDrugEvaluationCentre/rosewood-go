package rosewood

import "strings"

type sectionDescriptor int

const (
	sectionUnknown sectionDescriptor = iota
	sectionCaption
	sectionBody
	sectionFootNotes
	sectionControl
)

type section struct {
	//	scriptIdentifer string
	kind   sectionDescriptor
	offset int
	lines  []string
}

func newSection(kind sectionDescriptor, offset int) *section {
	return &section{kind: kind, offset: offset}
}

func (s *section) String() string {
	return strings.Join(s.lines, "\n")
}

func (s *section) LineCount() int {
	return len(s.lines)
}
