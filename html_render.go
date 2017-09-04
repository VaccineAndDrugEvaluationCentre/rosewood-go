package rosewood

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	htmlHeader = `
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<meta name="generator" content="Rosewood Carpenter %s" /> 
<link rel="stylesheet" href="carpenter.css">
</head>
<body>
`
	htmlFooter = `
</body>
</html>
`
	htmlOpenTable = `<table>
	`
	htmlCloseTable = `</table>
	`
	htmlOpenRow  = `<tr>`
	htmlCloseRow = `</tr>
	`
	htmlPara  = "<p>"
	htmlbreak = "<br>"
)

type HtmlRenderer struct {
	bw       io.Writer
	settings *Settings
	tables   []*table
}

func NewHtmlRenderer() *HtmlRenderer {
	return &HtmlRenderer{}
}

func (hr *HtmlRenderer) SetWriter(w io.Writer) error {
	hr.bw = w
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
	fmt.Fprintf(hr.bw, htmlHeader, VERSION)
	return nil
}

func (hr *HtmlRenderer) EndFile() error {
	fmt.Fprintf(hr.bw, htmlFooter)
	return nil
}

func (hr *HtmlRenderer) StartTable(t *table) error {
	fmt.Fprintf(hr.bw, htmlOpenTable)
	if t.caption != nil {
		fmt.Fprintf(hr.bw, "<caption>")
		for _, line := range t.caption.lines {
			fmt.Fprintf(hr.bw, "%s%s\n", line, htmlbreak)
		}
	}
	return nil
}

func (hr *HtmlRenderer) EndTable(t *table) error {
	fmt.Fprintf(hr.bw, htmlCloseTable)
	if t.footnotes == nil {
		return nil
	}
	fmt.Fprintf(hr.bw, htmlPara)
	for _, line := range t.footnotes.lines {
		fmt.Fprintf(hr.bw, "%s%s\n", line, htmlbreak)
	}
	return nil
}

func (hr *HtmlRenderer) StartRow(r *Row) error {
	fmt.Fprintf(hr.bw, htmlOpenRow)
	return nil
}

func (hr *HtmlRenderer) EndRow(r *Row) error {
	fmt.Fprintf(hr.bw, htmlCloseRow)
	return nil
}

func (hr *HtmlRenderer) OutputCell(c *Cell) error {
	if c.state == csMerged {
		return nil
	}
	fmt.Fprint(hr.bw, "<td")
	if c.rowSpan > 1 {
		fmt.Fprintf(hr.bw, ` rowspan="%d"`, c.rowSpan)
	}
	if c.colSpan > 1 {
		fmt.Fprintf(hr.bw, ` colspan="%d"`, c.colSpan)
	}
	if hr.settings.TrimCellContents {
		fmt.Fprint(hr.bw, ">", strings.TrimSpace(c.text), "</td>")
	} else {
		fmt.Fprint(hr.bw, ">", c.text, "</td>")
	}
	return nil
}

func render(w io.Writer, r *HtmlRenderer, settings *Settings, tables ...*table) error {
	//	trace.Printf("%#v", *t)
	bw := bufio.NewWriter(w)
	r.SetWriter(bw)
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
	bw.Flush()
	return nil
}
