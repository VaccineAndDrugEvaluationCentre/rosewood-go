package rosewood

import (
	"fmt"
	"io"

	"github.com/drgo/rosewood/parser"
	"github.com/drgo/rosewood/types"
	"github.com/drgo/rosewood/utils"
)

//VERSION of this library
const VERSION = "0.3.5"

type Interpreter struct {
	fileName string
	settings *utils.Settings
	tables   []*types.Table
	file     *parser.RwFile
}

func NewInterpreter(settings *utils.Settings) *Interpreter {
	ri := &Interpreter{}
	//if no custom settings use default ones
	if settings == nil {
		ri.settings = utils.DefaultSettings()
		if ri.settings == nil {
			panic("Interpreter failed to load settings")
		}
	} else {
		ri.settings = settings
	}
	return ri
}

//Errors returns a list of rosewood.EmError
func (ri *Interpreter) Errors(err error) (eList []error) {
	//	fmt.Printf("%T\n", err)
	return ri.file.Errors(err)
}

//Parse takes an io.Reader containing RoseWood script and an optional script identifier and return an error
func (ri *Interpreter) Parse(r io.Reader, scriptIdentifer string) error {
	ri.file = parser.NewRwFile(scriptIdentifer, ri.settings)
	var err error
	ri.tables, err = ri.file.Parse(r)
	if len(ri.tables) == 0 {
		//return utils.NewError(utils.ErrUnknown, , "unknown error in Interpreter.CreateTables()")
	}
	return err
}

func (ri *Interpreter) renderTables(w io.Writer, hr types.Renderer) error {
	var err error
	hr.SetWriter(w)
	hr.SetSettings(ri.settings)
	hr.SetTables(ri.tables)
	hr.StartFile()
	for _, t := range ri.tables {
		if err := t.Run(); err != nil {
			return fmt.Errorf("failed to run commands for table %s", err)
		}
		err = t.Render(w, hr)
	}
	hr.EndFile()
	return err
}

//Run takes an io.Reader streaming the contents of one or more Rosewood scripts
//and an io.Writer to output the formatted text.
func (ri *Interpreter) Run(src io.Reader, out io.Writer) error {
	var err error
	if err = ri.Parse(src, ""); err != nil {
		return err
	}
	if err = ri.renderTables(out, NewHtmlRenderer()); err != nil {
		return err
	}
	return nil
}

//Settings exported
type Settings = utils.Settings

func DefaultSettings() *utils.Settings {
	return utils.DefaultSettings()
}
