// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

import (
	"io"
	"strings"
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

//TODO: add debugstringer interface with debugString() and title() functions to use for debugPrint
func (t Table) String() string {
	var s strings.Builder
	s.Grow(1024 * 8) //arbitrary size to avoid reallocation
	if t.Caption != nil {
		s.WriteString("caption: " + t.Caption.String() + "\n")
	}
	if t.Header != nil {
		s.WriteString("header: " + t.Header.String() + "\n")
	}
	if t.Footnotes != nil {
		s.WriteString("footnotes: " + t.Footnotes.String() + "\n")
	}
	if t.grid != nil {
		s.WriteString("raw content: " + t.grid.DebugString() + "\n")
	}
	if t.Contents != nil {
		s.WriteString("proceeded content: " + t.Contents.DebugString() + "\n")
	}
	return s.String()
}

//Run applies all commands to table contents. Must be called before rendering the table
func (t *Table) Run() error {
	t.fixMissingRangeValues()
	//create a list of merge ranges
	mrlist, err := getAllRanges(t.CmdList, KwMerge)
	if err != nil {
		return err
	}
	t.grid, err = createMergedGridTable(t.Contents, mrlist)
	if err != nil {
		return err
	}
	//create a list of style ranges
	mrlist, err = getAllRanges(t.CmdList, KwStyle)
	if err != nil {
		return err
	}
	t.grid, err = applyStyles(t.grid, mrlist)
	return err
}

//Render use a types.Renderer to render table contents and write them to io.Writer
func (t *Table) Render(w io.Writer, hr Renderer) error {
	if err := hr.StartTable(t); err != nil {
		return err
	}
	for _, row := range t.grid.rows {
		if err := hr.StartRow(row); err != nil {
			return err
		}

		for _, cell := range row.cells {
			if err := hr.OutputCell(cell); err != nil {
				return err
			}
		}
		if err := hr.EndRow(row); err != nil {
			return err
		}

	}
	return hr.EndTable(t)

}

//fixMissingRangeValues fixes missing coordinates with reference to this table's dimensions
func (t *Table) fixMissingRangeValues() (err error) {
	for _, cmd := range t.CmdList {
		if !IsTableCommand(cmd.token) {
			continue
		}
		cmd.cellSpan.Normalize(t.Contents.RowCount(), t.Contents.MaxFieldCount())
	}
	return nil
}

func applyStyles(Contents *TableContents, mrlist []Range) (*TableContents, error) {
	if err := Contents.ValidateRanges(mrlist); err != nil {
		return nil, err
	}
	for _, mr := range mrlist {
		for i := mr.TopLeft.Row; i <= mr.BottomRight.Row; i++ {
			for j := mr.TopLeft.Col; j <= mr.BottomRight.Col; j++ {
				Contents.CellorPanic(i, j).AddStyle(mr.styles()...)
			}
		}
	}
	return Contents, nil
}

//TODO: remove unused code
// //spanToRangeList converts the spans specified in each command into a list of Type.Range ready for use
// func oldspanToRangeList(cmdList []*Command, cmdType RwKeyWord) (rList []Range, err error) {
// 	var sList []*Span
// 	for _, cmd := range cmdList {
// 		if cmd.token != cmdType {
// 			continue
// 		}
// 		tmpList, err := cmd.cellSpan.ExpandSpan()
// 		if err != nil {
// 			return nil, err
// 		}
// 		sList = append(sList, tmpList...)
// 	}
// 	sList = DeduplicateSpanList(sList)
// 	for _, s := range sList {
// 		rList = append(rList, SpanToRange(s))
// 	}
// 	sort.Slice(rList, func(i, j int) bool {
// 		return rList[i].less(rList[j])
// 	})
// 	return rList, nil
// }

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
