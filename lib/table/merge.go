package table

import (
	"fmt"

	"github.com/drgo/rosewood/lib/types"
)

//FIXME: remove once refactored to add ui.UI
const debug = 0

//createMergedGridTable creates the underlying grid table and applies merging ranges to it
func createMergedGridTable(src *TableContents, mlist []types.Range) (*TableContents, error) {
	grid := NewBlankTableContents(src.RowCount(), src.MaxFieldCount())
	//validate the ranges with respect to this table
	if err := grid.ValidateRanges(mlist); err != nil {
		return nil, err
	}

	for _, mr := range mlist {
		//the topleft cell will hold the row/col span info, so state = csSpanned. Error if it is previously merged or spanned
		topleft := grid.CellorPanic(mr.TopLeft.Row, mr.TopLeft.Col)
		if debug == types.DebugAll {
			fmt.Printf("processing range: %+v\n", mr) //DEBUG
		}
		if topleft.state == CsHMerged || topleft.state == CsVMerged {
			return nil, fmt.Errorf("invalid merge range [%s]: attempting to span a merged cell [%s]", mr.String(), topleft)
		}
		if topleft.state == CsSpanned {
			return nil, fmt.Errorf("invalid merge range [%s]: attempting to span a spanned cell [%s]", mr.String(), topleft)
		}
		topleft.state = CsSpanned
		topleft.colSpan = mr.BottomRight.Col - mr.TopLeft.Col + 1
		topleft.rowSpan = mr.BottomRight.Row - mr.TopLeft.Row + 1
		// find out type of merge range: horinzontal or vertical
		mergeType := CsHMerged
		if topleft.rowSpan > 1 {
			mergeType = CsVMerged
		}
		//hide the other cells in the merge range. Error if it is previously merged or spanned.
		for r := mr.TopLeft.Row; r <= mr.BottomRight.Row; r++ {
			for c := mr.TopLeft.Col; c <= mr.BottomRight.Col; c++ {
				cell := grid.CellorPanic(r, c)
				if cell == topleft {
					continue //skip the spanned cell
				}
				// if the to-be-merged range includes a spanned cell return an error
				if cell.state == CsSpanned {
					return nil, fmt.Errorf("invalid merge range [%s]: it hides a spanned cell [%s]", mr.String(), cell)
				}
				// cell is merged
				cell.state = mergeType
				if debug == types.DebugAll {
					fmt.Printf("cell %d,%d merged, vertically %t \n", r, c, mergeType == CsVMerged) //DEBUG
				}
			}
		}

	}
	if debug == types.DebugAll {
		fmt.Printf("Grid after merging:\n %+v\n", grid.DebugString()) //DEBUG
	}

	//now fill each non-merged cell in the grid with the content of available cells in the raw contents
	if debug == types.DebugAll {
		fmt.Println("copying contents") //DEBUG
	}
	for r := 1; r <= src.RowCount(); r++ {
		if debug == types.DebugAll {
			fmt.Printf("row %d:\n", r) //DEBUG
		}
		if err := copyRowContents(grid, src, r); err != nil {
			return nil, err
		}
	}
	if debug == types.DebugAll {
		fmt.Printf("Grid after copying contents:\n %+v\n", grid.DebugString()) //DEBUG
	}
	return grid, nil
}

func copyRowContents(grid, src *TableContents, r int) error {
	destRowLen := grid.Row(r).cellCount()
	// srcRowLen := src.Row(r).cellCount()
	srcC := 1
	for c := 1; c <= destRowLen; c++ {
		destCell := grid.cell(r, c)
		if destCell.State() == CsHMerged {
			if debug == types.DebugAll {
				fmt.Printf("     skipped horizontally merged cell %d,%d\n", r, c) //DEBUG
			}
			continue
		}
		if !src.isValidCoordinate(r, srcC) { //usually because src row has fewer cells
			break
		}
		srcCell := src.cell(r, srcC)
		destCell.text = srcCell.text
		if debug == types.DebugAll {
			fmt.Printf("     copied cell %d,%d to cell %d,%d\n", r, srcC, r, c) //DEBUG
		}
		srcC++
	}
	return nil
}

//TODO: delete old code
// func copyToValidCell(grid *TableContents, srcCell *Cell) error {
// 	r, c := srcCell.row, srcCell.col
// 	destCell := grid.CellorPanic(r, c)
// 	// if debug == types.DebugAll {
// 	// 	fmt.Printf("copying contents of cell %s", srcCell.DebugString()) //DEBUG
// 	// }
// 	if destCell.State() != CsMerged { //found cell and it's not merged
// 		destCell.text = srcCell.text
// 		if debug == types.DebugAll {
// 			fmt.Printf("     copied to cell %d,%d\n", r, c) //DEBUG
// 		}
// 		return nil
// 	}
// 	//otherwise, find the next non-merged cell in this row
// 	for j := c + 1; j <= grid.Row(r).cellCount(); j++ {
// 		if grid.cell(r, j).state == CsMerged { //find the next non-merged cell if any
// 			if debug == types.DebugAll {
// 				fmt.Printf("     skipped cell %d,%d\n", r, j) //DEBUG
// 			}
// 			continue
// 		}
// 		grid.cell(r, j).text = srcCell.text
// 		if debug == types.DebugAll {
// 			fmt.Printf("     copied to cell %d,%d\n", r, j) //DEBUG
// 		}
// 		return nil
// 	}
// 	return fmt.Errorf("could not copy the contents of cell %d:%d", r, c)
// }

// //now copy the text of each cell in contents to the corresponding cell in grid
// for r := 1; r <= Contents.RowCount(); r++ {
// 	for c := 1; c <= Contents.Row(r).cellCount(); c++ {
// 		srcCell := Contents.cell(r, c)
// 		if debug == types.DebugAll {
// 			fmt.Printf("copying contents of cell %s\n", srcCell.DebugString()) //DEBUG
// 		}
// 		if strings.TrimSpace(srcCell.text) == "" {
// 			if debug == types.DebugAll {
// 				fmt.Printf("     nothing to copy\n") //DEBUG
// 			}
// 			continue
// 		}
// 		if err := copyToValidCell(grid, srcCell); err != nil {
// 			return nil, err
// 		}
// 	}
// }
