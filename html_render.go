package rosewood

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
	tables   []*table
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

func (hr *HtmlRenderer) SetTables(tables []*table) error {
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

func (hr *HtmlRenderer) StartTable(t *table) error {
	if t.caption != nil && t.caption.LineCount() > 0 {
		io.WriteString(hr.w, t.caption.String()) //strings.Join(t.caption.lines, "")
	}
	io.WriteString(hr.w, htmlOpenTable)
	return nil
}

func (hr *HtmlRenderer) EndTable(t *table) error {
	io.WriteString(hr.w, htmlCloseTable)
	if t.footnotes != nil && t.footnotes.LineCount() > 0 {
		io.WriteString(hr.w, t.footnotes.String()) //strings.Join(t.caption.lines, "")
	}
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
	w := bufio.NewWriter(&b)
	fmt.Fprint(w, "<td")
	if c.rowSpan > 0 {
		fmt.Fprintf(w, ` rowspan="%d"`, c.rowSpan)
	}
	if c.colSpan > 0 {
		fmt.Fprintf(w, ` colspan="%d"`, c.colSpan)
	}
	fmt.Fprint(w, ">", c.text, "</td>")
	w.Flush()
	io.WriteString(hr.w, b.String())
	return nil
}

func render(w io.Writer, r *HtmlRenderer, settings *Settings, tables ...*table) error {
	//	fmt.Printf("%#v", *t)
	r.SetWriter(w)
	r.SetSettings(settings)
	r.SetTables(tables)
	r.StartFile()
	for _, t := range tables {
		r.StartTable(t)
		for _, row := range t.contents.rows {
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
