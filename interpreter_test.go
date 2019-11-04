// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package rosewood_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/drgo/core/debug"
	"github.com/drgo/rosewood"
	_ "github.com/drgo/rosewood/renderers/html" //include needed renderers
	"github.com/drgo/rosewood/types"
)

func TestInterpreter_Run(t *testing.T) {
	const pathPrefix = "test-files/"
	tests := []struct {
		srcFileName string
		outFileName string
		wantW       string
		wantErr     bool
	}{
		{"singletab.gold", "singletab", "", false},
		{"wrong1tab.gold", "", "", true},
	}
	job := rosewood.DefaultJob(types.DebugRosewoodSettings(debug.DebugAll))
	for _, tt := range tests {
		fmt.Println(strings.Repeat("*", 40))
		t.Run(tt.srcFileName, func(t *testing.T) {
			r, err := os.Open(pathPrefix + tt.srcFileName)
			if err != nil {
				t.Fatalf("could not open file [%s]: %s", tt.srcFileName, err)
			}
			defer r.Close()
			w := &bytes.Buffer{} // output
			err = rosewood.ToHTML(r.Name(), job, r, w)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Interpreter.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				fmt.Printf("error parsing file [%s]: %s\n", tt.srcFileName, err)
				return
			}
			fmt.Fprintf(os.Stderr, w.String())
		})
	}
}
