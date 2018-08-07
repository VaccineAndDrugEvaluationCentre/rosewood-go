package main

import (
	"fmt"
	"os"
	"path/filepath"

	rosewood "github.com/drgo/rosewood/lib"
)

//DoInit create a Rosewood job file
func DoInit(job *rosewood.Job) error {
	path := job.GetNameOfInputFile(0) //first input file is the name of new config file
	if path == "" {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to create config file in dir %s: %s", dir, err)
		}
		path = filepath.Join(dir, ConfigFileBaseName)
	}
	return job.SaveToMDSon(path, false)
}
