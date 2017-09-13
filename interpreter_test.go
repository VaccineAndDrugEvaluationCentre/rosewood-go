package rosewood

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/drgo/rosewood/utils"
)

func TestInterpreter_Run(t *testing.T) {
	const pathPrefix = "test-files/"
	tests := []struct {
		srcFileName string
		outFileName string
		settings    *utils.Settings
		wantW       string
		wantErr     bool
	}{
		{"correct1tab.rw", "correct1tab", utils.DebugSettings(true), "", false},
		{"wrong1tab.rw", "", utils.DebugSettings(true), "", true},
	}
	for _, tt := range tests {
		t.Run(tt.srcFileName, func(t *testing.T) {
			ri := NewInterpreter(tt.settings)
			r, err := os.Open(pathPrefix + tt.srcFileName)
			defer r.Close()
			if err != nil {
				t.Fatalf("could not open file [%s]: %s", tt.srcFileName, err)
			}
			fmt.Println(strings.Repeat("*", 40))
			w := &bytes.Buffer{}
			err = ri.Run(r, w)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Interpreter.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				fmt.Printf("error running file [%s]: %s\n", tt.srcFileName, err)
			}
			//fmt.Println(w.String())
			// if gotW := w.String(); gotW != tt.wantW {
			// 	t.Errorf("Interpreter.Run() = %v, want %v", gotW, tt.wantW)
			// }
			if tt.outFileName != "" {
				fn := path.Join(pathPrefix, tt.outFileName+"."+testFileExt)
				if err := ioutil.WriteFile(fn, w.Bytes(), 0644); err != nil {
					t.Errorf("failed to write to file %s: %v", fn, err)
				}
				fmt.Printf("Results saved to file://%s\n", fn)
			}
		})
	}
}

// //parseFile convenience function to parse a file
// func parseFile(ri *Interpreter, filename string) error {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse file %s", err)
// 	}
// 	defer file.Close()
// 	return ri.Parse(file, filename)
// }
