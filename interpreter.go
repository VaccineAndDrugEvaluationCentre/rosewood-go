package rosewood

import (
	"fmt"
	"io"

	"github.com/drgo/rosewood/parser"
	"github.com/drgo/rosewood/settings"
	"github.com/drgo/rosewood/types"
)

//Version of this library
const Version = "0.4.0"

//Settings is an alias for Rosewood settings
type Settings = settings.Settings

//Interpreter holds the state of a Rosewood interpreter
type Interpreter struct {
	settings *Settings
}

//NewInterpreter returns an initialized Rosewood interpreter
//
func NewInterpreter(settings *settings.Settings) *Interpreter {
	//if no custom settings use default ones
	if settings == nil {
		settings = DefaultSettings()
	}
	return &Interpreter{settings}
}

//Run takes an io.Reader streaming the contents of one or more Rosewood scripts
//and an io.Writer to output the formatted text.
func (ri *Interpreter) Run(src io.Reader, out io.Writer) error {
	file, err := ri.Parse(src, "")
	if err != nil {
		return err
	}
	return ri.RenderTables(out, file.Tables(), NewHTMLRenderer())
}

//Parse takes an io.Reader containing RoseWood script and an optional script identifier and returns parsed tables and an error
func (ri *Interpreter) Parse(r io.Reader, scriptIdentifer string) (*parser.File, error) {
	file := parser.NewFile(scriptIdentifer, ri.settings)
	if err := file.Parse(r); err != nil {
		return nil, err
	}
	if ri.settings.Debug {
		fmt.Printf("%d table(s) found\n", file.TableCount())
		tables := file.Tables()
		for i := 0; i < len(tables); i++ {
			fmt.Printf("%v\n", tables[i])
		}
	}
	return file, nil
}

//RenderTables renders 1 or more tables into a Writer using the passed Renderer
func (ri *Interpreter) RenderTables(w io.Writer, tables []*types.Table, hr types.Renderer) error {
	var err error
	hr.SetWriter(w)
	hr.SetSettings(ri.settings)
	hr.SetTables(tables)
	err = hr.StartFile()
	for _, t := range tables {
		if err = t.Run(); err != nil {
			return fmt.Errorf("failed to run one or more commands for table: %s", err)
		}
		if err = t.Render(w, hr); err != nil {
			return fmt.Errorf("failed to render table: %s", err)
		}
	}
	err = hr.EndFile()
	return err
}

//Settings returns currently active interpreter settings
func (ri *Interpreter) Settings() *settings.Settings {
	return ri.settings
}

//DefaultSettings returns a pointer to an initialized settings object
func DefaultSettings() *Settings {
	return settings.DefaultSettings()
}
