// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

const columnSeparator = "|"

var OsEOL string

func init() {
	OsEOL = "\n"
	if runtime.GOOS == "windows" {
		OsEOL = "\r\n"
	}
}

type TableContents struct {
	rows        []*Row
	maxFldCount RwInt
}

func (t *TableContents) String() string {
	var b bytes.Buffer
	for _, r := range t.rows {
		b.WriteString(r.String())
		b.WriteString(OsEOL)
	}
	return b.String()
}

//MaxFieldCount returns the maximum number of cells in a row
func (t *TableContents) MaxFieldCount() RwInt {
	return t.maxFldCount
}

type cellFunc func(c *Cell) error

func (t *TableContents) forEachCell(f cellFunc) error {
	for _, r := range t.rows {
		for _, c := range r.cells {
			if err := f(c); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *TableContents) isValidCoordinate(row, col RwInt) bool {
	if row < 1 || row > t.RowCount() {
		return false
	}
	if col < 1 || col > t.Row(row).cellCount() {
		return false
	}
	return true
}

//Row return the ith row (warning 1 based not zero based)
func (t *TableContents) Row(i RwInt) *Row {
	return t.rows[i-1]
}

//RowCount returns the number of rows in a table
func (t *TableContents) RowCount() RwInt {
	return RwInt(len(t.rows))
}

//CellorPanic returns the cell at row, col coordinates (warning 1 based not zero based)
//panics if the coordinates are not valid
func (t *TableContents) CellorPanic(row, col RwInt) *Cell {
	if !t.isValidCoordinate(row, col) {
		panic(fmt.Sprintf("invalid cell coordinates, row=%d, col=%d", row, col))
	}
	return t.cell(row, col)
}

//cell returns the cell at row, col coordinates; panics if coordinates are not valid
func (t *TableContents) cell(row, col RwInt) *Cell {
	return t.rows[row-1].cells[col-1]
}

//ValidateRanges takes a list of ranges and return an error if they are not valid coordinates within this table
func (t *TableContents) ValidateRanges(rList []Range) error {
	for _, r := range rList {
		if err := t.validateRange(r); err != nil {
			return err
		}
	}
	return nil
}

func (t *TableContents) validateRange(ra Range) error {
	//check that the range itself makes sense
	if err := ra.validate(); err != nil {
		return err
	}
	//check each coordinate against the table coordinates
	for r := ra.TopLeft.Row; r <= ra.BottomRight.Row; r++ {
		for c := ra.TopLeft.Col; c <= ra.BottomRight.Col; c++ {
			if !t.isValidCoordinate(r, c) {
				return fmt.Errorf("invalid coordinates: [%d,%d] ", r, c)
			}
		}
	}
	return nil
}

//NewTableContents parses a Rosewood table contents
func NewTableContents(text string) (*TableContents, error) {
	var (
		line, offset          RwInt
		fldCount, maxFldCount RwInt
		cells                 []*Cell
		rows                  []*Row
	)
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("empty table")
	}
	//add eol at the end if none
	if text[len(text)-1:] != "\n" {
		text = text + "\n"
	}
	line = 1
	for pos := 0; pos < len(text); pos++ {
		switch text[pos] {
		case '\r': //carriage return followed by linefeed as EOL sequence (Windows)

		case '\n': //linefeed for Linux and MacOs as EOL marker
			//TODO: add option to prevent cells with no field separator
			// if fldCount == 0 {
			// 	return nil, fmt.Errorf("row #%d has no cells", line)
			// }
			if fldCount > maxFldCount {
				maxFldCount = fldCount
			}
			rows = append(rows, &Row{cells: cells}) //create a row with currents cells and append to rows
			line++
			offset = RwInt(pos + 1) //offset is now just after the \n
			fldCount = 0            //reset fldcount
			cells = nil             //emtpy the cell slice
		case '|':
			fldCount++
			cell := NewCell(text[offset:pos], line, fldCount) //text from last offset to just before the separator
			cells = append(cells, cell)
			offset = RwInt(pos + 1) //offset is now just after the separator
		}
	}
	//TODO: fix situation where table has one field and no column separators
	if maxFldCount == 0 {
		return nil, fmt.Errorf("invalid data table: field count is 0")
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("invalid data table, row count is 0")
	}
	return &TableContents{rows: rows,
		maxFldCount: maxFldCount}, nil

}

//NewBlankTableContents creates an empty TableContents with Rowcount X colCount cells
func NewBlankTableContents(rowCount, colCount RwInt) *TableContents {
	rows := make([]*Row, rowCount)
	for i := RwInt(0); i < rowCount; i++ {
		rows[i] = newBlankRow(colCount)
	}
	return &TableContents{rows: rows,
		maxFldCount: colCount}
}

//MakeTableContents creates a TableContents from an array of Rows; use for testing only
func MakeTableContents(rows []*Row, maxFldCount RwInt) *Table {
	return &Table{Contents: &TableContents{
		rows:        rows,
		maxFldCount: maxFldCount,
	}}
}

func (t *TableContents) DebugString() string {
	var b bytes.Buffer
	for rn, r := range t.rows {
		for cn, _ := range r.cells {
			b.WriteString(fmt.Sprintf("%d,%d |", rn, cn))
		}
		b.WriteString(OsEOL)
	}
	return b.String()
}
