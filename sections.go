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
	scriptIdentifer string
	offset          int
	kind            sectionDescriptor
	lines           []string
}

//todo: change section to byte.buffer or []bytes
func newSection(scriptIdentifer string, offset int, kind sectionDescriptor) *section {
	return &section{scriptIdentifer: scriptIdentifer, offset: offset, kind: kind}
}

func (s *section) String() string {
	return strings.Join(s.lines, "\n")
}

func (s *section) LineCount() int {
	return len(s.lines)
}
