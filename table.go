package carpenter

import (
	"bytes"
	"fmt"
	"strings"
)

//TODO add coordinate struct for each cell
//TODO add range struct to hold a rectangular range of cells
//add parent link in each object to its parent: cell->range or row->table
//add interface to use in parent property; flds: topleft, bottomright coords, parent

type Table struct {
	rows        []*Row
	maxFldCount int
}

func (t *Table) String() string {
	var b bytes.Buffer
	if t.rows == nil {
		return "invalid: rows is nil"
	}
	for i := 0; i < len(t.rows); i++ {
		b.WriteString(t.rows[i].String())
		b.WriteString("\n")
	}
	return b.String()
}

type Row struct {
	cells []*Cell
}

func (r *Row) String() string {
	var b bytes.Buffer
	if r.cells == nil {
		return "invalid: cells is nil"
	}
	for i := 0; i < len(r.cells); i++ {
		b.WriteString(r.cells[i].String())
		b.WriteString(ColumnSeparator)
	}
	return b.String()
}

type Cell struct {
	text     string
	row, col int
}

func NewCell(text string, row, col int) *Cell {
	return &Cell{
		text: text,
		row:  row,
		col:  col}
}

func (c *Cell) String() string {
	return fmt.Sprintf("r%d c%d: %s", c.row, c.col, c.text)
}

func ParseTableData(s string) (*Table, error) {
	var (
		line, offset          int
		fldCount, maxFldCount int
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
			offset = pos + 1 //offset is now just after the \n
			if fldCount > maxFldCount {
				maxFldCount = fldCount
			}
			rows = append(rows, &Row{cells: cells}) //create a row with currents cells and append to rows
			fldCount = 0                            //reset fldcount
			cells = nil                             //emtpy cell slice
		case '|':
			fldCount++
			cell := NewCell(s[offset:pos], line, fldCount) //text from last offset to just before the separator
			cells = append(cells, cell)
			offset = pos + 1 //offset is now just after the separator
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
