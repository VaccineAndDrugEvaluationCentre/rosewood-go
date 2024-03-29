// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package table

import (
	"io"

	"github.com/drgo/rosewood/types"
)

//Renderer is an interface that groups all functions needed for rendering a Rosewood table
//see html_render.go for an implementation that renders tables into HTML
type Renderer interface {
	SetWriter(io.Writer) error
	SetSettings(settings *types.RosewoodSettings) error
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
