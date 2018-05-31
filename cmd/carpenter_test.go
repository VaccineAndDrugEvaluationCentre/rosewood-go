// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package main

import (
	"os"
	"strings"
	"testing"
)

//TODO: fix tests for empty and binary files: returned text is different than specified in tests
//TODO: add tests for rw v0.1
func Test_run(t *testing.T) {
	const (
		path = "../test-files" //all test files used below are stored here
		exe  = "carpenter "    //keep the space at the end
	)
	tests := []struct {
		name string //unique test identifer
		args string //simulate command line eg for "carpenter run ../test-files/empty.rw -f"
		// args should be "run %path/empty.rw -f"
		wantErr bool   //true if an error is expected
		errMsg  string //a fragment of the expected error description
	}{
		//test command arguments and flags
		{"no-cmd", "", true, ErrWrongCommand},
		{"wrong-cmd", "rnu", true, ""},
		{"version-cmd-ok", "version", false, ""},
		{"help-cmd-ok", "version", false, ""},
		{"wrong-global-flag", " -x run", true, "flag does not exist"},
		{"wrong-run-flag", "run %path/empty.gold -x", true, "flag does not exist"},

		{"missing-infile", "run ", true, ErrMissingInFile},
		{"outfile-exists", "run %path/singletab.gold", true, "singletab.html: file exists"},

		{"no-binary", "run " + exe, true, "file does not contain text (possibly a binary file)"},
		{"no-SmallBinFile", "check %path/smallbinfile.gold", true, "file does not contain text (possibly a binary file)"},
		{"run-empty-File", "run %path/empty.gold", true, "stream is empty"},
		{"check-empty-File", "check %path/empty.gold", true, "stream is empty"},
		{"run-empty-File", "run %path/notab.gold -r", true, "empty table"},
		{"check-empty-File", "check %path/notab.gold", true, "empty table"},
		{"run-section-Sep-Only", "run %path/sectionsepsonly.gold -r", true, "empty table"},
		{"check-section-Sep-Only", "check %path/sectionsepsonly.gold", true, "empty table"},
		{"check-command-Only", "check %path/commandonly.gold", true, "empty table"},
		{"too-many-scections", "check %path/toomanyscections.gold", true, "incorrect number of sections"},

		{"check-syntax-errors", "check %path/wrong1tab.gold", true, "unknown command xmerge"},
		{"check-3-dash-sep", "check %path/oldsyntax.gold", true, "does not start by a section separator: +++"}, //fails because of --- separator
		{"check-old-syntax", "check %path/oldsyntax.gold -S=---", true, "unknown command emphasize"},           //does not fail because of --- separator

		//All acceptable
		{"run-verysmalltab", "run %path/verysmalltab.gold -r", false, ""},
		{"run-oneCorrectTable", "run %path/singletab.gold -r", false, ""},
		{"run-typicalTable1", "run %path/TableT1A.gold -r", false, ""},
		{"run-two-tables", "run %path/correct2tabs.gold -r", false, ""},

		//crashes
		{"run-typicalTable1", "run %path/crash1.gold -r", true, "invalid coordinates"}, //passes check but fails on rendering
		// {"check-singleentryandcommandonly", "c", path + "singleentryandcommandonly.rw", false, "singleentryandcommandonly.rw]: file does not contain text"},
		// {"check-oneWrongTable", "c", path + "wrong1tab.rw", true, "wrong1tab.rw]: syntax error : unknown command xmerge"},
		// {"check-oneCorrectTable", "c", path + "singletab.rw", false, ""},
		// {"check-oneCorrectTableOnly", "c", path + "tabonly.rw", false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args = strings.Replace(tt.args, "%path", path, -1)
			tt.args = exe + tt.args //prepend exe name (args[0])
			os.Args = strings.Fields(tt.args)
			t.Logf("argument: %s\n", tt.args)
			err := RunApp()
			if (err != nil) != tt.wantErr {
				t.Errorf("RunApp() error = %v, wantErr %v\n", err, tt.wantErr)
				t.Logf("error: %v\n", err)
			}
			//a bit fragile as it depends on the exact error text not changing in the code
			if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("wrong error message, expected= %s, got %s\n", tt.errMsg, err.Error())
			}
		})
	}
}

// //setupSettings convert a semicolon separated string of valid Carpenter flags into a rosewood.Settings object
// func setupSettings(flagText string) *rosewood.Settings {
// 	s := rosewood.DefaultSettings()
// 	flags := strings.Split(strings.ToLower(flagText), ";")
// 	if len(flags) == 0 {
// 		return s
// 	}
// 	for _, flag := range flags {
// 		switch flag {
// 		case "c", "check":
// 			s.CheckSyntaxOnly = true
// 		case "r", "replace":
// 			s.OverWriteOutputFile = true
// 		default: //anything else is output file name
// 			s.OutputFileName = flag
// 		}
// 	}
// 	return s
// }
