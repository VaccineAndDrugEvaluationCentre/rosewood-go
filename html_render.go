package rosewood

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/drgo/rosewood/types"
	"github.com/drgo/rosewood/utils"
)

const (
	htmlHeader = `
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<meta name="generator" content="Rosewood Carpenter %s" /> 
<meta name="date-generated" content="%s" scheme="YYYY-MM-DD HH:MM:SS">
<link rel="stylesheet" href="%s">
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
	settings *utils.Settings
	tables   []*types.Table
}

func NewHtmlRenderer() *HtmlRenderer {
	return &HtmlRenderer{}
}

func (hr *HtmlRenderer) SetWriter(w io.Writer) error {
	hr.bw = w
	return nil
}

func (hr *HtmlRenderer) SetSettings(settings *utils.Settings) error {
	hr.settings = settings
	return nil
}

func (hr *HtmlRenderer) SetTables(tables []*types.Table) error {
	hr.tables = tables
	return nil
}

func (hr *HtmlRenderer) StartFile() error {
	cssFileName := hr.settings.StyleSheet
	if cssFileName == "" {
		cssFileName = "carpenter.css"
	}
	t := time.Now()
	fmt.Fprintf(hr.bw, htmlHeader, VERSION, t.Format("2006-01-02 15:04:05"), cssFileName)
	return nil
}

func (hr *HtmlRenderer) EndFile() error {
	fmt.Fprintf(hr.bw, htmlFooter)
	return nil
}

func (hr *HtmlRenderer) StartTable(t *types.Table) error {
	fmt.Fprintf(hr.bw, htmlOpenTable)
	if t.Caption != nil {
		fmt.Fprintf(hr.bw, "<caption>")
		for _, line := range t.Caption.Lines {
			fmt.Fprintf(hr.bw, "%s%s\n", line, htmlbreak)
		}
	}
	return nil
}

func (hr *HtmlRenderer) EndTable(t *types.Table) error {
	fmt.Fprintf(hr.bw, htmlCloseTable)
	if t.Footnotes == nil {
		return nil
	}
	fmt.Fprintf(hr.bw, htmlPara)
	for _, line := range t.Footnotes.Lines {
		fmt.Fprintf(hr.bw, "%s%s\n", line, htmlbreak)
	}
	return nil
}

func (hr *HtmlRenderer) StartRow(r *types.Row) error {
	fmt.Fprintf(hr.bw, htmlOpenRow)
	return nil
}

func (hr *HtmlRenderer) EndRow(r *types.Row) error {
	fmt.Fprintf(hr.bw, htmlCloseRow)
	return nil
}

func (hr *HtmlRenderer) OutputCell(c *types.Cell) error {
	if c.State() == types.CsMerged {
		return nil
	}
	fmt.Fprint(hr.bw, "<td")
	if c.RowSpan() > 1 {
		fmt.Fprintf(hr.bw, ` rowspan="%d"`, c.RowSpan())
	}
	if c.ColSpan() > 1 {
		fmt.Fprintf(hr.bw, ` colspan="%d"`, c.ColSpan())
	}
	if hr.settings.TrimCellContents {
		fmt.Fprint(hr.bw, ">", strings.TrimSpace(c.Text()), "</td>")
	} else {
		fmt.Fprint(hr.bw, ">", c.Text(), "</td>")
	}
	return nil
}

// func render(w io.Writer, r types.Renderer, settings *utils.Settings, tables ...*types.Table) error {
// 	//	trace.Printf("%#v", *t)
// 	bw := bufio.NewWriter(w)
// 	r.SetWriter(bw)
// 	r.SetSettings(settings)
// 	r.SetTables(tables)
// 	r.StartFile()
// 	for _, t := range tables {
// 		r.StartTable(t)
// 		for _, row := range t.contents.rows {
// 			r.StartRow(row)
// 			for _, cell := range row.cells {
// 				r.OutputCell(cell)
// 			}
// 			r.EndRow(row)
// 		}
// 		r.EndTable(t)
// 	}
// 	r.EndFile()
// 	bw.Flush()
// 	return nil
// }
