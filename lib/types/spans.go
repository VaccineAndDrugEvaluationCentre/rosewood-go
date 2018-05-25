// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

import (
	"fmt"
	"sort"
)

//Span holds info on a physical range of a Rosewood table
type Span struct {
	r1, r2, c1, c2 RwInt   //coordinates of the topleft and bottomrow cells in the span, maybe missing
	rby, cby       RwInt   //step increases in row and col
	rcl, ccl       []RwInt //list of row and cols when comma-separated list was specified
}

//NewSpan return new empty Span
func NewSpan() *Span {
	return &Span{r1: MissingRwInt, r2: MissingRwInt, c1: MissingRwInt, c2: MissingRwInt, rby: MissingRwInt, cby: MissingRwInt}
}

//MakeSpan returns a span from topleft and bottomright cell coordinates
func MakeSpan(r1, r2, c1, c2 RwInt) *Span {
	return &Span{r1: r1, r2: r2, c1: c1, c2: c2, rby: MissingRwInt, cby: MissingRwInt}
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
	return fmt.Sprintf("row %s:%s col %s:%s", formattedRwInt(s.r1), formattedRwInt(s.r2),
		formattedRwInt(s.c1), formattedRwInt(s.c2))
}

func (s Span) testString() string {
	return fmt.Sprintf("r(%s,%s)c(%s,%s)by(%s,%s)%v,%v", formattedRwInt(s.r1), formattedRwInt(s.r2),
		formattedRwInt(s.c1), formattedRwInt(s.c2), formattedRwInt(s.rby), formattedRwInt(s.cby), s.rcl, s.ccl)
}

//Validate performs simple validation of the range coordinates
func (s *Span) Validate() error {
	//TODO: add by and rcl validation
	//NOTE: row/col=0 is prevented by the parser
	if s.r1 > s.r2 {
		return fmt.Errorf("top row number (%d) must be smaller than bottom row number (%d)", s.r1, s.r2)
	}
	if s.c1 > s.c2 {
		return fmt.Errorf("Left column number (%d) must be smaller than Right column number (%d)", s.c1, s.c2)
	}
	return nil
}

//Normalize replace missing values with values defined by rowCount and colCount
func (s *Span) Normalize(rowCount, colCount RwInt) {
	if s.r1 == MissingRwInt && s.r2 == MissingRwInt { //span includes all rows
		s.r1 = 1
		s.r2 = rowCount
	}
	if s.c1 == MissingRwInt && s.c2 == MissingRwInt { //span includes all cols
		s.c1 = 1
		s.c2 = colCount
	}
	if s.r1 == MissingRwInt { // row x is equivalent to row x:x
		s.r1 = s.r2
	}
	if s.r2 == MissingRwInt { // row x is equivalent to row x:x
		s.r2 = s.r1
	}
	if s.c1 == MissingRwInt { // col x is equivalent to col x:x
		s.c1 = s.c2
	}
	if s.c2 == MissingRwInt { // col x is equivalent to col x:x
		s.c2 = s.c1
	}
}

//ExpandSpanToRanges convert by and comma list spans into one or more simple (topleft, bottomright) ranges
func (s *Span) ExpandSpanToRanges() (rList []Range, err error) {
	var rPoints, cPoints []RwInt
	//if skipped span (eg 1:2:10), generate Lists of all row and col points included
	if s.rby != MissingRwInt {
		if rPoints = genAllPossibleRangePoints(s.r1, s.r2, s.rby); rPoints == nil {
			return nil, fmt.Errorf("invalid span %s", s)
		}
	}
	if s.cby != MissingRwInt {
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

//TODO: optimize
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

//genAllPossibleRangePoints returns a list of all cell number between p1 and p2 incremented by step
func genAllPossibleRangePoints(p1, p2, step RwInt) (pList []RwInt) {
	if step == 0 || step > p2-p1 {
		return nil
	}
	for i := p1; i <= p2; i += step {
		pList = append(pList, i)
	}
	return pList
}

//TODO: Add to types.Table
//spanToRangeList converts the spans specified in each command of type cmdType into a list of Type.Range ready for use
func getAllRanges(cmdList []*Command, cmdType RwKeyWord) (allRangesList []Range, err error) {
	for _, cmd := range cmdList {
		if cmd.token != cmdType {
			continue
		}
		rList, err := cmd.cellSpan.ExpandSpanToRanges()
		fmt.Printf("%d %s --> len(rlist)=%d \n", cmdType, cmd.cellSpan, len(rList))
		if err != nil {
			return nil, err
		}
		//attach styles to all ranges
		for i := range rList {
			if cmdType == KwStyle {
				rList[i].addStyle(cmd.Args()...)
			}
		}
		allRangesList = append(allRangesList, rList...)
	}
	sort.Slice(allRangesList, func(i, j int) bool {
		return allRangesList[i].less(allRangesList[j])
	})
	return allRangesList, nil
}

//TODO: remove OLD CODE
//deduplicateSpanList returns deduplicated (unique r1,r2,c1,c2) span List
func deduplicateSpanList(sList []*Span) []*Span {
	set := make(map[string]*Span, len(sList))
	i := 0
	for _, s := range sList {
		if _, exists := set[s.String()]; exists {
			continue
		}
		set[s.String()] = s
		sList[i] = s
		i++
	}
	return sList[:i]
}

//ExpandSpan convert by and comma list spans into simple (topleft, bottomright only) spans
func (s *Span) ExpandSpan() (sList []*Span, err error) {
	var rPoints, cPoints []RwInt
	//if skipped span (eg 1:2:10), generate Lists of all row and col points included
	if s.rby != MissingRwInt {
		if rPoints = genAllPossibleRangePoints(s.r1, s.r2, s.rby); rPoints == nil {
			return nil, fmt.Errorf("invalid span %s", s)
		}
	}
	if s.cby != MissingRwInt {
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
		sList = append(sList, s)
	//scenario 2: both columns and rows are complex, create a span for each affected row and col combination
	//eg style row 1:2:6 col 1:2:6 --> row 1:1 col 1:1; row 3:3 col 1:1; row 5:5 col 1:1, repeat for col 3:3 & 5:5
	case len(rPoints) != 0 && len(cPoints) != 0:
		for _, r := range rPoints {
			for _, c := range cPoints {
				sList = append(sList, MakeSpan(r, r, c, c))
			}
		}
	//scenario 3: rows complex but cols simple, create a span for each affected row
	//eg style row 1:2:6 col 1:1 --> row 1:1 col 1:1; row 3:3 col 1:1; row 5:5 col 1:1
	case rPoints != nil:
		for _, r := range rPoints {
			sList = append(sList, MakeSpan(r, r, s.c1, s.c2))
		}
	//scenario 4: rows simple but cols complex, create a span for each affected col
	//eg row 1:1 col 1:2:6 --> row 1:1 col 1:1; row 1:1 col 3:3; row 1:1 col 5:5
	case cPoints != nil:
		for _, c := range cPoints {
			sList = append(sList, MakeSpan(s.r1, s.r2, c, c))
		}
	}
	return deduplicateSpanList(sList), nil
}
