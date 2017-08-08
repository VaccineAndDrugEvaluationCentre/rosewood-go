package carpenter

import (
	"bytes"
	"fmt"
	"strings"
)

//add parent link in each object to its parent: cell->range or row->table
//add interface to use in parent property; flds: topleft, bottomright coords, parent

type Table struct {
	rows        []*Row
	maxFldCount RwInt
}

func (t *Table) String() string {
	var b bytes.Buffer
	// if t.rows == nil {
	// 	return "invalid: rows is nil"
	// }
	for _, r := range t.rows {
		b.WriteString(r.String())
		b.WriteString("\n")
	}
	return b.String()
}

func (t *Table) ValidCoordinate(row, col RwInt) bool {
	if row < 1 || row > RwInt(len(t.rows)) {
		return false
	}
	if col < 1 || col > RwInt(len(t.rows[row-1].cells)) {
		return false
	}
	return true
}

func (t *Table) Cell(row, col RwInt) *Cell {
	if !t.ValidCoordinate(row, col) {
		return nil
	}
	return t.cell(row, col)
}

func (t *Table) cell(row, col RwInt) *Cell {
	return t.rows[row-1].cells[col-1]
}

func (t *Table) Merge(ra rwRange) error {
	if err := ra.Validate(); err != nil {
		return err
	}
	for r := ra.TopLeft.Row; r <= ra.BottomRight.Row; r++ {
		for c := ra.TopLeft.Col; c <= ra.BottomRight.Col; c++ {
			t.cell(r, c).hidden = true
		}
	}
	topleft := t.cell(ra.TopLeft.Row, ra.TopLeft.Col)
	topleft.hidden = false
	topleft.rowSpan = ra.BottomRight.Row - ra.TopLeft.Row
	topleft.colSpan = ra.BottomRight.Col - ra.TopLeft.Col
	return nil
}

type Row struct {
	cells []*Cell
}

func (r *Row) String() string {
	var b bytes.Buffer
	for _, c := range r.cells {
		b.WriteString(c.String())
		b.WriteString(ColumnSeparator)
	}
	return b.String()
}

type Cell struct {
	text             string
	row, col         RwInt
	hidden           bool
	rowSpan, colSpan RwInt
}

func NewCell(text string, row, col RwInt) *Cell {
	return &Cell{
		text: text,
		row:  row,
		col:  col}
}

// func NewHiddenCell(text string, row, col int) *Cell {
// 	return &Cell{
// 		text:   text,
// 		row:    row,
// 		col:    col,
// 		hidden: true}
// }

func (c *Cell) String() string {
	return fmt.Sprintf("r%d c%d: %s", c.row, c.col, c.text)
}

func ParseTableData(s string) (*Table, error) {
	var (
		line, offset          RwInt
		fldCount, maxFldCount RwInt
		cells                 []*Cell
		rows                  []*Row
	)
	if strings.TrimSpace(s) == "" {
		return nil, fmt.Errorf("empty table")
	}
	line = 1
	for pos := 0; pos < len(s); pos++ {
		switch s[pos] {
		case '\r': //carriage return followed by linefeed as EOL sequence (Windows)
			//todo: \r must be preceded by |
		case '\n': //linefeed for Linux and MacOs as EOL marker
			//todo: \n must be preceded by | or \r
			line++
			offset = RwInt(pos + 1) //offset is now just after the \n
			if fldCount > maxFldCount {
				maxFldCount = fldCount
			}
			rows = append(rows, &Row{cells: cells}) //create a row with currents cells and append to rows
			fldCount = 0                            //reset fldcount
			cells = nil                             //emtpy the cell slice
		case '|':
			fldCount++
			cell := NewCell(s[offset:pos], line, fldCount) //text from last offset to just before the separator
			cells = append(cells, cell)
			offset = RwInt(pos + 1) //offset is now just after the separator
		}
	}
	if maxFldCount == 0 || len(rows) == 0 { //no fields or rows found
		return nil, fmt.Errorf("invalid data table")
	}
	return &Table{rows: rows,
		maxFldCount: maxFldCount}, nil

}

//for testing only; it ignores errors
func monadicParseTableData(s string) *Table {
	t, _ := ParseTableData(s)
	return t
}

func genTestTableData(rows []*Row) *Table {
	return &Table{
		rows:        rows,
		maxFldCount: 5,
	}
}
