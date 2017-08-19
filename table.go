package rosewood

import (
	"sort"
)

type mergeRange struct {
	orgRange Range
}

//Table holds all the info needed to render a table
type table struct {
	identifier string
	contents   *tableContents
	grid       *tableContents
	caption    *section
	header     *section
	footnotes  *section
	cmdList    []*Command
	//	mergeList  []mergeRange
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

func createMergeRangeList(cmdList []*Command) (mrlist []mergeRange, err error) {
	for _, cmd := range cmdList {
		if cmd.token != kwMerge {
			continue
		}
		cr := spanToRange(cmd.cellSpan)
		mrlist = append(mrlist, mergeRange{cr})
	}
	sort.Slice(mrlist, func(i, j int) bool {
		return mrlist[i].orgRange.less(mrlist[j].orgRange)
	})
	return mrlist, nil
}

func createGridTable(contents *tableContents, mrlist []mergeRange) (*tableContents, error) {
	grid := newBlankTableContents(contents.rowCount(), contents.maxFldCount)
	for _, mr := range mrlist {
		for i := mr.orgRange.TopLeft.Row; i <= mr.orgRange.BottomRight.Row; i++ {
			for j := mr.orgRange.TopLeft.Col; j <= mr.orgRange.BottomRight.Col; j++ {
				grid.cell(i, j).hidden = true //hide the other cells in the merge range
			}
		}
		topleft := grid.cell(mr.orgRange.TopLeft.Row, mr.orgRange.TopLeft.Col)
		topleft.hidden = false
		topleft.colSpan = mr.orgRange.BottomRight.Col - mr.orgRange.TopLeft.Col + 1
		topleft.rowSpan = mr.orgRange.BottomRight.Row - mr.orgRange.TopLeft.Row + 1
	}
	for i := RwInt(1); i <= contents.rowCount(); i++ {
		r := contents.row(i)
		//fmt.Printf("in createGridTable i=%d, cellcount=%d\n", i, r.cellCount())
		for j := RwInt(1); j <= r.cellCount(); {
			if !contents.cell(i, j).hidden {
				grid.cell(i, j).text = contents.cell(i, j).text
				j++
				continue
			}
			for ; contents.cell(i, j).hidden && j <= grid.row(i).cellCount(); j++ {
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

func (t *table) Merge(ra Range) error {
	//	fmt.Println("range in Merge:", ra)

	//	return t.contents.merge(ra)
	return nil
}

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
// 				bottomRight := mrlist[index].orgRange.BottomRight
// 				//topLeft := mrlist[index].orgRange.TopLeft
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

func searchMRListByRange(mrlist []mergeRange, mr mergeRange) (index int, found bool) {
	index = sort.Search(len(mrlist), func(i int) bool { //see https://golang.org/pkg/sort/#Search
		return !mrlist[i].orgRange.less(mr.orgRange)
	})
	if index < len(mrlist) && mrlist[index].orgRange == mr.orgRange {
		return index, true
	}
	return -1, false
}

func searchMRListByCoordinate(mrlist []mergeRange, coord Coordinates) (index int) {
	if len(mrlist) == 0 {
		return -1
	}
	index = sort.Search(len(mrlist), func(i int) bool { //see https://golang.org/pkg/sort/#Search
		return mrlist[i].orgRange.TopLeft.Row >= coord.Row ||
			(mrlist[i].orgRange.TopLeft.Row == coord.Row && mrlist[i].orgRange.TopLeft.Col >= coord.Col)
	})
	if index < len(mrlist) && mrlist[index].orgRange.TopLeft == coord {
		return index
	}
	return -1
}
