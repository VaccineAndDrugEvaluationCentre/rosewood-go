// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package table

import (
	"bytes"
	"fmt"
	"strings"
)

//Row holds info on a Rosewood table row
type Row struct {
	cells []*Cell
}

func newBlankRow(row int, colCount int) *Row {
	cells := make([]*Cell, colCount)
	thisRow := &Row{cells}
	for i := 0; i < colCount; i++ {
		cells[i] = NewCell("", row, i+1, thisRow)
	}
	return thisRow
}

//MakeRow for testing
func makeRow(cells ...*Cell) *Row {
	return &Row{cells}
}

func (r *Row) String() string {
	var b bytes.Buffer
	for _, c := range r.cells {
		b.WriteString(c.String())
		b.WriteString("|") //just for string formatting
	}
	return b.String()
}

func (r *Row) cellCount() int {
	return len(r.cells)
}

// Number returns row number
func (r *Row) Number() int {
	if len(r.cells) > 0 {
		return r.cells[0].row
	}
	return -1
}

// LastCell returns pointer to last cell in the row
func (r *Row) LastCell() *Cell {
	if r.cellCount() > 0 {
		return r.cells[r.cellCount()-1]
	}
	return nil
}

// LastVisibleCell returns pointer to last cell in the row that is not merged or spanned
func (r *Row) LastVisibleCell() *Cell {
	for i := len(r.cells) - 1; i >= 0; i-- {
		if !r.cells[i].Merged() {
			return r.cells[i]
		}
	}
	return nil
}

// Header returns true if the first cell in a row has a header style
func (r *Row) Header() bool {
	return r.cellCount() > 1 && r.cells[0].Header()
}

//CellState describes whether a cell is merged, spanned or otherwise
type CellState int

const (
	//CsNormal regular cell
	CsNormal CellState = iota
	//CsSpanned cell included in a span
	CsSpanned
	//CsHMerged cell is part of a col merge eg merge col 1:3
	CsHMerged
	//CsVMerged cell is part of a row merge eg merge row 1:3
	CsVMerged
	//CsVHMerged cell is part of a row and col merge eg merge row 1:2 col 2:3
	CsVHMerged
)

// cellStateLabel: normal, spanned or merged
var cellStateLabel = []string{"n", "s", "h", "v", "b"}

//Cell holds information on each table cell
type Cell struct {
	text             string
	row, col         int
	state            CellState
	rowSpan, colSpan int
	styleList        []string
	header           bool //optimization for header cells
	parentRow        *Row
}

//NewCell returns a pointer to a new Cell
func NewCell(text string, row, col int, parentRow *Row) *Cell {
	return &Cell{
		text:      text,
		row:       row,
		col:       col,
		parentRow: parentRow}
}

//MakeCell creates a new cell for testing
func MakeCell(text string, row, col int, state CellState, rowSpan, colSpan int) *Cell {
	return &Cell{
		text: text, row: row, col: col,
		state: state, rowSpan: rowSpan, colSpan: colSpan}
}

func (c *Cell) clone(src *Cell) *Cell {
	*c = *src
	return c
}
func (c *Cell) State() CellState {
	return c.state
}
func (c *Cell) Text() string {
	return c.text
}
func (c *Cell) RowSpan() int {
	return c.rowSpan
}
func (c *Cell) ColSpan() int {
	return c.colSpan
}

func (c *Cell) Header() bool {
	return c.header
}

func (c *Cell) Styles() []string {
	return c.styleList
}

func (c *Cell) Merged() bool {
	return c.state >= CsHMerged
}

// Row returns pointer to the row holding this cell
func (c *Cell) Row() *Row {
	return c.parentRow
}

// LastCell true if this cells is last cell in the row
func (c *Cell) LastCell() bool {
	return c.parentRow.LastCell() == c
}

// LastVisibleCell true if this cells is last cell in the row that is not merged or spanned
func (c *Cell) LastVisibleCell() bool {
	return c.parentRow.LastVisibleCell() == c
}

//AddStyle adds one or more style names if they do not already exist in the list
//sufficiently efficient for short lists and avoids allocating a map
func (c *Cell) AddStyle(styles ...string) error {
outer:
	for _, s := range styles {
		for _, ss := range c.styleList { //skip s if it already exists in the list
			if ss == s {
				continue outer
			}
		}
		if s == "header" {
			c.header = true //optimization for header cells
		} else {
			c.styleList = append(c.styleList, s)
		}
	}
	return nil
}

func (c *Cell) String() string {
	return fmt.Sprintf("r%d c%d: %s", c.row, c.col, c.text)
}

func (c *Cell) DebugString() string {
	return fmt.Sprintf("%d,%d (%s): %s", c.row, c.col, cellStateLabel[c.state], strings.TrimSpace(c.text))
}
