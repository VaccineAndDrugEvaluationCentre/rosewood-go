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
		cmd.cellSpan.normalize(t.contents.rowCount(), t.contents.maxFldCount)
		trace.Printf("normalized: %v\n", cmd.cellSpan.testString())
	}
	return nil
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
		//trace.Printf("in createGridTable i=%d, cellcount=%d\n", i, r.cellCount())
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
