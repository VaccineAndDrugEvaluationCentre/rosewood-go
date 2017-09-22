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
	// htmlOpenTable = `<table>
	// `
	// htmlCloseTable = `</table>
	// `
	// htmlOpenRow  = `<tr>`
	// htmlCloseRow = `</tr>
	// `
	htmlPara     = "<p>"
	htmlbreak    = "<br>"
	htmlOpenDiv  = "<div>"
	htmlCloseDiv = "</div>"
)

//htmlRenderer implements types.Renderer for HTML output
type htmlRenderer struct {
	bw        io.Writer
	settings  *utils.Settings
	tables    []*types.Table
	htmlError error //tracks errors
}

//NewHTMLRenderer create a new htmlRenderer and return a Renderer
func NewHTMLRenderer() types.Renderer {
	return &htmlRenderer{}
}

func (hr *htmlRenderer) SetWriter(w io.Writer) error {
	hr.bw = w
	return nil
}

func (hr *htmlRenderer) SetSettings(settings *utils.Settings) error {
	hr.settings = settings
	return nil
}

func (hr *htmlRenderer) SetTables(tables []*types.Table) error {
	hr.tables = tables
	return nil
}

func (hr *htmlRenderer) Err() error {
	return hr.htmlError
}

//write does all the writing to the writer and handles errors but stopping any further writing
func (hr *htmlRenderer) write(format string, a ...interface{}) error {
	if hr.htmlError == nil {
		_, hr.htmlError = fmt.Fprintf(hr.bw, format, a...)
	}
	return hr.htmlError
}

func (hr *htmlRenderer) StartFile() error {
	cssFileName := hr.settings.StyleSheetName
	if cssFileName == "" {
		cssFileName = "carpenter.css"
	}
	ExecutableVersion := fmt.Sprintf("Exe Version %s, Lib Version %s", hr.settings.ExecutableVersion, Version)
	t := time.Now()
	return hr.write(htmlHeader, ExecutableVersion, t.Format("2006-01-02 15:04:05"), cssFileName)
}

func (hr *htmlRenderer) EndFile() error {
	return hr.write(htmlFooter)
}

func (hr *htmlRenderer) StartTable(t *types.Table) error {
	hr.write("<table>")
	if t.Caption != nil {
		hr.write("<caption>")
		for _, line := range t.Caption.Lines {
			hr.write("%s%s\n", line, htmlbreak)
		}
		hr.write("</caption>") //added for completeness
	}
	return hr.Err()
}

func (hr *htmlRenderer) EndTable(t *types.Table) error {
	hr.write("</table>")
	if t.Footnotes != nil {
		hr.write(`<div class="footnotes">`)
		for _, line := range t.Footnotes.Lines {
			hr.write("%s%s\n", line, htmlbreak)
		}
		hr.write("</div>\n")
	}
	return hr.Err()
}

func (hr *htmlRenderer) StartRow(r *types.Row) error {
	return hr.write(`<tr>`)
}

func (hr *htmlRenderer) EndRow(r *types.Row) error {
	return hr.write("</tr>")
}
func (hr *htmlRenderer) OutputCell(c *types.Cell) error {
	if c.State() == types.CsMerged {
		return nil
	}
	hr.write("<td")
	if len(c.Styles()) > 0 {
		hr.write(` class="%s"`, strings.Join(c.Styles(), " "))
	}
	if c.RowSpan() > 1 {
		hr.write(` rowspan="%d"`, c.RowSpan())
	}
	if c.ColSpan() > 1 {
		hr.write(` colspan="%d"`, c.ColSpan())
	}
	txt := escapeString(c.Text())
	if hr.settings.TrimCellContents {
		txt = strings.TrimSpace(txt)
	}
	hr.write("%s%s%s", ">", txt, "</td>")
	return hr.Err()
}

// escapeString escapes special characters like "<" to become "&lt;". It
// modified from html.EscapeString to escape <, <=, >, >=,  &, ' and ".
func escapeString(s string) string {
	return htmlEscaper.Replace(s)
}

var htmlEscaper = strings.NewReplacer(
	`&`, "&amp;",
	`'`, "&#39;",
	`<=`, "&le;",
	`<=`, "&le;",
	`<`, "&lt;",
	`>=`, "&ge;",
	`=>`, "&ge;",
	`>`, "&gt;",
	`"`, "&#34;",
)

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
