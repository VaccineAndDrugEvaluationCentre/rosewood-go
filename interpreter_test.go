package rosewood

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestInterpreter_Run(t *testing.T) {
	const pathPrefix = "test-files/"
	testSettings := func() *Settings { return DefaultSettings() }
	tests := []struct {
		srcFileName string
		outFileName string
		settings    *Settings
		wantW       string
		wantErr     bool
	}{
		{"correct1tab.rw", "correct1tab", testSettings(), "", false},
	}
	for _, tt := range tests {
		t.Run(tt.srcFileName, func(t *testing.T) {
			ri := NewInterpreter(tt.settings)
			r, err := os.Open(pathPrefix + tt.srcFileName)
			defer r.Close()
			if err != nil {
				t.Fatalf("could not open file [%s]: %s", tt.srcFileName, err)
			}
			w := &bytes.Buffer{}
			if err = ri.Run(r, w); (err != nil) != tt.wantErr {
				t.Fatalf("Interpreter.Run() error = %v, wantErr %v", err, tt.wantErr)
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

func ExampleNewInterpreter() {
	ri := NewInterpreter(nil)
	if err := parseFile(ri, "test-files/correct1tab.rw"); err != nil {
		trace.Printf("error parsing file: %s\n", err)
	}
	fmt.Println(ri.sectionCount())
	if ri.sectionCount() != 4 {
		return
	}
	trace.Printf("%d\n", ri.sections[0].offset)
	fmt.Printf("%d\n", ri.sections[2].offset)
	fmt.Printf("%d\n", ri.sections[0].LineCount())
	fmt.Printf("%d\n", ri.sections[3].LineCount())
	fmt.Printf("%d\n", len(ri.tables))
	if len(ri.tables) == 0 {
		return
	}
	t := ri.tables[0]
	fmt.Printf("#rows:%d\ncells per row\n", t.contents.rowCount())
	for _, r := range t.contents.rows {
		fmt.Printf("%d\n", len(r.cells))
	}

	// Output:
	// 4
	// 2
	// 11
	// 1
	// 5
	// 1
	// #rows:6
	// cells per row
	// 1
	// 3
	// 4
	// 4
	// 4
	// 4
}

func ExampleNewInterpreter2() {
	ri := NewInterpreter(nil)
	if err := parseFile(ri, "test-files/correct2tabs.rw"); err != nil {
		fmt.Printf("error parsing file: %s\n", err)
	}
	fmt.Println(ri.sectionCount())
	if ri.sectionCount() == 8 {
		fmt.Printf("%d\n", ri.sections[0].offset)
		fmt.Printf("%d\n", ri.sections[2].offset)
		fmt.Printf("%d\n", ri.sections[0].LineCount())
		fmt.Printf("%d\n", ri.sections[3].LineCount())
	}
	// Output:
	//8
	//2
	//11
	//1
	//4
}

//parseFile convenience function to parse a file
func parseFile(ri *Interpreter, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to parse file %s", err)
	}
	defer file.Close()
	return ri.Parse(file, filename)
}
