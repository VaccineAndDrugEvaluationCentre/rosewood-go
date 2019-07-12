// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package table

import (
	"fmt"
	"io"
	"strings"

	"github.com/drgo/core/ui"
	"github.com/drgo/rosewood/types"
)

//Table holds all the info needed to render a Table
type Table struct {
	ui.UI
	identifier string
	Contents   *TableContents // source grid
	grid       *TableContents //output grid
	Caption    *types.Section
	Header     *types.Section
	Footnotes  *types.Section
	CmdList    []*types.Command
}

//NewTable returns a new empty Table
func NewTable(debug int) *Table {
	return &Table{
		UI: ui.NewUI(debug),
	}
}

//ProcessedTableContents returns a pointer to table contents after applying all commands
func (t *Table) ProcessedTableContents() *TableContents {
	return t.grid
}

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
	if t.Contents != nil {
		s.WriteString("raw content: " + "\n" + t.Contents.DebugString() + "\n")
	}
	if t.grid != nil {
		s.WriteString("processed content: " + "\n" + t.grid.DebugString() + "\n")
	}
	return s.String()
}

//Run applies all commands to table contents. Must be called before rendering the table
func (t *Table) Run() error {
	t.fixMissingRangeValues()
	//create a list of merge ranges
	rlist, err := types.GetAllRanges(t.CmdList, types.KwMerge)
	if err != nil {
		return err
	}
	if t.Level() == ui.DebugAll {
		fmt.Println("Table.Run()- Merged ranges:")
		for _, r := range rlist {
			fmt.Printf("%v\n", r)
		}
	}
	if err = t.createMergedGridTable(rlist); err != nil {
		return err
	}
	//create a list of style ranges
	if rlist, err = types.GetAllRanges(t.CmdList, types.KwStyle); err != nil {
		return err
	}
	return t.applyStyles(rlist)
}

//Render use a types.Renderer to render table contents and write them to io.Writer
func (t *Table) Render(w io.Writer, hr Renderer) error {
	t.Log("***starting rendering table")
	if err := hr.StartTable(t); err != nil {
		return err
	}
	for r, row := range t.grid.rows {
		if err := hr.StartRow(row); err != nil {
			return err
		}

		for c, cell := range row.cells {
			t.Logf("%d,%d:%v\n", r, c, cell.DebugString())
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
		if !types.IsTableCommand(cmd) {
			continue
		}
		cmd.Span().Normalize(t.Contents.RowCount(), t.Contents.MaxFieldCount())
	}
	return nil
}

func (t *Table) applyStyles(rlist []types.Range) error {
	if err := t.grid.ValidateRanges(rlist); err != nil {
		return err
	}
	for _, mr := range rlist {
		for i := mr.TopLeft.Row; i <= mr.BottomRight.Row; i++ {
			for j := mr.TopLeft.Col; j <= mr.BottomRight.Col; j++ {
				t.grid.CellorPanic(i, j).AddStyle(mr.Styles()...)
			}
		}
	}
	return nil
}
