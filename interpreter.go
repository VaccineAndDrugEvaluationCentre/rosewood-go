package rosewood

import (
	"fmt"
	"io"

	"github.com/drgo/rosewood/parser"
	"github.com/drgo/rosewood/types"
	"github.com/drgo/rosewood/utils"
)

//Version of this library
const Version = "0.4.0"

//Interpreter holds the state of a Rosewood interpreter
type Interpreter struct {
	// file     *parser.File
	settings *utils.Settings
}

//NewInterpreter returns an initialized Rosewood interpreter
//
func NewInterpreter(settings *utils.Settings) *Interpreter {
	ri := &Interpreter{}
	//if no custom settings use default ones
	if settings == nil {
		ri.settings = utils.DefaultSettings()
	} else {
		ri.settings = settings
	}
	return ri
}

//Run takes an io.Reader streaming the contents of one or more Rosewood scripts
//and an io.Writer to output the formatted text.
func (ri *Interpreter) Run(src io.Reader, out io.Writer) error {
	file, err := ri.Parse(src, "")
	if err != nil {
		return err
	}
	if err = ri.RenderTables(out, file.Tables(), NewHTMLRenderer()); err != nil {
		return err
	}
	return nil
}

// //Errors returns a list of parsing and run errors
// func (ri *Interpreter) Errors() (eList []error) {
// 	return ri.file.Errors()
// }

// //Err returns a compact list of parsing and run errors
// func (ri *Interpreter) Err() error {
// 	return ri.file.Err()
// }

//Parse takes an io.Reader containing RoseWood script and an optional script identifier and returns parsed tables and an error
func (ri *Interpreter) Parse(r io.Reader, scriptIdentifer string) (*parser.File, error) {
	file := parser.NewFile(scriptIdentifer, ri.settings)
	if err := file.Parse(r); err != nil {
		return nil, err
	}
	// if len(file.tables) == 0 {
	// 	return nil, fmt.Errorf("unknown error in Interpreter.Parse()")
	// }
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
			return fmt.Errorf("failed to run commands for table %s", err)
		}
		if err = t.Render(w, hr); err != nil {
			return fmt.Errorf("failed to render table %s", err)
		}
	}
	err = hr.EndFile()
	return err
}

//Settings returns currently active interpreter settings
func (ri *Interpreter) Settings() *utils.Settings {
	return ri.settings
}

//Settings holds Rosewood settings
type Settings = utils.Settings

//DefaultSettings returns default settings
func DefaultSettings() *utils.Settings {
	return utils.DefaultSettings()
}
