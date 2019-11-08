// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package latex

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/drgo/core/md"
	"github.com/drgo/core/ui"
	"github.com/drgo/rosewood"
	"github.com/drgo/rosewood/table"
	"github.com/drgo/rosewood/types"
)

const (
	header = `
\documentclass{article}
\usepackage{booktabs}
\usepackage{longtable}
\usepackage{crimson}
\usepackage[T1]{fontenc}
\usepackage{pdflscape}
\usepackage{multirow}
\usepackage[margin=1in]{geometry}
\usepackage{array}
\newcolumntype{L}[1]{>{\raggedright\arraybackslash}p{#1}}
\begin{document}
`

	footer = `
\end{document}
`
)

var defaultCSS []byte

//init, run automatically, registers Latex renderer with Rosewood
func init() {
	config := rosewood.RendererConfig{
		Name:     "latex",
		Renderer: makeLatexRenderer,
	}
	// load default css once per exe run
	// exeDir, err := files.GetExeDir()
	// if err == nil {
	// 	defaultCSS, _ = ioutil.ReadFile(filepath.Join(exeDir, defaultCSSFileName))
	// }
	rosewood.RegisterRenderer(&config)
}

//latexRenderer implements table.Renderer for Latex output
type latexRenderer struct {
	bw         io.Writer
	settings   *types.RosewoodSettings
	tables     []*table.Table
	props      *table.Properties // properties of current table
	currentRow *table.Row
	err        error //tracks errors
}

//makeLatexRenderer factory function according to the renderer registration requirements
func makeLatexRenderer() (table.Renderer, error) {
	return NewLatexRenderer()
}

//NewLatexRenderer create a new latexRenderer and return it as a Renderer
func NewLatexRenderer() (table.Renderer, error) {
	return &latexRenderer{}, nil
}

func (hr *latexRenderer) SetWriter(w io.Writer) error {
	hr.bw = w
	return nil
}

func (hr *latexRenderer) SetSettings(settings *types.RosewoodSettings) error {
	hr.settings = settings
	return nil
}

func (hr *latexRenderer) SetTables(tables []*table.Table) error {
	hr.tables = tables
	return nil
}

func (hr *latexRenderer) Err() error {
	return hr.err
}

// write does all the writing to the writer and handles errors by stopping any further writing
// and returning the error
func (hr *latexRenderer) write(s string) error { //TODO:optimize
	if hr.err == nil {
		_, hr.err = hr.bw.Write([]byte(s))
	}
	return hr.err
}

func (hr *latexRenderer) StartFile() error {
	var b strings.Builder //optimization for golang >= 1.10
	b.Grow(1024 * 100)    //preallocate 100kb to avoid additional allocations
	b.WriteString(header)
	b.WriteString("%%date-generated" + time.Now().Format("2006-01-02 15:04:05") + "\n")
	if hr.settings.Debug >= ui.DebugAll {
		b.WriteString(time.Now().Format("2006-01-02 15:04:05") + "\n")
	}
	hr.write(b.String())
	return hr.Err()
}

func (hr *latexRenderer) EndFile() error {
	return hr.write(footer)
}

func (hr *latexRenderer) StartTable(t *table.Table) error {
	hr.props = t.Properties()
	hr.write(`
\begin{table}[htbp]
\centering
\setlength{\tabcolsep}{0.5em}
\def\arraystretch{1.5}  
`)
	if t.Caption != nil {
		hr.write(`\caption{`)
		for _, line := range t.Caption.Lines {
			hr.write(hr.renderText(line))
		}
		hr.write("}\n")
	}
	hr.write(`\begin{tabular}{` + hr.props.TableSpecs + "}\n")
	hr.write(`\toprule` + "\n")
	return hr.Err()
}

func (hr *latexRenderer) EndTable(t *table.Table) error {
	hr.write(`
\bottomrule
\end{tabular}
\label{table:mr}
\end{table}
`)
	// if t.Footnotes != nil {
	// 	hr.write(`<div class="rw-footnotes">` + "\n")
	// 	for _, line := range t.Footnotes.Lines {
	// 		hr.write(hr.renderText(line) + "<br>\n")
	// 	}
	// 	hr.write("</div>\n")
	// }
	return hr.Err()
}

func (hr *latexRenderer) StartRow(r *table.Row) error {
	hr.currentRow = r
	return hr.Err()
}

func (hr *latexRenderer) EndRow(r *table.Row) error {
	hr.write(`\\` + "\n")
	if hr.currentRow.Number() == hr.props.HeaderRowsCount { //last row in header
		hr.write(`\midrule` + "\n")
	}
	return hr.Err()
}

func (hr *latexRenderer) OutputCell(c *table.Cell) error {
	if c.State() == table.CsHMerged { //skip col merged cells
		return hr.Err()
	}
	var b strings.Builder
	b.Grow(1024)
	// b.WriteString("  <" + tag) //open td or th tag
	// // write styles
	// switch len(c.Styles()) {
	// case 0: //donothing
	// case 1: //optimization for the common scenario with only one style
	// 	b.WriteString(` class="` + c.Styles()[0] + `"`) //eg class="style1"
	// default:
	// 	b.WriteString(` class="` + strings.Join(c.Styles(), " ") + `"`)
	// }
	// if c.RowSpan() > 1 {
	// 	b.WriteString(fmt.Sprintf(" rowspan=\"%d\"", c.RowSpan())) //eg rowspan="3"
	// }
	if c.ColSpan() > 1 {
		b.WriteString(`\multicolumn{` + strconv.Itoa(c.ColSpan()) + `}{c}{`)
	}
	if c.RowSpan() > 1 {
		b.WriteString(`\multirow{` + strconv.Itoa(c.RowSpan()) + `}{*}{`)
	}
	// trim cell contents b/c html ignores white space anyway
	b.WriteString(hr.renderText(c.Text()))
	if c.RowSpan() > 1 {
		b.WriteString("}")
	}
	if c.ColSpan() > 1 {
		b.WriteString("}")
	}
	if !c.LastVisibleCell() {
		b.WriteString("&")
	}
	hr.write(b.String())
	return hr.Err()
}

func (hr *latexRenderer) renderText(s string) string {
	return md.EscapeAsTex(s)
	switch hr.settings.MarkdownRender {
	case "standard", "":
		txt, _ := md.InlinedMdToHTML(s, nil)
		return string(txt)
	case "strict":
		txt, err := md.InlinedMdToHTML(s, nil)
		if err != nil {
			hr.err = fmt.Errorf("error in parsing the following text: %s; error is %s ", strconv.Quote(s), err)
		}
		return string(txt)
	default: //including "disabled"
		return s
	}
}
