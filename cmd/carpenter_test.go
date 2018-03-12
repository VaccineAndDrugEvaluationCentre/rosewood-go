//  Copyright 2017 VDEC. All rights reserved.

//End-to-end testing of Carpenter main functions.
//Requires files stored currently in "../test-files/"
//Any Rosewood file that is not processed correctly should be added to the above folder and the fix must be tested below
package main

import (
	"strings"
	"testing"

	"github.com/drgo/rosewood"
)

func Test_run(t *testing.T) {
	const (
		path = "../test-files/" //all test files used below are stored here
	)
	tests := []struct {
		name    string   //unique test identifer
		flags   string   //semicolon separated list of valid carpenter flags
		args    []string //array of input filenames passed to run()
		wantErr bool     //true if run() must fail
		errMsg  string   //a fragment of the expected error description
	}{
		//default settings
		{"def-nothing", "", []string{}, true, "nothing is piped in"},
		{"def-emptyFile", "", []string{path + "empty.rw"}, true, "empty.rw]: stream is empty"},
		{"def-oneCorrectTable", "", []string{path + "singletab.rw"}, true, "singletab.html: file exists"},

		//replace enabled
		{"replace-nothing", "r", []string{}, true, "nothing is piped in"},
		{"replace-emptyFile", "r", []string{path + "empty.rw"}, true, "empty.rw]: stream is empty"},
		{"replace-sectionSepOnly", "r", []string{path + "sectionsepsonly.rw"}, true, "error parsing table in section # 2: empty table"},
		{"replace-verysmalltab", "r", []string{path + "verysmalltab.rw"}, false, ""},
		{"replace-oneCorrectTable", "r", []string{path + "singletab.rw"}, false, ""},
		{"replace-typicalTable1", "r", []string{path + "TableT1A.rw"}, false, ""},
		{"replace-fileOutput", "r;" + path + "newtabt1a.html", []string{path + "TableT1A.rw"}, false, ""},
		//syntax checking enabled
		{"check-nothing", "c", []string{}, true, "nothing is piped in"},
		{"check-emptyFile", "c", []string{path + "empty.rw"}, true, "empty.rw]: stream is empty"},
		{"check-SectionSepOnly", "c", []string{path + "sectionsepsonly.rw"}, true, "error parsing table in section # 2: empty table"},
		{"check-SmallBinFile", "c", []string{path + "smallbinfile"}, true, "smallbinfile]: file does not contain text"},
		{"check-commandOnly", "c", []string{path + "commandonly.rw"}, true, "error parsing table in section # 2: empty table"},
		{"check-singleentryandcommandonly", "c", []string{path + "singleentryandcommandonly.rw"}, false, "singleentryandcommandonly.rw]: file does not contain text"},
		{"check-oneWrongTable", "c", []string{path + "wrong1tab.rw"}, true, "wrong1tab.rw]: syntax error : unknown command xmerge"},
		{"check-oneCorrectTable", "c", []string{path + "singletab.rw"}, false, ""},
		{"check-oneCorrectTableOnly", "c", []string{path + "tabonly.rw"}, false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := setupSettings(tt.flags)
			err := Run(settings, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v\n", err, tt.wantErr)
				t.Logf("error: %v\n", err)
			}
			if err != nil && !strings.Contains(err.Error(), tt.errMsg) { //a bit fragile as it depends on the exact error text not changing in the code
				t.Errorf("wrong error message, expected= %s, got %s\n", tt.errMsg, err.Error())
			}
		})
	}
}

//setupSettings convert a semicolon separated string of valid Carpenter flags into a rosewood.Settings object
func setupSettings(flagText string) *rosewood.Settings {
	s := rosewood.DefaultSettings()
	flags := strings.Split(strings.ToLower(flagText), ";")
	if len(flags) == 0 {
		return s
	}
	for _, flag := range flags {
		switch flag {
		case "c", "check":
			s.CheckSyntaxOnly = true
		case "r", "replace":
			s.OverWriteOutputFile = true
		default: //anything else is output file name
			s.OutputFileName = flag
		}
	}
	return s
}
