// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package html

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/drgo/core/files"
	"github.com/drgo/core/md"
	"github.com/drgo/core/ui"
	"github.com/drgo/rosewood"
	"github.com/drgo/rosewood/table"
	"github.com/drgo/rosewood/types"
)

const (
	htmlHeader = `<!DOCTYPE html>
<head>
<meta charset="utf-8">
<meta name="generator" content="Rosewood Carpenter" />
`

	htmlBody = `	
</head>
<body>
`
	htmlFooter = `
</body>
</html>
`

	defaultCSSFileName = "carpenter.css"
)

var defaultCSS []byte

//init, run automatically, registers HTML renderer with Rosewood
func init() {
	config := rosewood.RendererConfig{
		Name:     "html",
		Renderer: makeHTMLRenderer,
	}
	// load default css once per exe run
	exeDir, err := files.GetExeDir()
	if err == nil {
		defaultCSS, _ = ioutil.ReadFile(filepath.Join(exeDir, defaultCSSFileName))
	}
	rosewood.RegisterRenderer(&config)
}

//htmlRenderer implements table.Renderer for HTML output
type htmlRenderer struct {
	bw        io.Writer
	settings  *types.RosewoodSettings
	tables    []*table.Table
	htmlError error  //tracks errors
	css       []byte //holds css text
}

//makeHTMLRenderer factory function according to the renderer registration requirements
func makeHTMLRenderer() (table.Renderer, error) {
	return NewHTMLRenderer()
}

//NewHTMLRenderer create a new htmlRenderer and return it as a Renderer
func NewHTMLRenderer() (table.Renderer, error) {
	return &htmlRenderer{}, nil
}

func (hr *htmlRenderer) SetWriter(w io.Writer) error {
	hr.bw = w
	return nil
}

func (hr *htmlRenderer) SetSettings(settings *types.RosewoodSettings) error {
	hr.settings = settings
	cssFileName := strings.TrimSpace(hr.settings.StyleSheetName)
	if cssFileName == "" { // use default css
		hr.css = defaultCSS
		return nil
	}
	hr.css = []byte(cssFileName)
	var err error
	if hr.settings.DoNotInlineCSS == false {
		if hr.css, err = ioutil.ReadFile(cssFileName); err != nil {
			return fmt.Errorf("failed to load css file %s, %s", cssFileName, err)
		}
	}
	return nil
}

func (hr *htmlRenderer) SetTables(tables []*table.Table) error {
	hr.tables = tables
	return nil
}

func (hr *htmlRenderer) Err() error {
	return hr.htmlError
}

// write does all the writing to the writer and handles errors by stopping any further writing
// and returning the error
func (hr *htmlRenderer) write(s string) error { //TODO:optimize
	if hr.htmlError == nil {
		_, hr.htmlError = hr.bw.Write([]byte(s))
	}
	return hr.htmlError
}

func (hr *htmlRenderer) StartFile() error {
	var b strings.Builder //optimization for golang >= 1.10
	b.Grow(1024 * 100)    //preallocate 100kb to avoid additional allocations
	b.WriteString(htmlHeader)
	b.WriteString(`<meta name="date-generated" content="`)
	b.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	b.WriteString(`" scheme="YYYY-MM-DD HH:MM:SS">` + "\n")
	// FIXME: add settings.HeaderText to support writing anything by the caller to the header
	// ExecutableVersion := fmt.Sprintf("Exe Version %s, Lib Version %s", hr.settings.ExecutableVersion, hr.settings.LibVersion)
	if hr.settings.DoNotInlineCSS {
		b.WriteString(`<link rel="stylesheet" type="text/css" href="` + string(hr.css) + `">`)
	} else {
		b.WriteString("<style>\n")
		b.Write(hr.css)
		b.WriteString("\n</style>\n")
	}
	b.WriteString(htmlBody)
	if hr.settings.Debug >= ui.DebugAll {
		b.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	}
	hr.write(b.String())
	return hr.Err()
}

func (hr *htmlRenderer) EndFile() error {
	return hr.write(htmlFooter)
}

func (hr *htmlRenderer) StartTable(t *table.Table) error {
	hr.write(`<table class="rw-table">`)
	if t.Caption != nil {
		hr.write("<caption>")
		for _, line := range t.Caption.Lines {
			hr.write(hr.renderText(line))
		}
		hr.write("</caption>\n") //added for completeness
	}
	return hr.Err()
}

func (hr *htmlRenderer) EndTable(t *table.Table) error {
	hr.write("</table>\n")
	if t.Footnotes != nil {
		hr.write(`<div class="rw-footnotes">` + "\n")
		for _, line := range t.Footnotes.Lines {
			hr.write(hr.renderText(line) + "<br>\n")
		}
		hr.write("</div>\n")
	}
	return hr.Err()
}

func (hr *htmlRenderer) StartRow(r *table.Row) error {
	return hr.write(`<tr class="rw-row">` + "\n")
}

func (hr *htmlRenderer) EndRow(r *table.Row) error {
	return hr.write("</tr>\n")
}

func (hr *htmlRenderer) OutputCell(c *table.Cell) error {
	if c.Merged() { //skip merged cells
		return nil
	}
	tag := "td"
	if c.Header() {
		tag = "th"
	}
	var b strings.Builder //optimization for golang >= 1.10
	b.Grow(1024)
	b.WriteString("  <" + tag) //open td or th tag
	// write styles
	switch len(c.Styles()) {
	case 0: //donothing
	case 1: //optimization for the common scenario with only one style
		b.WriteString(` class="` + c.Styles()[0] + `"`) //eg class="style1"
	default:
		b.WriteString(` class="` + strings.Join(c.Styles(), " ") + `"`)
	}
	if c.RowSpan() > 1 {
		b.WriteString(fmt.Sprintf(" rowspan=\"%d\"", c.RowSpan())) //eg rowspan="3"
	}
	if c.ColSpan() > 1 {
		b.WriteString(fmt.Sprintf(" colspan=\"%d\"", c.ColSpan())) // eg colspan="2"
	}
	// trim cell contents b/c html ignores white space anyway
	b.WriteString(">" + hr.renderText(strings.TrimSpace(c.Text())) + "</" + tag + ">\n") //eg "> text </td>"
	hr.write(b.String())
	return hr.Err()
}

func (hr *htmlRenderer) renderText(s string) string {
	switch hr.settings.MarkdownRender {
	case "standard", "":
		txt, _ := md.InlinedMdToHTML(s, nil)
		return string(txt)
	case "strict":
		txt, err := md.InlinedMdToHTML(s, nil)
		if err != nil {
			hr.htmlError = fmt.Errorf("error in parsing the following text: %s; error is %s ", strconv.Quote(s), err)
		}
		return string(txt)
	default: //including "disabled"
		return s
	}
}

// // escapeString escapes special characters like "<" to become "&lt;". It
// // modified from stdlib html.EscapeString() to escape <, <=, >, >=,  &, ' and ".
// func escapeString(s string) string {
// 	return htmlEscaper.Replace(s)
// }

// var htmlEscaper = strings.NewReplacer(
// 	`&`, "&amp;",
// 	`'`, "&#39;",
// 	`<=`, "&le;",
// 	`<`, "&lt;",
// 	`>=`, "&ge;",
// 	`=>`, "&ge;",
// 	`>`, "&gt;",
// 	`"`, "&#34;",
// )
