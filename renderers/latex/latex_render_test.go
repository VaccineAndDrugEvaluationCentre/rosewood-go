// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package latex

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/drgo/core/ui"
	"github.com/drgo/rosewood"
	"github.com/drgo/rosewood/parser"
	"github.com/drgo/rosewood/types"
)

func TestLatexRenderer(t *testing.T) {
	const pathPrefix = "../../test-files/"
	tests := []struct {
		srcFileName string
		outFileName string
		wantW       string
		wantErr     bool
	}{
		{"singletab.gold", "singletab", "", false},
		{"wrong1tab.gold", "", "", true},
	}
	job := rosewood.DefaultJob(types.DebugRosewoodSettings(ui.DebugAll))
	for _, tt := range tests {
		fmt.Println(strings.Repeat("*", 40))
		t.Run(tt.srcFileName, func(t *testing.T) {
			ri := rosewood.NewInterpreter(job).SetScriptIdentifer(pathPrefix + tt.srcFileName)
			file, err := tParseFile(t, ri)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Interpreter.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				fmt.Printf("error parsing file [%s]: %s\n", tt.srcFileName, err)
				return
			}
			hr, err := rosewood.GetRendererByName("latex")
			if err != nil {
				t.Fatalf(" rosewood.GetRendererByName error = %v", err)
			}
			out := &bytes.Buffer{} // output
			err = ri.ReportError(ri.Render(out, file, hr))
			if err != nil {
				t.Fatalf(" rosewood.Interepeter.Render error = %v", err)
			}
			fmt.Fprintf(os.Stderr, out.String())
		})
	}
}

func tParseFile(t *testing.T, ri *rosewood.Interpreter) (*parser.File, error) {
	r, err := os.Open(ri.ScriptIdentifer())
	if err != nil {
		t.Fatalf("could not open file [%s]: %s", ri.ScriptIdentifer(), err)
	}
	defer r.Close()
	file, err := ri.Parse(r, ri.ScriptIdentifer())
	if err != nil || ri.Settings().CheckSyntaxOnly {
		return nil, ri.ReportError(err)
	}
	return file, nil
}
