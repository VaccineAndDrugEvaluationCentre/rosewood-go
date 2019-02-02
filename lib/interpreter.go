// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package rosewood

import (
	"bufio"
	"fmt"
	"io"

	"github.com/drgo/core/errors"
	"github.com/drgo/rosewood/lib/parser"
	"github.com/drgo/rosewood/lib/table"
	"github.com/drgo/rosewood/lib/types"
)

//Version of this library
const version = "0.5.1"

//LibVersion version of Rosewood lib
func LibVersion() string {
	return version
}

//Interpreter holds the state of a Rosewood interpreter
type Interpreter struct {
	settings        *Settings
	scriptIdentifer string
}

//NewInterpreter returns an initialized Rosewood interpreter
func NewInterpreter(settings *types.RosewoodSettings) *Interpreter {
	//if no custom settings use default ones
	if settings == nil {
		settings = DefaultSettings()
	}
	// set debug flag of internal packages
	table.SetDebug(types.DebugAll)
	return &Interpreter{settings, ""}
}

//Parse takes an io.Reader containing RoseWood script and an optional script identifier and returns
// parsed tables and an error
func (ri *Interpreter) Parse(r io.ReadSeeker, scriptIdentifer string) (*parser.File, error) {
	file := parser.NewFile(scriptIdentifer, ri.settings)
	if err := file.Parse(r); err != nil {
		return nil, err
	}
	if ri.settings.Debug == types.DebugAll {
		fmt.Printf("***Parsing finished: %d table(s) found\n", file.TableCount())
		for i, t := range file.Tables() {
			fmt.Printf("****Contents of table %d\n", i+1)
			fmt.Printf("%v\n", t)
		}
	}
	return file, nil
}

//Render renders 1 or more tables into a Writer using the passed Renderer
func (ri *Interpreter) Render(w io.Writer, file *parser.File, hr table.Renderer) error {
	var err error
	bw := bufio.NewWriter(w) //buffer the writer to speed up writing
	tables := file.Tables()
	hr.SetWriter(bw)
	hr.SetSettings(ri.settings)
	hr.SetTables(tables)
	err = hr.StartFile()
	for i, t := range tables {
		if err = t.Run(); err != nil {
			return fmt.Errorf("failed to run one or more commands for table: %s", err)
		}
		if ri.settings.Debug == types.DebugAll {
			fmt.Printf("****processed contents of table %d\n", i+1)
			fmt.Printf("%v\n", t.ProcessedTableContents().DebugString())
		}
		if err = t.Render(w, hr); err != nil {
			return fmt.Errorf("failed to render table %d: %s", i+1, err)
		}
	}
	err = hr.EndFile()
	bw.Flush() //flush to ensure all changes are written to the writer
	return err
}

//ReportError returns a list of errors encountered during running
func (ri *Interpreter) ReportError(err error) error {
	return errors.ErrorsToError(err)
}

//ScriptIdentifer returns currently processed ScriptIdentifer
func (ri *Interpreter) ScriptIdentifer() string {
	return ri.scriptIdentifer
}

//SetScriptIdentifer sets the name of the running script
func (ri *Interpreter) SetScriptIdentifer(scriptIdentifer string) *Interpreter {
	ri.scriptIdentifer = scriptIdentifer
	return ri
}

//Settings returns currently active interpreter settings
func (ri *Interpreter) Settings() *types.RosewoodSettings {
	return ri.settings
}

//ConvertToCurrentVersion utility to convert older versions of Rosewood to current version
func ConvertToCurrentVersion(settings *types.RosewoodSettings, in io.Reader, out io.Writer) error {
	switch settings.ConvertFromVersion {
	case "v0.1":
		return parser.ConvertToCurrentVersion(settings, parser.RWSyntaxVdotzero1, in, out)
	}
	return fmt.Errorf("invalid version number: %s", settings.ConvertFromVersion)
}

//TODO: remove old code
// //RenderTables renders 1 or more tables into a Writer using the passed Renderer
// func (ri *Interpreter) RenderTables(w io.Writer, tables []*types.Table, hr table.Renderer) error {
// 	var err error
// 	hr.SetWriter(w)
// 	hr.SetSettings(ri.settings)
// 	hr.SetTables(tables)
// 	err = hr.StartFile()
// 	for _, t := range tables {
// 		if err = t.Run(); err != nil {
// 			return fmt.Errorf("failed to run one or more commands for table: %s", err)
// 		}
// 		if err = t.Render(w, hr); err != nil {
// 			return fmt.Errorf("failed to render table: %s", err)
// 		}
// 	}
// 	err = hr.EndFile()
// 	return err
// }

// //Run takes an io.Reader streaming the contents of one or more Rosewood scripts
// //and an io.Writer to output the formatted text.
// func (ri *Interpreter) Run(src io.Reader, out io.Writer) error {
// 	file, err := ri.Parse(src, "")
// 	if err != nil {
// 		return err
// 	}
// 	return ri.RenderTables(out, file.Tables(), html.NewHTMLRenderer())
// }
