package carpenter

import (
	"bytes"
	"io"
	"strconv"
)

const (
	htmlHeader = `
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<link rel="stylesheet" href="offset-v0_1_0.css">
</head>
<body>
`
	htmlFooter = `
</body>
</html>
`
	htmlOpenTable = `<table class="">
	`
	htmlCloseTable = `</table>
	`
	htmlOpenRow  = `<tr>`
	htmlCloseRow = `</tr>
	`
)

type HtmlRenderer struct {
	w        io.Writer
	settings *Settings
	tables   []*Table
}

func NewHtmlRenderer() *HtmlRenderer {
	return &HtmlRenderer{}
}

func (hr *HtmlRenderer) SetWriter(w io.Writer) error {
	hr.w = w
	return nil
}

func (hr *HtmlRenderer) SetSettings(settings *Settings) error {
	hr.settings = settings
	return nil
}

func (hr *HtmlRenderer) SetTables(tables []*Table) error {
	hr.tables = tables
	return nil
}

func (hr *HtmlRenderer) StartFile() error {
	io.WriteString(hr.w, htmlHeader)
	return nil
}

func (hr *HtmlRenderer) EndFile() error {
	io.WriteString(hr.w, htmlFooter)
	return nil
}

func (hr *HtmlRenderer) StartTable(t *Table) error {
	io.WriteString(hr.w, htmlOpenTable)
	return nil
}

func (hr *HtmlRenderer) EndTable(t *Table) error {
	io.WriteString(hr.w, htmlCloseTable)
	return nil
}

func (hr *HtmlRenderer) StartRow(r *Row) error {
	io.WriteString(hr.w, htmlOpenRow)
	return nil
}

func (hr *HtmlRenderer) EndRow(r *Row) error {
	io.WriteString(hr.w, htmlCloseRow)
	return nil
}

func (hr *HtmlRenderer) OutputCell(c *Cell) error {
	if c.hidden {
		return nil
	}
	var b bytes.Buffer
	b.WriteString("<td")
	if c.rowSpan > 0 {
		b.WriteString(` rowspan="` + strconv.Itoa(int(c.rowSpan)) + `"`)
	}
	if c.colSpan > 0 {
		b.WriteString(` colspan="` + strconv.Itoa(int(c.colSpan)) + `"`)
	}
	b.WriteString(">")
	b.WriteString(c.text)
	b.WriteString("</td>")
	io.WriteString(hr.w, b.String())
	return nil
}

func Render(w io.Writer, r *HtmlRenderer, settings *Settings, tables ...*Table) error {
	//	fmt.Printf("%#v", *t)
	r.SetWriter(w)
	r.SetSettings(settings)
	r.SetTables(tables)
	r.StartFile()
	for _, t := range tables {
		r.StartTable(t)
		for _, row := range t.rows {
			r.StartRow(row)
			for _, cell := range row.cells {
				r.OutputCell(cell)
			}
			r.EndRow(row)
		}
		r.EndTable(t)
	}
	r.EndFile()
	return nil
}
