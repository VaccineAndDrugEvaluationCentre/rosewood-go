// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package rosewood

import (
	"fmt"
	"io"

	"github.com/drgo/rosewood/parser"
	"github.com/drgo/rosewood/types"
)

//Version of this library
const version = "0.5.6"

//LibVersion version of Rosewood lib
func LibVersion() string {
	return version
}

//ConvertToCurrentVersion utility to convert older versions of Rosewood to current version
func ConvertToCurrentVersion(settings *types.RosewoodSettings, in io.Reader, out io.Writer) error {
	switch settings.ConvertFromVersion {
	case "v0.1":
		return parser.ConvertToCurrentVersion(settings, parser.RWSyntaxVdotzero1, in, out)
	}
	return fmt.Errorf("invalid version number: %s", settings.ConvertFromVersion)
}

// ToHTML runs a task
func ToHTML(inputFileName string, job *Job, in io.ReadSeeker, out io.Writer) error {
	ri := NewInterpreter(job).SetScriptIdentifer(inputFileName)
	file, err := ri.Parse(in, ri.ScriptIdentifer())
	if err != nil || ri.Settings().CheckSyntaxOnly {
		return ri.ReportError(err)
	}
	hr, err := GetRendererByName("html")
	if err != nil {
		return err
	}
	return ri.ReportError(ri.Render(out, file, hr))
}
