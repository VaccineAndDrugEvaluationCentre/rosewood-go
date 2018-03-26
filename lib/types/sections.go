// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

import "strings"

type SectionDescriptor int

const (
	SectionUnknown SectionDescriptor = iota
	SectionCaption
	SectionBody
	SectionFootNotes
	SectionControl
)

type Section struct {
	//	scriptIdentifer string
	Kind   SectionDescriptor
	Offset int
	Lines  []string
}

func NewSection(Kind SectionDescriptor, Offset int) *Section {
	return &Section{Kind: Kind, Offset: Offset}
}

//newControlSection mostly used for
func NewControlSection(Lines []string) *Section {
	return &Section{Kind: SectionControl, Offset: 1, Lines: Lines}
}

func (s *Section) String() string {
	return strings.Join(s.Lines, "\n")
}

func (s *Section) LineCount() int {
	return len(s.Lines)
}
