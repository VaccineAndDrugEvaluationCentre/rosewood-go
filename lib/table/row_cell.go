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
	for i := 0; i < colCount; i++ {
		cells[i] = NewCell("", row, i+1)
	}
	return &Row{cells}
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
)

// cellStateLabel: normal, spanned or merged
var cellStateLabel = []string{"n", "s", "hm", "vm"}

//Cell holds information on each table cell
type Cell struct {
	text             string
	row, col         int
	state            CellState
	rowSpan, colSpan int
	styleList        []string
	header           bool //optimization for header cells
}

//NewCell returns a pointer to a new Cell
func NewCell(text string, row, col int) *Cell {
	return &Cell{
		text: text,
		row:  row,
		col:  col}
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
