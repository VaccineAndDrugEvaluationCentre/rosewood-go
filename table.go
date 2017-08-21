package rosewood

import (
	"fmt"
	"sort"
)

//Table holds all the info needed to render a table
type table struct {
	identifier string
	contents   *tableContents
	grid       *tableContents
	caption    *section
	header     *section
	footnotes  *section
	cmdList    []*Command
}

func newTable() *table {
	return &table{}
}

func (t *table) normalizeMergeRanges() (err error) {
	for _, cmd := range t.cmdList {
		if cmd.token != kwMerge {
			continue
		}
		cmd.cellSpan = normalizeSpan(cmd.cellSpan, t.contents.rowCount(), t.contents.maxFldCount)
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

func createMergeRangeList(cmdList []*Command) (mrlist []Range, err error) {
	for _, cmd := range cmdList {
		if cmd.token != kwMerge {
			continue
		}
		mrlist = append(mrlist, spanToRange(cmd.cellSpan))
	}
	sort.Slice(mrlist, func(i, j int) bool {
		return mrlist[i].less(mrlist[j])
	})
	return mrlist, nil
}

func createGridTable(contents *tableContents, mrlist []Range) (*tableContents, error) {
	grid := newBlankTableContents(contents.rowCount(), contents.maxFldCount)
	for _, mr := range mrlist {
		for i := mr.TopLeft.Row + 1; i <= mr.BottomRight.Row; i++ {
			for j := mr.TopLeft.Col + 1; j <= mr.BottomRight.Col; j++ {
				if grid.cell(i, j).state == csSpanned {
					return nil, fmt.Errorf("invalid merge range [%s]: it hides a spanned cell [%s]", mr.testString(), grid.cell(i, j))
				}
				grid.cell(i, j).state = csMerged //hide the other cells in the merge range
			}
		}
		topleft := grid.cell(mr.TopLeft.Row, mr.TopLeft.Col)
		if topleft.state == csMerged {
			return nil, fmt.Errorf("invalid merge range [%s]: attempting to span a merged cell [%s]", mr.testString(), topleft)
		}
		topleft.state = csSpanned
		topleft.colSpan = mr.BottomRight.Col - mr.TopLeft.Col + 1
		topleft.rowSpan = mr.BottomRight.Row - mr.TopLeft.Row + 1
	}
	//now copy contents to unmerged cells
	for i := RwInt(1); i <= contents.rowCount(); i++ {
		r := contents.row(i)
		//fmt.Printf("in createGridTable i=%d, cellcount=%d\n", i, r.cellCount())
		for j := RwInt(1); j <= r.cellCount(); {
			if contents.cell(i, j).state != csMerged {
				grid.cell(i, j).text = contents.cell(i, j).text
				j++
				continue
			}
			for ; contents.cell(i, j).state == csMerged && j <= grid.row(i).cellCount(); j++ {
				//skip merged cells
			}
			if j <= grid.row(i).cellCount() {
				grid.cell(i, j).text = contents.cell(i, j).text
				j++
			}
		}
	}
	return grid, nil
}

func (t *table) run() error {
	t.normalizeMergeRanges()
	mrlist, err := createMergeRangeList(t.cmdList)
	if err != nil {
		return err
	}
	t.grid, err = createGridTable(t.contents, mrlist)
	return err
}

// func (t *table) Merge(ra Range) error {
// 	//	fmt.Println("range in Merge:", ra)

// 	//	return t.contents.merge(ra)
// 	return nil
// }

// func createGridTable(contents *tableContents, mrlist []mergeRange) (*tableContents, error) {
// 	grid := newBlankTableContents(contents.rowCount(), contents.maxFldCount)
// 	for i := RwInt(1); i <= contents.rowCount(); i++ {
// 		r := contents.row(i)
// 		//fmt.Printf("in createGridTable i=%d, cellcount=%d\n", i, r.cellCount())
// 		for j := RwInt(1); j <= r.cellCount(); {
// 			grid.cell(i, j).clone(&Cell{row: i, col: j, text: contents.cell(i, j).text})
// 			//fmt.Printf("in createGridTable i=%d, j=%d old %s new %s\n", i, j, contents.cell(i, j).text, grid.cell(i, j).text)
// 			index := searchMRListByCoordinate(mrlist, Coordinates{i, j}) //topleft cell in a merge range?
// 			if index != -1 {
// 				fmt.Printf("in createGridTable index=%d, %+v\n", index, mrlist[index])
// 				bottomRight := mrlist[index].BottomRight
// 				//topLeft := mrlist[index].TopLeft
// 				k := j + 1
// 				for ; k <= bottomRight.Col && k <= grid.row(i).cellCount(); k++ {
// 					grid.cell(i, k).hidden = true //hide the other cells in the merge range
// 				}
// 				//update the row and colspan attributes of the topleft cell
// 				grid.cell(i, j).colSpan = k - j //
// 				//				grid.cell(i, j).rowSpan = bottomRight.Row - topLeft.Row

// 				j = k - 1 //j now is the last processed cell
// 			}
// 			j++
// 		}
// 	}
// 	return grid, nil
// }

func searchMRListByRange(mrlist []Range, mr Range) (index int, found bool) {
	index = sort.Search(len(mrlist), func(i int) bool { //see https://golang.org/pkg/sort/#Search
		return !mrlist[i].less(mr)
	})
	if index < len(mrlist) && mrlist[index] == mr {
		return index, true
	}
	return -1, false
}

func searchMRListByCoordinate(mrlist []Range, coord Coordinates) (index int) {
	if len(mrlist) == 0 {
		return -1
	}
	index = sort.Search(len(mrlist), func(i int) bool { //see https://golang.org/pkg/sort/#Search
		return mrlist[i].TopLeft.Row >= coord.Row ||
			(mrlist[i].TopLeft.Row == coord.Row && mrlist[i].TopLeft.Col >= coord.Col)
	})
	if index < len(mrlist) && mrlist[index].TopLeft == coord {
		return index
	}
	return -1
}
