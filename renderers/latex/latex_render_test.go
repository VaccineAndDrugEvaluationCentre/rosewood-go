// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package latex

import (
	"bytes"
	"fmt"
	"io"
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
		{"big-table-long.txt", "bigtab", "", false},
		{"bug-v1tov2-extra-style-converted-v1-2-v2rw", "bug-v1v2", "", false},
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
			io.Copy(os.Stderr, out)
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
