// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

import (
	"fmt"
	"strings"
)

//SectionDescriptor describes section type
type SectionDescriptor int

const (
	SectionUnknown SectionDescriptor = iota
	SectionCaption
	SectionBody
	SectionFootNotes
	SectionControl
)

var sectionDescriptorText = [...]string{"Unknown", "Caption", "Body", "FootNotes", "Control"}

//Section holds info on a Rosewood file section
type Section struct {
	Kind   SectionDescriptor
	Offset int
	Lines  []string
}

//NewSection creates a new section
func NewSection(Kind SectionDescriptor, Offset int) *Section {
	return &Section{Kind: Kind, Offset: Offset}
}

//NewControlSection creates a new control section
func NewControlSection(Lines []string) *Section {
	return &Section{Kind: SectionControl, Offset: 1, Lines: Lines}
}

func (s *Section) String() string {
	return strings.Join(s.Lines, "\n")
}

// DebugString prints useful debug info
func (s *Section) DebugString() string {
	return fmt.Sprintf("%s section starts in line %d and ends in line %d", sectionDescriptorText[s.Kind],
		s.Offset, s.Offset+len(s.Lines))
}

//LineCount returns the number of text lines in the section
func (s *Section) LineCount() int {
	return len(s.Lines)
}
