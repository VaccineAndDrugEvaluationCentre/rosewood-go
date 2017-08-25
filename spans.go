package rosewood

import (
	"fmt"
	"sort"
)

type subspan struct {
	kind        string
	left, right RwInt
	by          RwInt
	list        []RwInt
}

func newSubSpan() subspan {
	return subspan{left: MissingRwInt, right: MissingRwInt, by: MissingRwInt}
}

type span struct {
	r1, r2, c1, c2 RwInt
	rby, cby       RwInt
	rcl, ccl       []RwInt
}

func newSpan() span {
	return span{r1: MissingRwInt, r2: MissingRwInt, c1: MissingRwInt, c2: MissingRwInt, rby: MissingRwInt, cby: MissingRwInt}
}

func makeSpan(r1, r2, c1, c2 RwInt) span {
	return span{r1: r1, r2: r2, c1: c1, c2: c2, rby: MissingRwInt, cby: MissingRwInt}
}

func spanToRange(cs span) Range {
	return makeRange(cs.r1, cs.c1, cs.r2, cs.c2)
}

func (s span) String() string {
	return fmt.Sprintf("row %s:%s col %s:%s", formattedRwInt(s.r1), formattedRwInt(s.r2),
		formattedRwInt(s.c1), formattedRwInt(s.c2))
}

func (s span) testString() string {
	return fmt.Sprintf("r(%s,%s)c(%s,%s)by(%s,%s)%v,%v", formattedRwInt(s.r1), formattedRwInt(s.r2),
		formattedRwInt(s.c1), formattedRwInt(s.c2), formattedRwInt(s.rby), formattedRwInt(s.cby), s.rcl, s.ccl)
}

//validate performs simple validation of the range coordinates
func (s span) validate() error {
	if s.r1 > s.r2 {
		return fmt.Errorf("top row number (%d) must be smaller than bottom row number (%d)", s.r1, s.r2)
	}
	if s.c1 > s.c2 {
		return fmt.Errorf("left column number (%d) must be smaller than right column number (%d)", s.c1, s.c2)
	}
	return nil
}

func normalizeSpan(cs span, rowCount, colCount RwInt) span {
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
	return cs
}

func createMergeRangeList(cmdList []*Command) (rList []Range, err error) {
	var sList []span
	for _, cmd := range cmdList {
		if cmd.token != kwMerge {
			continue
		}
		tmpList, err := expandSpan(cmd.cellSpan)
		if err != nil {
			return nil, err
		}
		sList = append(sList, tmpList...)
	}
	sList = deduplicateSpanList(sList)
	//todo: convert sList to rList
	sort.Slice(rList, func(i, j int) bool {
		return rList[i].less(rList[j])
	})
	return rList, nil
}

func expandSpan(cs span) (sList []span, err error) {
	var rPoints, cPoints []RwInt
	//if skipped span, generate lists of all row and col points included
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
	//if comma-separated, add to above lists (which could be empty)
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
				sList = append(sList, makeSpan(r, MissingRwInt, c, MissingRwInt))
			}
		}
	//scenario 3: rows complex but cols simple
	case rPoints != nil:
		for _, r := range rPoints {
			sList = append(sList, makeSpan(r, MissingRwInt, cs.c1, cs.c2))
		}
	//scenario 4: rows simple but cols complex
	case cPoints != nil:
		for _, c := range cPoints {
			sList = append(sList, makeSpan(cs.r1, cs.r2, c, MissingRwInt))
		}
	}
	return sList, nil
}

//returns deduplicated (unique r1,r2,c1,c2) span list
func deduplicateSpanList(sList []span) []span {
	set := make(map[string]span, len(sList))
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
