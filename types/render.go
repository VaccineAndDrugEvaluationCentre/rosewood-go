package types

import (
	"io"

	"github.com/drgo/rosewood/utils"
)

type Renderer interface {
	SetWriter(io.Writer) error
	SetSettings(*utils.Settings) error
	SetTables([]*Table) error
	StartFile() error
	EndFile() error

	StartTable(t *Table) error

	EndTable(t *Table) error

	StartRow(r *Row) error

	EndRow(r *Row) error
	OutputCell(c *Cell) error
}
