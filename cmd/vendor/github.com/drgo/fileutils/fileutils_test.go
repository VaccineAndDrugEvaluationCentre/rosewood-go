// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package fileutils

import (
	"testing"
)

// see function definition for rules
func Test_getUsableFilePath(t *testing.T) {
	type args struct {
		filePath string
		base     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		//if empty, return empty, regardless of base
		{"emptypath", args{"", ""}, ""},
		{"baseemptypath", args{"", "/test"}, "/test"},

		//if absolute and valid return filepath as is, regardless of base
		{"abspath", args{"/Users/salah/Dropbox/code/go/src/drgo/fileutils/fileutils_test.go", ""}, "/Users/salah/Dropbox/code/go/src/drgo/fileutils/fileutils_test.go"},
		{"baseabspath", args{"/Users/salah/Dropbox/code/go/src/drgo/fileutils/fileutils_test.go", "/test"}, "/Users/salah/Dropbox/code/go/src/drgo/fileutils/fileutils_test.go"},

		//relative from current dir and exists
		{"relpath1", args{"test/test.md", ""}, "test/test.md"},
		{"relpath2", args{"./test/import-stata/grant-list-tracker.tmpl", ""}, "./test/import-stata/grant-list-tracker.tmpl"},
		{"glob1", args{"./test/import-stata/*.xml", ""}, "./test/import-stata/*.xml"},

		{"baserelpath1", args{"/test/test.md", "/test"}, "/test/test.md"},
		{"baserelpath2", args{"../test/import-stata/grant-list-tracker.tmpl", "/test"}, "../test/import-stata/grant-list-tracker.tmpl"},
		{"baseglob1", args{"../test/import-stata/*.xml", "/test"}, "./test/import-stata/*.xml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUsableFilePath(tt.args.filePath, tt.args.base); got != tt.want {
				t.Errorf("GetUsableFilePath(): %v = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
