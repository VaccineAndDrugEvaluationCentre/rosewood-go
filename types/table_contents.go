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

//row return the ith row (warning 1 based not zero based)
func (t *TableContents) Row(i RwInt) *Row {
	return t.rows[i-1]
}

func (t *TableContents) RowCount() RwInt {
	return RwInt(len(t.rows))
}

func (t *TableContents) Cell(row, col RwInt) *Cell {
	if !t.isValidCoordinate(row, col) {
		return nil
	}
	return t.cell(row, col)
}

func (t *TableContents) cell(row, col RwInt) *Cell {
	return t.rows[row-1].cells[col-1]
}

func (t *TableContents) validateRange(ra Range) (Range, error) {
	if err := ra.validate(); err != nil {
		return ra, err
	}
	for r := ra.TopLeft.Row; r <= ra.BottomRight.Row; r++ {
		for c := ra.TopLeft.Col; c <= ra.BottomRight.Col; c++ {
			//trace.Printf("r%d c%d \n", r, c)
			if !t.isValidCoordinate(r, c) {

			}
		}
	}
	// normalize := func(value, Default RwInt) RwInt {
	// 	if value < MinRwInt || value == MissingRwInt {
	// 		return Default
	// 	}
	// 	return value
	// }
	// ra.TopLeft.Row = normalize(ra.TopLeft.Row, MinRwInt)
	// ra.TopLeft.Col = normalize(ra.TopLeft.Col, MinRwInt)

	// ra.BottomRight.Row = normalize(ra.BottomRight.Row, t.maxFldCount)
	// ra.BottomRight.Col = normalize(ra.BottomRight.Col, t.maxFldCount)

	return ra, nil
}

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
			//TODO: \r must be preceded by |
		case '\n': //linefeed for Linux and MacOs as EOL marker
			//TODO: \n must be preceded by | or \r
			if fldCount == 0 {
				return nil, fmt.Errorf("row #%d has no cells", line)
			}
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
	if maxFldCount == 0 {
		return nil, fmt.Errorf("invalid data table: field count is 0")
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("invalid data table, row count is 0")
	}
	return &TableContents{rows: rows,
		maxFldCount: maxFldCount}, nil

}

func NewBlankTableContents(RowCount, colCount RwInt) *TableContents {
	rows := make([]*Row, RowCount)
	for i := RwInt(0); i < RowCount; i++ {
		rows[i] = newBlankRow(colCount)
	}
	//	trace.Printf("in newBlankTableContents %v\n", rows)
	return &TableContents{rows: rows,
		maxFldCount: colCount}
}

// func (t *TableContents) merge(ra Range) error {
// 	var err error
// 	if ra, err = t.validateRange(ra); err != nil {
// 		return fmt.Errorf("merge failed: %s", err)
// 	}
// 	trace.Printf("ra=%s  %s tl=%s  br=%s\n", ra.testString(), ra, ra.TopLeft, ra.BottomRight)
// 	for r := ra.TopLeft.Row; r <= ra.BottomRight.Row; r++ {
// 		for c := ra.TopLeft.Col; c <= ra.BottomRight.Col; c++ {
// 			trace.Printf("r%d c%d \n", r, c)
// 			if t.Cell(r, c) != nil {
// 				t.Cell(r, c).hidden = true
// 			}
// 		}
// 	}
// 	topleft := t.Cell(ra.TopLeft.Row, ra.TopLeft.Col)
// 	if topleft == nil {
// 		trace.Printf("topleft is nil, r%d c%d \n", ra.TopLeft.Row, ra.TopLeft.Col)
// 	} else {
// 		topleft.hidden = false
// 		topleft.rowSpan = ra.BottomRight.Row - ra.TopLeft.Row
// 		topleft.colSpan = ra.BottomRight.Col - ra.TopLeft.Col
// 	}
// 	return nil
// }
