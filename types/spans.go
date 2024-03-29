// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

import (
	"fmt"
)

//Span holds info on a physical range of a Rosewood table
type Span struct {
	r1, r2, c1, c2 int   //coordinates of the topleft and bottomrow cells in the span, maybe missing
	rby, cby       int   //step increases in row and col
	rcl, ccl       []int //list of row and cols when comma-separated list was specified
}

//NewSpan return new empty Span
func NewSpan() *Span {
	return &Span{r1: RwMissing, r2: RwMissing, c1: RwMissing, c2: RwMissing, rby: RwMissing, cby: RwMissing}
}

//MakeSpan returns a span from topleft and bottomright cell coordinates
func MakeSpan(r1, r2, c1, c2 int) *Span {
	return &Span{r1: r1, r2: r2, c1: c1, c2: c2, rby: RwMissing, cby: RwMissing}
}

//NewSpanFromSpanSegments converts two span segments into one span for ease of validation
func NewSpanFromSpanSegments(spanSegments []*SpanSegment) *Span {
	s := NewSpan()
	for _, segment := range spanSegments {
		switch segment.kind {
		case "row":
			s.r1 = segment.Left
			s.r2 = segment.Right
			s.rby = segment.By
			s.rcl = segment.List
		case "col":
			s.c1 = segment.Left
			s.c2 = segment.Right
			s.cby = segment.By
			s.ccl = segment.List
		default:
			panic("invalid SpanSegment in NewSpanFromSpanSegments()") //should never happen
		}
	}
	return s
}

func SpanToRange(s *Span) Range {
	return makeRange(s.r1, s.c1, s.r2, s.c2)
}

func (s Span) String() string {
	return fmt.Sprintf("row %s:%s col %s:%s", formattedCellCoord(s.r1), formattedCellCoord(s.r2),
		formattedCellCoord(s.c1), formattedCellCoord(s.c2))
}

func (s Span) testString() string {
	return fmt.Sprintf("r(%s,%s)c(%s,%s)by(%s,%s)%v,%v", formattedCellCoord(s.r1), formattedCellCoord(s.r2),
		formattedCellCoord(s.c1), formattedCellCoord(s.c2), formattedCellCoord(s.rby), formattedCellCoord(s.cby), s.rcl, s.ccl)
}

//Validate performs simple validation of the range coordinates
func (s *Span) Validate() error {
	//TODO: add by and rcl validation
	//NOTE: row/col < 1 is prevented by the parser
	if s.r1 > s.r2 {
		return fmt.Errorf("top row number (%d) must be smaller than bottom row number (%d)", s.r1, s.r2)
	}
	if s.c1 > s.c2 {
		return fmt.Errorf("Left column number (%d) must be smaller than Right column number (%d)", s.c1, s.c2)
	}
	return nil
}

//Normalize replace missing values with values defined by rowCount and colCount
func (s *Span) Normalize(rowCount, colCount int) {
	if s.r2 == RwMax { //eg style 1:2:max; max is converted to RwMissing
		s.r2 = RwMissing
	}
	if s.c2 == RwMax {
		s.c2 = RwMissing
	}
	if s.r1 == RwMissing && s.r2 == RwMissing { //span includes all rows eg style col 1
		s.r1 = 1
		s.r2 = rowCount
	}
	if s.c1 == RwMissing && s.c2 == RwMissing { //span includes all cols eg style row 1
		s.c1 = 1
		s.c2 = colCount
	}
	//TODO: verify need for this
	if s.r1 == RwMissing { // row x is equivalent to row x:x,
		s.r1 = s.r2
	}
	if s.r2 == RwMissing { // row x is equivalent to row x:x, eg style row 1 col 1,3
		s.r2 = s.r1
	}
	//TODO: verify need for this
	if s.c1 == RwMissing { // col x is equivalent to col x:x
		s.c1 = s.c2
	}
	if s.c2 == RwMissing { // col x is equivalent to col x:x, eg style row 1,3 col 1
		s.c2 = s.c1
	}
}

//ExpandSpanToRanges convert by and comma list spans into one or more simple (topleft, bottomright) ranges
func (s *Span) ExpandSpanToRanges() (rList []Range, err error) {
	var rPoints, cPoints []int
	//if skipped span (eg 1:2:10), generate Lists of all row and col points included
	if s.rby != RwMissing {
		if rPoints = genAllPossibleRangePoints(s.r1, s.r2, s.rby); rPoints == nil {
			return nil, fmt.Errorf("invalid span %s", s)
		}
	}
	if s.cby != RwMissing {
		if cPoints = genAllPossibleRangePoints(s.c1, s.c2, s.cby); cPoints == nil {
			return nil, fmt.Errorf("invalid span %s", s)
		}
	}
	//if comma-separated, add to above Lists (which could be empty)
	rPoints = append(rPoints, s.rcl...)
	cPoints = append(cPoints, s.ccl...)

	switch {
	//scenario 1: simple (no steps or comma list) span, eg style row 1:3 col 1:2, return it
	case len(rPoints) == 0 && len(cPoints) == 0:
		rList = append(rList, SpanToRange(s))
	//scenario 2: both columns and rows are complex, create a span for each affected row and col combination
	//eg style row 1:2:6 col 1:2:6 --> row 1:1 col 1:1; row 3:3 col 1:1; row 5:5 col 1:1, repeat for col 3:3 & 5:5
	case len(rPoints) != 0 && len(cPoints) != 0:
		for _, r := range rPoints {
			for _, c := range cPoints {
				//TODO: create range directly
				rList = append(rList, SpanToRange(MakeSpan(r, r, c, c)))
			}
		}
	//scenario 3: rows complex but cols simple, create a span for each affected row
	//eg style row 1:2:6 col 1:1 --> row 1:1 col 1:1; row 3:3 col 1:1; row 5:5 col 1:1
	case rPoints != nil:
		for _, r := range rPoints {
			rList = append(rList, SpanToRange(MakeSpan(r, r, s.c1, s.c2)))
		}
	//scenario 4: rows simple but cols complex, create a span for each affected col
	//eg row 1:1 col 1:2:6 --> row 1:1 col 1:1; row 1:1 col 3:3; row 1:1 col 5:5
	case cPoints != nil:
		for _, c := range cPoints {
			rList = append(rList, SpanToRange(MakeSpan(s.r1, s.r2, c, c)))
		}
	}
	return deduplicateRangeList(rList), nil
}

//deduplicateRangeList returns deduplicated (unique topleft, bottomright combo) range List
func deduplicateRangeList(rList []Range) []Range {
	if len(rList) < 2 { //nothing to deduplicate
		return rList
	}
	set := make(map[string]Range, len(rList))
	i := 0
	for _, r := range rList {
		if _, exists := set[r.String()]; exists {
			continue
		}
		set[r.String()] = r
		rList[i] = r
		i++
	}
	return rList[:i]
}

//genAllPossibleRangePoints returns a list of all cell number between p1 and p2 incremented/decremented by step
func genAllPossibleRangePoints(p1, p2, step int) (pList []int) {
	//NOTE: step=zero and abs(step) > p2-p1 is prevented by the parser
	if step > 0 {
		for i := p1; i <= p2; i += step {
			pList = append(pList, i)
		}
	} else {
		for i := p2; i >= p1; i -= step {
			pList = append(pList, i)
		}
	}
	return pList
}
