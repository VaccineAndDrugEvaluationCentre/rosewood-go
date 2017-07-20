package carpenter

import "io"

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
	htmlOpenTable = `
<table class="">`
	htmlCloseTable = `
</table>`
	htmlOpenRow = `
<tr>`
	htmlCloseRow = `
</tr>`
)

type HtmlRenderer struct {
	w io.Writer
}

func NewHtmlRenderer() *HtmlRenderer {
	return &HtmlRenderer{}
}

func (hr *HtmlRenderer) SetWriter(w io.Writer) error {
	hr.w = w
	return nil
}

func (hr *HtmlRenderer) StartFile(t *Table) error {
	io.WriteString(hr.w, htmlHeader)
	return nil
}

func (hr *HtmlRenderer) EndFile(t *Table) error {
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
	io.WriteString(hr.w, "<td>"+c.String()+"</td>")
	return nil
}

func Render(w io.Writer, r *HtmlRenderer, t *Table) error {
	r.SetWriter(w)

	r.StartFile(t)
	r.StartTable(t)
	for i := 0; i < len(t.rows); i++ {
		r.StartRow(t.rows[i])
		for j := 0; j < len(t.rows[i].cells); j++ {
			r.OutputCell(t.rows[i].cells[j])
		}
		r.EndRow(t.rows[i])
	}
	r.EndTable(t)
	r.EndFile(t)
	return nil
}
