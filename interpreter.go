// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package rosewood

import (
	"bufio"
	"fmt"
	"io"

	"github.com/drgo/core/errors"
	"github.com/drgo/core/ui"
	"github.com/drgo/rosewood/parser"
	"github.com/drgo/rosewood/table"
	"github.com/drgo/rosewood/types"
)

//Interpreter holds the state of a Rosewood interpreter
type Interpreter struct {
	job             *Job
	settings        *Settings
	scriptIdentifer string
}

//NewInterpreter returns an initialized Rosewood interpreter
func NewInterpreter(job *Job) *Interpreter {
	//if no custom settings use default ones
	if job == nil || job.RosewoodSettings == nil {
		panic("rosewood.NewInterpreter: job and job.RosewoodSettings must not be null")
	}
	// set debug flag of internal packages
	//table.SetDebug(ui.DebugAll)
	return &Interpreter{job, job.RosewoodSettings, ""}
}

//Parse takes an io.Reader containing RoseWood script and an optional script identifier and returns
// parsed tables and an error
func (ri *Interpreter) Parse(r io.ReadSeeker, scriptIdentifer string) (*parser.File, error) {
	file := parser.NewFile(scriptIdentifer, ri.job)
	if err := file.Parse(r); err != nil {
		return nil, err
	}
	if ri.job.RunOptions.Debug == ui.DebugAll {
		ri.job.UI.Logf("***Parsing finished: %d table(s) found\n", file.TableCount())
		for i, t := range file.Tables() {
			ri.job.UI.Logf("****Contents of table %d\n%v\n", i+1, t)
		}
	}
	return file, nil
}

//Render renders 1 or more tables into a Writer using the passed Renderer
func (ri *Interpreter) Render(w io.Writer, file *parser.File, hr table.Renderer) error {
	var err error
	bw := bufio.NewWriter(w) //buffer the writer to speed up writing
	tables := file.Tables()
	_ = hr.SetWriter(bw)
	if err = hr.SetSettings(ri.settings); err != nil {
		return fmt.Errorf("failed to render table: %s", err)
	}
	_ = hr.SetTables(tables)
	if err = hr.StartFile(); err != nil {
		return fmt.Errorf("failed to render table: %s", err)
	}
	for i, t := range tables {
		if err = t.Run(); err != nil {
			return fmt.Errorf("failed to run one or more commands for table: %s", err)
		}
		ri.job.UI.Logf("****processed contents of table %d\n%v\n", i+1, t.ProcessedTableContents().DebugString())
		if err = t.Render(w, hr); err != nil {
			return fmt.Errorf("failed to render table %d: %s", i+1, err)
		}
	}
	if err = hr.EndFile(); err != nil {
		return fmt.Errorf("failed to render table: %s", err)
	}
	return bw.Flush() //flush to ensure all changes are written to the writer
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
