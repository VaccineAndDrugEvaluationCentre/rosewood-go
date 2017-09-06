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

//newControlSection mostly used for
func newControlSection(lines []string) *section {
	return &section{kind: sectionControl, offset: 1, lines: lines}
}

func (s *section) String() string {
	return strings.Join(s.lines, "\n")
}

func (s *section) LineCount() int {
	return len(s.lines)
}
