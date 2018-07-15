// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

import (
	"bytes"
	"fmt"
)

type Row struct {
	cells []*Cell
}

func newBlankRow(colCount int) *Row {
	cells := make([]*Cell, colCount)
	for i := 0; i < colCount; i++ {
		cells[i] = &Cell{}
	}
	//	trace.Printf("in newBlankRow %v\n", cells)
	return &Row{cells}
}

//MakeRow for testing
func MakeRow(cells ...*Cell) *Row {
	return &Row{cells}
}

func (r *Row) String() string {
	var b bytes.Buffer
	for _, c := range r.cells {
		b.WriteString(c.String())
		b.WriteString(columnSeparator)
	}
	return b.String()
}

func (t *Row) cellCount() int {
	return len(t.cells)
}

type CellState int

const (
	CsNormal CellState = iota
	CsSpanned
	CsMerged
)

var CellStateLabel = []string{"normal", "spanned", "merged"}

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
	return fmt.Sprintf("r:%d-c:%d=%s-> %s", c.row, c.col, CellStateLabel[c.state], c.text)
}
