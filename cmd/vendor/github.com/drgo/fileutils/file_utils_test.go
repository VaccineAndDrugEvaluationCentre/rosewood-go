// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package fileutils

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCreateFile(t *testing.T) {
	tests := []struct {
		name      string
		fileName  string
		overWrite bool
		wantErr   bool
	}{
		{"already exists", createTempFile(), false, true},
		{"already exists-overwrite", createTempFile(), true, false},
		{"does not exist", "newfile123212", false, false},
		{"dir does not exist", "newdir123212/newfile", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(tt.fileName) // clean up
			gotOut, err := CreateFile(tt.fileName, tt.overWrite)
			defer gotOut.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func createTempFile() string {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}
	return tmpfile.Name()
}
