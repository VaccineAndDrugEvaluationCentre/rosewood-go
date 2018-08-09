package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/drgo/fileutils"
	"github.com/drgo/htmldocx"
	rosewood "github.com/drgo/rosewood/lib"
)

func cfgFailed(err error) error {
	return fmt.Errorf("failed to save configuration: %s", err)
}

//DoInit create a Rosewood job file
func DoInit(job *rosewood.Job) (string, error) {
	fileName := job.ConfigFileName()
	if fileName == "" {
		dir, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to create config file in dir %s: %s", dir, err)
		}
		fileName = filepath.Join(dir, ConfigFileBaseName)
	}
	return genConfigFile(job, fileName)
}

func genConfigFile(job *rosewood.Job, pFileName string) (fileName string, err error) {
	file, err := ioutil.TempFile("", "rw-temp")
	if err != nil {
		return "", cfgFailed(err)
	}
	defer func() {
		fileName = pFileName
		if fileName == "" {
			fileName = file.Name()
		}
		errC := fileutils.CloseAndRename(file, fileName, job.OverWriteOutputFile)
		if errC != nil {
			err = cfgFailed(errC)
			fileName = ""
		}
	}()
	if err := job.SaveToMDSon(file); err != nil {
		return "", cfgFailed(err)
	}
	//append htmldocx default document configuration
	doc, err := htmldocx.DefaultDocumentMDSon(job.HTMLFileNames, job.HTMLWorkDirName)
	if err != nil {
		return "", cfgFailed(err)
	}
	_, err = file.WriteString(doc)
	if err != nil {
		return "", cfgFailed(err)
	}
	return fileName, nil
}
