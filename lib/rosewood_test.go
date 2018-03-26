//rosewood_test contains integration (black box tests) for // Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package rosewood
// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package rosewood_test

// import (
// 	"bytes"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path"
// 	"testing"

// 	"github.com/drgo/rosewood"
// )

// func TestInterpreter_Run(t *testing.T) {
// 	const pathPrefix = "test-files/"
// 	const testFileExt = "html"
// 	testSettings := func() *rosewood.Settings { return rosewood.DefaultSettings() }
// 	tests := []struct {
// 		srcFileName string
// 		outFileName string
// 		settings    *rosewood.Settings
// 		wantW       string
// 		wantErr     bool
// 	}{
// 		{"correct1tab.rw", "correct1tab", testSettings(), "", false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.srcFileName, func(t *testing.T) {
// 			ri := rosewood.NewInterpreter(tt.settings)
// 			r, err := os.Open(pathPrefix + tt.srcFileName)
// 			defer r.Close()
// 			if err != nil {
// 				t.Fatalf("could not open file [%s]: %s", tt.srcFileName, err)
// 			}
// 			w := &bytes.Buffer{}
// 			if err = ri.Run(r, w); (err != nil) != tt.wantErr {
// 				t.Fatalf("Interpreter.Run() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			//fmt.Println(w.String())
// 			// if gotW := w.String(); gotW != tt.wantW {
// 			// 	t.Errorf("Interpreter.Run() = %v, want %v", gotW, tt.wantW)
// 			// }
// 			if tt.outFileName != "" {
// 				fn := path.Join(pathPrefix, tt.outFileName+"."+testFileExt)
// 				if err := ioutil.WriteFile(fn, w.Bytes(), 0644); err != nil {
// 					t.Errorf("failed to write to file %s: %v", fn, err)
// 				}
// 				fmt.Printf("Results saved to file://%s\n", fn)
// 			}
// 		})
// 	}
// }
