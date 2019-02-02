package main

import (
	"fmt"
	"os"
	"time"

	"github.com/drgo/core/errors"
	"github.com/drgo/core/files"
	rosewood "github.com/drgo/rosewood/lib"
)

//V1toV2 convert v 0.1 to v0.2 file
func V1toV2(job *rosewood.Job) error {
	report := func(result result) {
		fmt.Printf("\n--------------\nconverting %s:", result.task.inputFileName)
		if result.error != nil {
			fmt.Printf("\nErrors: %v\n", result.error)
		} else {
			fmt.Printf("...Done\n")
			fmt.Printf("output file: %s\n", result.task.outputFileName)
		}
	}
	start := time.Now()
	if job.RosewoodSettings.Debug >= rosewood.DebugUpdates {
		fmt.Printf("Processing %d files\n", len(job.RwFileNames))
	}
	resChan := make(chan result)
	job.RosewoodSettings.ConvertFromVersion = "v0.1"
	var err error
	//A counting semaphore to limit number of concurrently open files
	tokens := NewCountingSemaphore(job.RosewoodSettings.MaxConcurrentWorkers)
	go func() {
		for _, inputFile := range job.RwFileNames {
			//FIXME: constructfilename does not work for files with path
			task := getTask(inputFile, files.ConstructFileName(inputFile, "rw", "", "-converted-v1-2-v2"), job)
			tokens.Reserve(1)
			go convertFile(*task, resChan)
		}
	}()
	for i := 0; i < len(job.RwFileNames); i++ {
		res := <-resChan
		tokens.Free(1) //release a reserved worker
		if job.RosewoodSettings.Debug >= rosewood.DebugUpdates || res.error != nil {
			report(res)
		}
		err = res.error
	}
	if err == nil {
		fmt.Printf("Completed with no errors in %s\n", time.Since(start).String())
	} else {
		fmt.Printf("Completed with errors in %s\n", time.Since(start).String())
	}
	return err
}

//
func convertFile(task task, resChan chan result) {
	var (
		in  *os.File
		out *os.File
		err error
	)
	//TODO: remove inputDescriptor?
	in, err = getValidInputReader(DefaultRwInputDescriptor(task.settings).SetFileName(task.inputFileName))
	if err != nil {
		resChan <- task.getResult(err)
		return
	}
	defer in.Close()
	if _, err = os.Stat(task.outputFileName); err == nil && !task.OverWriteOutputFile {
		resChan <- task.getResult(fmt.Errorf("file already exists: %s", task.outputFileName))
		return
	}
	if out, err = getOutputWriter(task.outputFileName, task.OverWriteOutputFile); err != nil {
		resChan <- task.getResult(err)
		return
	}
	defer func() {
		if err == nil { //do not save file if runFile below failed
			if closeErr := files.CloseAndRename(out, task.outputFileName, task.OverWriteOutputFile); closeErr != nil {
				resChan <- task.getResult(closeErr)
			}
		}
	}()
	if err = rosewood.ConvertToCurrentVersion(task.settings, in, out); err != nil {
		resChan <- task.getResult(errors.ErrorsToError(err))
		return
	}
	resChan <- task.getResult(err)
}
