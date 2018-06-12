package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func DoInit(job *Job) error {
	var path string
	if job.GetNameOfInputFile(0) == "" { //first input file is the name of new config file
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to create config file in dir %s: %s", dir, err)
		}
		path = filepath.Join(dir, ConfigFileBaseName)
	} else {
		path = job.InputFiles[0].Name
	}
	return SaveJob(job, path, true)
}
