package types

import (
	"io"

	"github.com/drgo/rosewood/utils"
)

//Renderer is an interface that groups all functions needed for rendering a Rosewood table
//see html_render.go for an implementation that renders tables into HTML
type Renderer interface {
	SetWriter(io.Writer) error
	SetSettings(*utils.Settings) error
	SetTables([]*Table) error
	Err() error
	StartFile() error
	EndFile() error

	StartTable(t *Table) error

	EndTable(t *Table) error

	StartRow(r *Row) error

	EndRow(r *Row) error
	OutputCell(c *Cell) error
}
