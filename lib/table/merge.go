package table

import (
	"fmt"
	"strings"

	"github.com/drgo/rosewood/lib/types"
)

//createMergedGridTable creates the underlying grid table and applies merging ranges to it
func createMergedGridTable(Contents *TableContents, mrlist []types.Range) (*TableContents, error) {
	grid := NewBlankTableContents(Contents.RowCount(), Contents.MaxFieldCount())
	//validate the ranges with respect to this table
	if err := grid.ValidateRanges(mrlist); err != nil {
		return nil, err
	}
	//fmt.Printf("%+v\n", grid.DebugString())
	for _, mr := range mrlist {
		//the topleft cell will hold the row/col span info so state = csSpanned. Error if it is previously merged or spanned
		topleft := grid.CellorPanic(mr.TopLeft.Row, mr.TopLeft.Col)
		if debug == types.DebugAll {
			fmt.Printf("processing range: %+v\n", mr) //DEBUG
		}
		if topleft.state == CsMerged {
			return nil, fmt.Errorf("invalid merge range [%s]: attempting to span a merged cell [%s]", mr.String(), topleft)
		}
		if topleft.state == CsSpanned {
			return nil, fmt.Errorf("invalid merge range [%s]: attempting to span a spanned cell [%s]", mr.String(), topleft)
		}
		topleft.state = CsSpanned
		topleft.colSpan = mr.BottomRight.Col - mr.TopLeft.Col + 1
		topleft.rowSpan = mr.BottomRight.Row - mr.TopLeft.Row + 1
		//hide the other cells in the merge range. Error if it is previously merged or spanned.
		for r := mr.TopLeft.Row; r <= mr.BottomRight.Row; r++ {
			for c := mr.TopLeft.Col; c <= mr.BottomRight.Col; c++ {
				cell := grid.CellorPanic(r, c)
				if cell == topleft {
					continue //skip the spanned cell
				}
				if cell.state == CsSpanned {
					return nil, fmt.Errorf("invalid merge range [%s]: it hides a spanned cell [%s]", mr.String(), cell)
				}
				cell.state = CsMerged
				if debug == types.DebugAll {
					fmt.Printf("cell %d,%d status= merged\n", r, c) //DEBUG
				}
			}
		}

	}
	if debug == types.DebugAll {
		fmt.Printf("Grid after merging:\n %+v\n", grid.DebugString()) //DEBUG
	}
	//now copy the text of each cell in contents to the corresponding cell in grid
	for r := 1; r <= Contents.RowCount(); r++ {
		for c := 1; c <= Contents.Row(r).cellCount(); c++ {
			srcCell := Contents.cell(r, c)
			if debug == types.DebugAll {
				fmt.Printf("copying contents of cell %s\n", srcCell.DebugString()) //DEBUG
			}
			if strings.TrimSpace(srcCell.text) == "" {
				if debug == types.DebugAll {
					fmt.Printf("     nothing to copy\n") //DEBUG
				}
				continue
			}
			if err := copyToValidCell(grid, srcCell); err != nil {
				return nil, err
			}
		}
	}
	if debug == types.DebugAll {
		fmt.Printf("Grid after copying contents:\n %+v\n", grid.DebugString()) //DEBUG
	}
	return grid, nil
}

//TODO: move to TableContents?
func copyToValidCell(grid *TableContents, srcCell *Cell) error {
	r, c := srcCell.row, srcCell.col
	destCell := grid.CellorPanic(r, c)
	// if debug == types.DebugAll {
	// 	fmt.Printf("copying contents of cell %s", srcCell.DebugString()) //DEBUG
	// }
	if destCell.State() != CsMerged { //found cell and it's not merged
		destCell.text = srcCell.text
		if debug == types.DebugAll {
			fmt.Printf("     copied to cell %d,%d\n", r, c) //DEBUG
		}
		return nil
	}
	//otherwise, find the next non-merged cell in this row
	for j := c + 1; j <= grid.Row(r).cellCount(); j++ {
		if grid.cell(r, j).state == CsMerged { //find the next non-merged cell if any
			if debug == types.DebugAll {
				fmt.Printf("     skipped cell %d,%d\n", r, j) //DEBUG
			}
			continue
		}
		grid.cell(r, j).text = srcCell.text
		if debug == types.DebugAll {
			fmt.Printf("     copied to cell %d,%d\n", r, j) //DEBUG
		}
		return nil
	}
	return fmt.Errorf("could not copy the contents of cell %d:%d", r, c)
}
