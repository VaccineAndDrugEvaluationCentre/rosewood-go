package main

import (
	"fmt"
	"os"
	"time"

	"github.com/drgo/fileutils"

	"github.com/drgo/errors"
	rosewood "github.com/drgo/rosewood/lib"
	"github.com/drgo/rosewood/lib/setter"
)

func V1toV2(settings *rosewood.Settings, inFileNames []string) error {
	start := time.Now()
	if settings.Debug >= setter.DebugUpdates {
		fmt.Printf("Processing %d files\n", len(inFileNames))
	}
	resChan := make(chan result)
	var err error
	for _, fileName := range inFileNames {
		go ConvertFile(settings, fileName, resChan)
	}
	for i := 0; i < len(inFileNames); i++ {
		res := <-resChan
		//		if settings.Debug >= setter.DebugUpdates || res.err != nil {
		fmt.Printf("processing " + res.inputFileName)
		if res.err != nil {
			fmt.Printf("\nErrors: %v\n", res.err)
		} else {
			fmt.Printf("...Done\n")
		}
		err = res.err
	}
	if err == nil {
		fmt.Printf("Completed with no errors in %s\n", time.Since(start).String())
	} else {
		fmt.Printf("Completed with errors in %s\n", time.Since(start).String())
	}
	return err
}

func ConvertFile(settings *rosewood.Settings, fileName string, resChan chan result) {
	var (
		in  *os.File
		out *os.File
	)
	//TODO: remove inputDescriptor?
	iDesc := DefaultRwInputDescriptor(settings)
	var err error
	in, err = getValidInputReader(iDesc.SetFileName(fileName))
	if err != nil {
		resChan <- result{fileName, "", err}
		return
	}
	defer in.Close()
	outputFileName := fileutils.ConstructFileName(fileName, "rw", "", "-converted-v1-2-v2")
	if out, err = getOutputWriter(outputFileName, settings.OverWriteOutputFile); err != nil {
		resChan <- result{fileName, outputFileName, err}
		return
	}
	defer func() {
		if err == nil { //do not save file if runFile below failed
			if closeErr := fileutils.CloseAndRename(out, outputFileName, settings.OverWriteOutputFile); closeErr != nil {
				resChan <- result{fileName, outputFileName, closeErr}
			}
		}
	}()
	if err = rosewood.ConvertToCurrentVersion(settings, in, out); err != nil {
		resChan <- result{fileName, outputFileName, errors.ErrorsToError(err)}
		return
	}
	resChan <- result{fileName, outputFileName, nil}
}

//TODO: fix use of annotate error
func annotateError(fileName string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("----------\nerror running file [%s]:\n%s", fileName, err)
}
