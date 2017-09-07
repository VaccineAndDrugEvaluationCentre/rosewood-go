package types

import (
	"fmt"
	"io"
)

//Table holds all the info needed to render a Table
type Table struct {
	identifier string
	Contents   *TableContents
	grid       *TableContents
	Caption    *Section
	Header     *Section
	Footnotes  *Section
	CmdList    []*Command
}

//NewTable returns a new empty Table
func NewTable() *Table {
	return &Table{}
}

func (t *Table) normalizeMergeRanges() (err error) {
	//trace := utils.NewTrace(true, nil)
	for _, cmd := range t.CmdList {
		if cmd.token != KwMerge {
			continue
		}
		cmd.cellSpan.Normalize(t.Contents.RowCount(), t.Contents.MaxFieldCount())
		//trace.Printf("normalized: %v\n", cmd.cellSpan.TestString())
	}
	return nil
}

func createGridTable(Contents *TableContents, mrlist []Range) (*TableContents, error) {
	grid := NewBlankTableContents(Contents.RowCount(), Contents.MaxFieldCount())
	for _, mr := range mrlist {
		for i := mr.TopLeft.Row + 1; i <= mr.BottomRight.Row; i++ {
			for j := mr.TopLeft.Col + 1; j <= mr.BottomRight.Col; j++ {
				if grid.Cell(i, j).state == CsSpanned {
					return nil, fmt.Errorf("invalid merge range [%s]: it hides a spanned cell [%s]", mr.TestString(), grid.Cell(i, j))
				}
				grid.Cell(i, j).state = CsMerged //hide the other cells in the merge range
			}
		}
		topleft := grid.Cell(mr.TopLeft.Row, mr.TopLeft.Col)
		if topleft.state == CsMerged {
			return nil, fmt.Errorf("invalid merge range [%s]: attempting to span a merged cell [%s]", mr.TestString(), topleft)
		}
		topleft.state = CsSpanned
		topleft.colSpan = mr.BottomRight.Col - mr.TopLeft.Col + 1
		topleft.rowSpan = mr.BottomRight.Row - mr.TopLeft.Row + 1
	}
	//now copy Contents to unmerged cells
	for i := RwInt(1); i <= Contents.RowCount(); i++ {
		r := Contents.Row(i)
		//trace.Printf("in createGridTable i=%d, cellcount=%d\n", i, r.cellCount())
		for j := RwInt(1); j <= r.cellCount(); {
			if Contents.cell(i, j).state != CsMerged {
				grid.cell(i, j).text = Contents.cell(i, j).text
				j++
				continue
			}
			for ; Contents.cell(i, j).state == CsMerged && j <= grid.Row(i).cellCount(); j++ {
				//skip merged cells
			}
			if j <= grid.Row(i).cellCount() {
				grid.cell(i, j).text = Contents.cell(i, j).text
				j++
			}
		}
	}
	return grid, nil
}

func (t *Table) Run() error {
	t.normalizeMergeRanges()
	mrlist, err := createMergeRangeList(t.CmdList)
	if err != nil {
		return err
	}
	t.grid, err = createGridTable(t.Contents, mrlist)
	return err
}

func (t *Table) Render(w io.Writer, hr Renderer) error {
	hr.StartTable(t)
	for _, row := range t.grid.rows {
		hr.StartRow(row)
		for _, cell := range row.cells {
			hr.OutputCell(cell)
		}
		hr.EndRow(row)
	}
	hr.EndTable(t)
	return nil
}

// func searchMRListByRange(mrlist []Range, mr Range) (index int, found bool) {
// 	index = sort.Search(len(mrlist), func(i int) bool { //see https://golang.org/pkg/sort/#Search
// 		return !mrlist[i].less(mr)
// 	})
// 	if index < len(mrlist) && mrlist[index] == mr {
// 		return index, true
// 	}
// 	return -1, false
// }

// func searchMRListByCoordinate(mrlist []Range, coord Coordinates) (index int) {
// 	if len(mrlist) == 0 {
// 		return -1
// 	}
// 	index = sort.Search(len(mrlist), func(i int) bool { //see https://golang.org/pkg/sort/#Search
// 		return mrlist[i].TopLeft.Row >= coord.Row ||
// 			(mrlist[i].TopLeft.Row == coord.Row && mrlist[i].TopLeft.Col >= coord.Col)
// 	})
// 	if index < len(mrlist) && mrlist[index].TopLeft == coord {
// 		return index
// 	}
// 	return -1
// }
