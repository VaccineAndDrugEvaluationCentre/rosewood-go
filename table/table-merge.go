package table

import (
	"fmt"

	"github.com/drgo/rosewood/types"
)

//createMergedGridTable creates the underlying grid table and applies merging ranges to it
func (t *Table) createMergedGridTable(mlist []types.Range) error {
	t.grid = NewBlankTableContents(t.Contents.RowCount(), t.Contents.MaxFieldCount())
	//validate the ranges with respect to this table
	if err := t.grid.ValidateRanges(mlist); err != nil {
		return err
	}
	t.Log("creating merged grid") //DEBUG
	for _, mr := range mlist {
		t.Logf("processing range: %+v\n", mr) //DEBUG
		//the topleft cell will hold the row/col span info. Error if it is previously merged or spanned
		topleft := t.grid.CellorPanic(mr.TopLeft.Row, mr.TopLeft.Col)
		if topleft.Merged() {
			return fmt.Errorf("invalid merge range [%s]: attempting to span a merged cell [%s]", mr.String(), topleft)
		}
		if topleft.state == CsSpanned {
			return fmt.Errorf("invalid merge range [%s]: attempting to span a spanned cell [%s]", mr.String(), topleft)
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
				cell := t.grid.CellorPanic(r, c)
				if cell == topleft {
					continue //skip the spanned cell
				}
				// if the to-be-merged range includes a spanned cell return an error
				if cell.state == CsSpanned {
					return fmt.Errorf("invalid merge range [%s]: it hides a spanned cell [%s]", mr.String(), cell)
				}
				// cell is merged
				cell.state = mergeType
				//if vh merge range, copy state from first row except for the first cell of each row
				if mergeType == CsVHMerged && c != mr.TopLeft.Col {
					cell.state = CsHMerged
					if firstRowNum > -1 {
						cell.state = t.grid.CellorPanic(firstRowNum, c).state
					}
				}
				t.Logf("cell %d,%d %s merged:%t\n", r, c, cellStateLabel[cell.state], cell.Merged()) //DEBUG
			}
			if firstRowNum == -1 && mergeType == CsVHMerged { //capture first row number of a vertical/horizontal merge
				firstRowNum = r
			}
		}
	} //loop mlist
	t.Logf("Grid after merging:\n%+v\n", t.grid.DebugString()) //DEBUG
	t.Log("copying contents")                                  //DEBUG

	//now fill each non-merged cell in the grid with the content of available cells in the raw contents
	for r := 1; r <= t.Contents.RowCount(); r++ {
		t.Logf("row %d:\n", r) //DEBUG
		if err := t.copyRowContents(r); err != nil {
			return err
		}
	}
	t.Logf("Grid after copying contents:\n %+v\n", t.grid.DebugString()) //DEBUG
	return nil
}

func (t *Table) copyRowContents(r int) error {
	destRowLen := t.grid.Row(r).cellCount()
	srcC := 1
	for c := 1; c <= destRowLen; c++ {
		destCell := t.grid.cell(r, c)
		if destCell.State() == CsHMerged {
			t.Logf("     skipped horizontally merged cell %d,%d\n", r, c) //DEBUG
			continue
		}
		if !t.Contents.isValidCoordinate(r, srcC) { //usually because src row has fewer cells
			break
		}
		srcCell := t.Contents.cell(r, srcC)
		destCell.text = srcCell.text
		t.Logf("     copied cell %d,%d to cell %d,%d\n", r, srcC, r, c) //DEBUG
		srcC++
	}
	return nil
}
