// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package html

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	rosewood "github.com/drgo/rosewood/lib"
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
	settings  *types.RosewoodSettings
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

func (hr *htmlRenderer) SetSettings(settings *types.RosewoodSettings) error {
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
func (hr *htmlRenderer) write(s string) error { //TODO:optimize
	if hr.htmlError == nil {
		_, hr.htmlError = hr.bw.Write([]byte(s))
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
	if hr.settings.Debug >= types.DebugAll {
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
			hr.write(hr.renderText(line))
		}
		hr.write("</caption>\n") //added for completeness
	}
	return hr.Err()
}

func (hr *htmlRenderer) EndTable(t *types.Table) error {
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

func (hr *htmlRenderer) StartRow(r *types.Row) error {
	return hr.write(`<tr class="rw-row">` + "\n")
}

func (hr *htmlRenderer) EndRow(r *types.Row) error {
	return hr.write("</tr>\n")
}

func (hr *htmlRenderer) renderText(s string) string {
	txt, err := InlinedMdToHTML(s, nil)
	if err != nil {
		hr.htmlError = fmt.Errorf("error in parsing the following text: %s; error is %s ", strconv.Quote(s), err)
	}
	return string(txt)
}

func (hr *htmlRenderer) OutputCell(c *types.Cell) error {
	//fmt.Printf("%s\n", c.DebugString()) //DEBUG
	if c.State() == types.CsMerged {
		return nil
	}
	tag := "td"
	if c.Header() {
		tag = "th"
	}
	var b strings.Builder //optimization for golang >= 1.10
	b.Grow(1024)
	b.WriteString("<" + tag) //open td or th tag
	switch len(c.Styles()) {
	case 0: //donothing
	case 1: //optimization for the common scenario with only style
		b.WriteString(` class="` + c.Styles()[0] + string('"')) //eg class="style1"
	default:
		b.WriteString(` class="` + strings.Join(c.Styles(), " ") + string('"')) // replace with \"
	}
	if c.RowSpan() > 1 {
		b.WriteString(` rowspan="` + strconv.Itoa(c.RowSpan()) + string('"')) //eg rowspan="3"
	}
	if c.ColSpan() > 1 {
		b.WriteString(` colspan=` + strconv.Itoa(c.ColSpan()) + string('"')) // eg colspan="2"
	}
	// txt, err := InlinedMdToHTML(c.Text(), nil)
	// if err != nil {
	// 	hr.htmlError = fmt.Errorf("error in parsing the following text: %s; error is %s ", strconv.Quote(c.Text()), err)
	// }
	//	if hr.settings.TrimCellContents {
	//txt = strings.TrimSpace(txt)
	//	}
	b.WriteString(">" + hr.renderText(c.Text()) + "</" + tag + ">\n") //eg "> text </td>"
	hr.write(b.String())
	return hr.Err()
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
