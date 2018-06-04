package types

import (
	"fmt"
	"strings"
)

//createMergedGridTable creates the underlying grid table and applies merging ranges to it
func createMergedGridTable(Contents *TableContents, mrlist []Range) (*TableContents, error) {
	grid := NewBlankTableContents(Contents.RowCount(), Contents.MaxFieldCount())
	//validate the ranges with respect to this table
	if err := grid.ValidateRanges(mrlist); err != nil {
		return nil, err
	}
	//fmt.Printf("%+v\n", grid.DebugString())
	for _, mr := range mrlist {
		//fmt.Printf("%+v\n", mr)  //DEBUG
		//the topleft cell will hold the row/col span info so state = csSpanned. Error if it is previously merged or spanned
		topleft := grid.CellorPanic(mr.TopLeft.Row, mr.TopLeft.Col)
		if topleft.state == CsMerged {
			return nil, fmt.Errorf("invalid merge range [%s]: attempting to span a merged cell [%s]", mr.testString(), topleft)
		}
		if topleft.state == CsSpanned {
			return nil, fmt.Errorf("invalid merge range [%s]: attempting to span a spanned cell [%s]", mr.testString(), topleft)
		}
		topleft.state = CsSpanned
		topleft.colSpan = mr.BottomRight.Col - mr.TopLeft.Col + 1
		topleft.rowSpan = mr.BottomRight.Row - mr.TopLeft.Row + 1
		//hide the other cells in the merge range. Error if it is previously merged or spanned.
		for i := mr.TopLeft.Row; i <= mr.BottomRight.Row; i++ {
			for j := mr.TopLeft.Col; j <= mr.BottomRight.Col; j++ {
				cell := grid.CellorPanic(i, j)
				if cell == topleft {
					continue //skip the spanned cell
				}
				if cell.state == CsSpanned {
					return nil, fmt.Errorf("invalid merge range [%s]: it hides a spanned cell [%s]", mr.testString(), cell)
				}
				cell.state = CsMerged
				//fmt.Printf("%d:%d= %s\n", i, j, CellStateLabel[cell.State()]) //DEBUG
			}
		}

	}
	//fmt.Printf("Contents:\n %+v\n", Contents.DebugString()) //DEBUG
	//fmt.Printf("Destination:\n %+v\n", grid.DebugString())  //DEBUG
	/*
	   Contents:
	    1,1 (normal) |
	   2,1 (normal) |2,2 (normal) |2,3 (normal) |
	   3,1 (normal) |3,2 (normal) |3,3 (normal) |3,4 (normal) |
	   4,1 (normal) |4,2 (normal) |4,3 (normal) |4,4 (normal) |
	   5,1 (normal) |5,2 (normal) |5,3 (normal) |5,4 (normal) |
	   6,1 (normal) |6,2 (normal) |6,3 (normal) |6,4 (normal) |

	   Destination:
	    1,1 (spanned) |1,2 (merged) |1,3 (merged) |1,4 (merged) |
	   2,1 (spanned) |2,2 (spanned) |2,3 (merged) |2,4 (spanned) |
	   3,1 (merged) |3,2 (normal) |3,3 (normal) |3,4 (merged) |
	   4,1 (normal) |4,2 (normal) |4,3 (normal) |4,4 (normal) |
	   5,1 (normal) |5,2 (normal) |5,3 (normal) |5,4 (normal) |
	   6,1 (normal) |6,2 (normal) |6,3 (normal) |6,4 (normal) |
	*/
	//now copy the text of each cell in contents to the corresponding cell in grid
	//skipping any merged cells
	for i := 1; i <= Contents.RowCount(); i++ {
		for j := 1; j <= Contents.Row(i).cellCount(); j++ {
			if strings.TrimSpace(Contents.cell(i, j).text) == "" {
				continue
			}
			if err := copyToValidCell(grid, Contents.cell(i, j)); err != nil {
				return nil, err
			}
		}
	}
	return grid, nil
}

//TODO: merge if and for loop; move to TableContents
func copyToValidCell(grid *TableContents, srcCell *Cell) error {
	r, c := srcCell.row, srcCell.col
	destCell := grid.CellorPanic(r, c)
	if destCell.State() != CsMerged { //found cell and it's not merged
		destCell.text = srcCell.text
		return nil
	}
	//otherwise, find the next non-merged cell in this row
	for j := c + 1; j <= grid.Row(r).cellCount(); j++ {
		if grid.cell(r, j).state == CsMerged { //find the next non-merged cell if any
			continue
		}
		grid.cell(r, j).text = srcCell.text
		return nil
	}
	return fmt.Errorf("could not copy the contents of cell %d:%d", r, c)
}
