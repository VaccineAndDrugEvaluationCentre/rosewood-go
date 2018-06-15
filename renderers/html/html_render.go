// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package html

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	rosewood "github.com/drgo/rosewood/lib"
	"github.com/drgo/rosewood/lib/setter"
	"github.com/drgo/rosewood/lib/types"
)

const (
	htmlHeader = `<!DOCTYPE html>
<head>
<meta charset="utf-8">
<meta name="generator" content="Rosewood Carpenter new" />
`

	htmlBody = `	
</head>
<body>
`
	htmlFooter = `
</body>
</html>
`
)

//init registers HTML renderer with Rosewood
func init() {
	config := rosewood.RendererConfig{
		Name:     "html",
		Renderer: makeHTMLRenderer,
	}
	rosewood.Register(&config)
}

//htmlRenderer implements types.Renderer for HTML output
type htmlRenderer struct {
	bw        io.Writer
	settings  *setter.Settings
	tables    []*types.Table
	htmlError error //tracks errors
}

//makeHTMLRenderer factory functions according to the renderer registration requirements
func makeHTMLRenderer() (types.Renderer, error) {
	return NewHTMLRenderer(), nil
}

//NewHTMLRenderer create a new htmlRenderer and return it as a Renderer
func NewHTMLRenderer() types.Renderer {
	return &htmlRenderer{}
}

func (hr *htmlRenderer) SetWriter(w io.Writer) error {
	hr.bw = w
	return nil
}

func (hr *htmlRenderer) SetSettings(settings *setter.Settings) error {
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

//write does all the writing to the writer and handles errors by stopping any further writing
func (hr *htmlRenderer) write(format string, a ...interface{}) error { //TODO:optimize
	if hr.htmlError == nil {
		_, hr.htmlError = fmt.Fprintf(hr.bw, format, a...)
	}
	return hr.htmlError
}

func (hr *htmlRenderer) StartFile() error {
	var err error
	var b strings.Builder //optimization for golang >= 1.10
	b.Grow(1024 * 100)    //preallocate 100kb to avoid additional allocations
	b.WriteString(htmlHeader)
	b.WriteString(`<meta name="date-generated" content="`)
	b.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	b.WriteString(`" scheme="YYYY-MM-DD HH:MM:SS">
		`)
	//FIXME:
	//	ExecutableVersion := fmt.Sprintf("Exe Version %s, Lib Version %s", hr.settings.ExecutableVersion, hr.settings.LibVersion)
	cssFileName := hr.settings.StyleSheetName
	if cssFileName == "" {
		cssFileName = "carpenter.css"
	}
	css := []byte(cssFileName)
	if hr.settings.DoNotInlineCSS == false {
		if css, err = ioutil.ReadFile(cssFileName); err != nil { //optimize using io.copy
			return err
		}
	}
	b.WriteString("<style>\n")
	b.Write(css)
	b.WriteString("\n</style>\n")
	b.WriteString(htmlBody)
	if hr.settings.Debug >= setter.DebugAll {
		b.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	}
	hr.write(b.String())
	return hr.Err()
}

func (hr *htmlRenderer) EndFile() error {
	return hr.write(htmlFooter)
}

func (hr *htmlRenderer) StartTable(t *types.Table) error {
	hr.write(`<table class="rw-table">
		`)
	if t.Caption != nil {
		hr.write("<caption>")
		for _, line := range t.Caption.Lines {
			hr.write("%s%s\n", line, "<br>")
		}
		hr.write("</caption>\n") //added for completeness
	}
	return hr.Err()
}

func (hr *htmlRenderer) EndTable(t *types.Table) error {
	hr.write("</table>\n")
	if t.Footnotes != nil {
		hr.write(`<div class="rw-footnotes">%s`, "\n")
		for _, line := range t.Footnotes.Lines {
			hr.write("%s%s\n", line, "<br>")
		}
		hr.write("</div>\n")
	}
	return hr.Err()
}

func (hr *htmlRenderer) StartRow(r *types.Row) error {
	return hr.write(`<tr class="rw-row">%s`, "\n")
}

func (hr *htmlRenderer) EndRow(r *types.Row) error {
	return hr.write("</tr>\n")
}

func (hr *htmlRenderer) OutputCell(c *types.Cell) error {
	//TODO: not working; merged cells are still printed
	//fmt.Printf("%s\n", c.DebugString()) //DEBUG
	if c.State() == types.CsMerged {
		return nil
	}
	hr.write("<td")
	//EXPERIMENTAL: add base cell style rw-cell
	hr.write(` class="rw-cell %s"`, strings.Join(c.Styles(), " "))

	if c.RowSpan() > 1 {
		hr.write(` rowspan="%d"`, c.RowSpan())
	}
	if c.ColSpan() > 1 {
		hr.write(` colspan="%d"`, c.ColSpan())
	}
	txt := escapeString(c.Text())
	//	if hr.settings.TrimCellContents {
	txt = strings.TrimSpace(txt)
	//	}
	hr.write(">%s</td>\n", txt)
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
