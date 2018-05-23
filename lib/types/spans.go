// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

import (
	"bytes"
	"fmt"
)

type Subspan struct {
	kind        string //
	Left, Right RwInt
	By          RwInt
	List        []RwInt
}

func NewSubSpan(kind string) Subspan {
	return Subspan{kind: kind, Left: MissingRwInt, Right: MissingRwInt, By: MissingRwInt}
}

func (ss *Subspan) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%s ", ss.kind)
	if ss.Left != MissingRwInt {
		fmt.Fprintf(buf, "%s", formattedRwInt(ss.Left))
		if ss.By != MissingRwInt {
			fmt.Fprintf(buf, ":%s", formattedRwInt(ss.By))
		}
		fmt.Fprintf(buf, ":%s", formattedRwInt(ss.Right))
		if len(ss.List) > 0 {
			fmt.Fprintf(buf, ", ") //add comma if we also have a comma separated list
		}
	}
	for _, item := range ss.List {
		fmt.Fprintf(buf, "%s,", formattedRwInt(item))
	}
	//remove last comma if any
	if bytes.HasSuffix(buf.Bytes(), []byte{','}) {
		buf.Truncate(buf.Len() - 1)
	}
	return buf.String()
}

//SubSpansToSpan converts two subspan into one span for ease of validation
func SubSpansToSpan(SubSpans []*Subspan) *Span {
	s := NewSpan()
	for _, ss := range SubSpans {
		switch ss.kind {
		case "row":
			s.r1 = ss.Left
			s.r2 = ss.Right
			s.rby = ss.By
			s.rcl = ss.List
		case "col":
			s.c1 = ss.Left
			s.c2 = ss.Right
			s.cby = ss.By
			s.ccl = ss.List
		default:
			panic("invalid Subspan in SubspanToSpan()") //should never happen
		}
	}
	return s
}

type Span struct {
	r1, r2, c1, c2 RwInt
	rby, cby       RwInt
	rcl, ccl       []RwInt
}

func NewSpan() *Span {
	return &Span{r1: MissingRwInt, r2: MissingRwInt, c1: MissingRwInt, c2: MissingRwInt, rby: MissingRwInt, cby: MissingRwInt}
}

func MakeSpan(r1, r2, c1, c2 RwInt) *Span {
	return &Span{r1: r1, r2: r2, c1: c1, c2: c2, rby: MissingRwInt, cby: MissingRwInt}
}

func SpanToRange(cs *Span) Range {
	return makeRange(cs.r1, cs.c1, cs.r2, cs.c2)
}

func (s Span) String() string {
	return fmt.Sprintf("row %s:%s col %s:%s", formattedRwInt(s.r1), formattedRwInt(s.r2),
		formattedRwInt(s.c1), formattedRwInt(s.c2))
}

func (s Span) testString() string {
	return fmt.Sprintf("r(%s,%s)c(%s,%s)by(%s,%s)%v,%v", formattedRwInt(s.r1), formattedRwInt(s.r2),
		formattedRwInt(s.c1), formattedRwInt(s.c2), formattedRwInt(s.rby), formattedRwInt(s.cby), s.rcl, s.ccl)
}

//Validate performs simple validation of the range coordinates
func (s *Span) Validate() error {
	if s.r1 > s.r2 {
		return fmt.Errorf("top row number (%d) must be smaller than bottom row number (%d)", s.r1, s.r2)
	}
	if s.c1 > s.c2 {
		return fmt.Errorf("Left column number (%d) must be smaller than Right column number (%d)", s.c1, s.c2)
	}
	return nil
}

//Normalize replace missing values with values defined by rowCount and colCount
func (cs *Span) Normalize(rowCount, colCount RwInt) {
	if cs.r1 == MissingRwInt && cs.r2 == MissingRwInt {
		cs.r1 = 1
		cs.r2 = rowCount
	}
	if cs.c1 == MissingRwInt && cs.c2 == MissingRwInt {
		cs.c1 = 1
		cs.c2 = colCount
	}
	if cs.r1 == MissingRwInt {
		cs.r1 = cs.r2
	}
	if cs.r2 == MissingRwInt {
		cs.r2 = cs.r1
	}
	if cs.c1 == MissingRwInt {
		cs.c1 = cs.c2
	}
	if cs.c2 == MissingRwInt {
		cs.c2 = cs.c1
	}
}

//ExpandSpan convert by and comma list spans into simple spans
func (cs *Span) ExpandSpan() (sList []*Span, err error) {
	var rPoints, cPoints []RwInt
	//if skipped span, generate Lists of all row and col points included
	if cs.rby != MissingRwInt {
		if rPoints = genAllPossibleRangePoints(cs.r1, cs.r2, cs.rby); rPoints == nil {
			return nil, fmt.Errorf("invalid span %s", cs)
		}
	}
	if cs.cby != MissingRwInt {
		if cPoints = genAllPossibleRangePoints(cs.c1, cs.c2, cs.cby); cPoints == nil {
			return nil, fmt.Errorf("invalid span %s", cs)
		}
	}
	//if comma-separated, add to above Lists (which could be empty)
	rPoints = append(rPoints, cs.rcl...)
	cPoints = append(cPoints, cs.ccl...)

	switch {
	//scenario 1: simple span, return it
	case len(rPoints) == 0 && len(cPoints) == 0:
		sList = append(sList, cs)
	//scenario 2: both columns and rows are complex
	case len(rPoints) != 0 && len(cPoints) != 0:
		for _, r := range rPoints {
			for _, c := range cPoints {
				sList = append(sList, MakeSpan(r, MissingRwInt, c, MissingRwInt))
			}
		}
	//scenario 3: rows complex but cols simple
	case rPoints != nil:
		for _, r := range rPoints {
			sList = append(sList, MakeSpan(r, MissingRwInt, cs.c1, cs.c2))
		}
	//scenario 4: rows simple but cols complex
	case cPoints != nil:
		for _, c := range cPoints {
			sList = append(sList, MakeSpan(cs.r1, cs.r2, c, MissingRwInt))
		}
	}
	return sList, nil
}

//DeduplicateSpanList returns deduplicated (unique r1,r2,c1,c2) span List
func DeduplicateSpanList(sList []*Span) []*Span {
	set := make(map[string]*Span, len(sList))
	i := 0
	for _, cs := range sList {
		if _, exists := set[cs.String()]; exists {
			continue
		}
		set[cs.String()] = cs
		sList[i] = cs
		i++
	}
	return sList[:i]
}

func genAllPossibleRangePoints(p1, p2, by RwInt) (pList []RwInt) {
	if by == 0 || by > p2-p1 {
		return nil
	}
	for i := p1; i <= p2; i += by {
		pList = append(pList, i)
	}
	return pList
}
