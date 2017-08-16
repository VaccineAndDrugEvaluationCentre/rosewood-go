package rosewood

import (
	"bytes"
	"fmt"
	"strings"
)

//add parent link in each object to its parent: cell->range or row->table
//add interface to use in parent property; flds: topleft, bottomright coords, parent

type tableContents struct {
	rows        []*Row
	maxFldCount RwInt
}

func (t *tableContents) String() string {
	var b bytes.Buffer
	for _, r := range t.rows {
		b.WriteString(r.String())
		b.WriteString(OSEOL)
	}
	return b.String()
}

func (t *tableContents) isValidCoordinate(row, col RwInt) bool {
	if row < 1 || row > RwInt(len(t.rows)) {
		return false
	}
	if col < 1 || col > RwInt(len(t.rows[row-1].cells)) {
		return false
	}
	return true
}

func (t *tableContents) Cell(row, col RwInt) *Cell {
	if !t.isValidCoordinate(row, col) {
		return nil
	}
	return t.cell(row, col)
}

func (t *tableContents) cell(row, col RwInt) *Cell {
	return t.rows[row-1].cells[col-1]
}

func (t *tableContents) validateRange(ra Range) (Range, error) {
	if err := ra.validate(); err != nil {
		return ra, err
	}
	for r := ra.TopLeft.Row; r <= ra.BottomRight.Row; r++ {
		for c := ra.TopLeft.Col; c <= ra.BottomRight.Col; c++ {
			fmt.Printf("r%d c%d \n", r, c)
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

func (t *tableContents) merge(ra Range) error {
	var err error
	if ra, err = t.validateRange(ra); err != nil {
		return fmt.Errorf("merge failed: %s", err)
	}
	fmt.Printf("ra=%s  %s tl=%s  br=%s\n", ra.testString(), ra, ra.TopLeft, ra.BottomRight)
	for r := ra.TopLeft.Row; r <= ra.BottomRight.Row; r++ {
		for c := ra.TopLeft.Col; c <= ra.BottomRight.Col; c++ {
			fmt.Printf("r%d c%d \n", r, c)
			if t.Cell(r, c) != nil {
				t.Cell(r, c).hidden = true
			}
		}
	}
	topleft := t.Cell(ra.TopLeft.Row, ra.TopLeft.Col)
	if topleft == nil {
		fmt.Printf("topleft is nil, r%d c%d \n", ra.TopLeft.Row, ra.TopLeft.Col)
	} else {
		topleft.hidden = false
		topleft.rowSpan = ra.BottomRight.Row - ra.TopLeft.Row
		topleft.colSpan = ra.BottomRight.Col - ra.TopLeft.Col
	}
	return nil
}

type Row struct {
	cells []*Cell
}

func (r *Row) String() string {
	var b bytes.Buffer
	for _, c := range r.cells {
		b.WriteString(c.String())
		b.WriteString(columnSeparator)
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

func (c *Cell) String() string {
	return fmt.Sprintf("r%d c%d: %s", c.row, c.col, c.text)
}

func NewTableContents(text string) (*tableContents, error) {
	var (
		line, offset          RwInt
		fldCount, maxFldCount RwInt
		cells                 []*Cell
		rows                  []*Row
	)
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("empty table")
	}
	line = 1
	for pos := 0; pos < len(text); pos++ {
		switch text[pos] {
		case '\r': //carriage return followed by linefeed as EOL sequence (Windows)
			//todo: \r must be preceded by |
		case '\n': //linefeed for Linux and MacOs as EOL marker
			//todo: \n must be preceded by | or \r
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
	if maxFldCount == 0 || len(rows) == 0 { //no fields or rows found
		return nil, fmt.Errorf("invalid data table")
	}
	return &tableContents{rows: rows,
		maxFldCount: maxFldCount}, nil

}
