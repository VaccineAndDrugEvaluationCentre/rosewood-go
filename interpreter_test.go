// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package rosewood

// func TestInterpreter_Run(t *testing.T) {
// 	const pathPrefix = "test-files/"
// 	tests := []struct {
// 		srcFileName string
// 		outFileName string
// 		settings    *types.RosewoodSettings
// 		wantW       string
// 		wantErr     bool
// 	}{
// 		{"singletab.rw", "singletab", types.DebugRosewoodSettings(true), "", false},
// 		{"wrong1tab.rw", "", types.DebugRosewoodSettings(true), "", true},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.srcFileName, func(t *testing.T) {
// 			ri := NewInterpreter(tt.settings)
// 			r, err := os.Open(pathPrefix + tt.srcFileName)
// 			defer r.Close()
// 			if err != nil {
// 				t.Fatalf("could not open file [%s]: %s", tt.srcFileName, err)
// 			}
// 			fmt.Println(strings.Repeat("*", 40))
// 			file, err := ri.Parse(r, tt.srcFileName)
// 			if (err != nil) != tt.wantErr {
// 				t.Fatalf("Interpreter.Parse() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if err != nil {
// 				fmt.Printf("error parsing file [%s]: %s\n", tt.srcFileName, err)
// 				return
// 			}
// 			w := &bytes.Buffer{}
// 			err = ri.RenderTables(w, file.Tables(), NewHTMLRenderer())
// 			if (err != nil) != tt.wantErr {
// 				t.Fatalf("Interpreter.RenderTables() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if err != nil {
// 				fmt.Printf("error rendering file [%s]: %s\n", tt.srcFileName, err)
// 			}
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

// //parseFile convenience function to parse a file
// func parseFile(ri *Interpreter, filename string) error {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse file %s", err)
// 	}
// 	defer file.Close()
// 	return ri.Parse(file, filename)
// }
