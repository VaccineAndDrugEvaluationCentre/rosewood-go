package table

import (
	"fmt"

	"github.com/drgo/rosewood/types"
)

//FIXME: remove once refactored to add ui.UI
const debug = types.DebugAll

//createMergedGridTable creates the underlying grid table and applies merging ranges to it
func createMergedGridTable(src *TableContents, mlist []types.Range) (*TableContents, error) {
	grid := NewBlankTableContents(src.RowCount(), src.MaxFieldCount())
	//validate the ranges with respect to this table
	if err := grid.ValidateRanges(mlist); err != nil {
		return nil, err
	}
	if debug == types.DebugAll {
		fmt.Println("creating merged grid") //DEBUG
	}

	for _, mr := range mlist {
		if debug == types.DebugAll {
			fmt.Printf("processing range: %+v\n", mr) //DEBUG
		}
		//the topleft cell will hold the row/col span info. Error if it is previously merged or spanned
		topleft := grid.CellorPanic(mr.TopLeft.Row, mr.TopLeft.Col)
		if topleft.Merged() {
			return nil, fmt.Errorf("invalid merge range [%s]: attempting to span a merged cell [%s]", mr.String(), topleft)
		}
		if topleft.state == CsSpanned {
			return nil, fmt.Errorf("invalid merge range [%s]: attempting to span a spanned cell [%s]", mr.String(), topleft)
		}
		// determine topcell state
		topleft.state = CsSpanned
		topleft.colSpan = mr.BottomRight.Col - mr.TopLeft.Col + 1
		topleft.rowSpan = mr.BottomRight.Row - mr.TopLeft.Row + 1
		// find out type of merge range: horizontal or vertical
		mergeType := CsNormal
		switch {
		case topleft.rowSpan > 1 && topleft.colSpan > 1:
			mergeType = CsVHMerged
		case topleft.rowSpan > 1:
			mergeType = CsVMerged
		case topleft.colSpan > 1:
			mergeType = CsHMerged
		}
		//hide the other cells in the merge range. Error if a cell in the range is previously merged or spanned.
		firstRowNum := -1
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
				//if vh merge range, copy state from first row except for the first cell of each row
				if mergeType == CsVHMerged && c != mr.TopLeft.Col {
					cell.state = CsHMerged
					if firstRowNum > -1 {
						cell.state = grid.CellorPanic(firstRowNum, c).state
					}
				}
				if debug == types.DebugAll {
					fmt.Printf("cell %d,%d %s merged:%t\n", r, c, cellStateLabel[cell.state], cell.Merged()) //DEBUG
				}
			}
			if firstRowNum == -1 && mergeType == CsVHMerged { //capture first row number of a vertical/horizontal merge
				firstRowNum = r
			}
		}
	} //loop mlist
	if debug == types.DebugAll {
		fmt.Printf("Grid after merging:\n%+v\n", grid.DebugString()) //DEBUG
		fmt.Println("copying contents")                              //DEBUG
	}
	//now fill each non-merged cell in the grid with the content of available cells in the raw contents
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
