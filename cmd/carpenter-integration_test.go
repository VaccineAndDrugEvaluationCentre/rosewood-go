// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package main

import (
	"os"
	"strings"
	"testing"
)

//TODO: remove all existing output html files
//TODO: add tests for rw v0.1
//TODO: add test for file with few empty lines
//TODO: switch to using error #
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
		// {"no-cmd", "", true, ErrWrongCommand}, //running carpenter without options is now allowed
		{"wrong-cmd", "rnu", true, ""},
		{"version-cmd-ok", "version", false, ""},
		{"help-cmd-ok", "version", false, ""},
		{"wrong-global-flag", " -x run", true, "flag does not exist"},
		{"wrong-run-flag", "run %path/empty.gold -x", true, "flag does not exist"},

		{"missing-infile", "run ", true, ErrMissingInFile},
		//TODO: find a way to test this because now we are removing all existing output html
		//{"outfile-exists", "run %path/singletab.gold", true, "singletab.html: file exists"},

		{"no-binary", "check " + exe, true, "file does not start by a valid section separator"},
		{"no-SmallBinFile", "check %path/smallbinfile.gold", true, "file does not start by a valid section separator"},
		{"run-empty-File", "run %path/empty.gold", true, "file is empty"},
		{"check-empty-File", "check %path/empty.gold", true, "file is empty"},
		{"run-empty-File", "run %path/notab.gold -r", true, "empty table"},
		{"check-empty-File", "check %path/notab.gold", true, "empty table"},
		{"run-section-Sep-Only", "run %path/sectionsepsonly.gold -r", true, "empty table"},
		{"check-section-Sep-Only", "check %path/sectionsepsonly.gold", true, "empty table"},
		{"check-command-Only", "check %path/commandonly.gold", true, "empty table"},
		{"too-many-scections", "check %path/toomanyscections.gold", true, "incorrect number of sections"},

		{"check-syntax-errors", "check %path/wrong1tab.gold", true, "unknown command xmerge"},

		//All acceptable
		{"run-verysmalltab", "run %path/verysmalltab.gold -r", false, ""},
		{"run-oneCorrectTable", "run %path/singletab.gold -r", false, ""},
		{"run-typicalTable1", "run %path/TableT1A.gold -r", false, ""},
		//FIXME: fix this test
		//{"run-two-tables", "run %path/correct2tabs.gold -r", false, ""},

		//testing crashers
		//FIXME: fix this test
		//{"run-typicalTable1", "run %path/crash1.gold -r", true, "invalid coordinates"}, //passes check but fails on rendering
		// {"check-singleentryandcommandonly", "c", path + "singleentryandcommandonly.rw", false, "singleentryandcommandonly.rw]: file does not contain text"},
		// {"check-oneWrongTable", "c", path + "wrong1tab.rw", true, "wrong1tab.rw]: syntax error : unknown command xmerge"},
		// {"check-oneCorrectTable", "c", path + "singletab.rw", false, ""},
		// {"check-oneCorrectTableOnly", "c", path + "tabonly.rw", false, ""},
	}
	//TODO: add clean older output files and other setup code here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args = strings.Replace(tt.args, "%path", path, -1)
			tt.args = exe + tt.args //prepend exe name (args[0])
			t.Logf("command line: %s\n", tt.args)
			os.Args = strings.Fields(tt.args)
			//			fmt.Printf("len=%d ==>%v\n", len(cmdArgs), cmdArgs)
			err := RunApp() //skip the command name cmdArgs[1:]
			if (err != nil) != tt.wantErr {
				t.Errorf("RunApp() error = %v, wantErr %v\n", err, tt.wantErr)
				t.Logf("error: %v\n", err)
			}
			//a bit fragile as it depends on the exact error text not changing in the code
			if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("wrong error message, \nexpected= %q, \ngot= %q\n", tt.errMsg, err)
			}
			// if err != nil && !strings.Contains(errors.ErrorsToError(err).Error(), tt.errMsg) {
			// 	t.Errorf("wrong error message, expected= %s, got %s\n", tt.errMsg, errors.ErrorsToError(err))
			// }
		})
	}
}
