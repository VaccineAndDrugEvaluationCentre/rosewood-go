package main

import (
	"fmt"
	"os"
	"time"

	"github.com/drgo/errors"
	"github.com/drgo/fileutils"
	rosewood "github.com/drgo/rosewood/lib"
	"github.com/drgo/rosewood/lib/setter"
)

func V1toV2(job *Job) error {
	report := func(result result) {
		fmt.Printf("\n--------------\nconverting %s:", result.task.inputFile.Name)
		if result.error != nil {
			fmt.Printf("\nErrors: %v\n", result.error)
		} else {
			fmt.Printf("...Done\n")
			fmt.Printf("output file: %s\n", result.task.outputFile.Name)
		}
	}
	start := time.Now()
	if job.Settings.Debug >= setter.DebugUpdates {
		fmt.Printf("Processing %d files\n", len(job.InputFiles))
	}
	resChan := make(chan result)
	job.Settings.ConvertFromVersion = "v0.1"
	var err error
	//A counting semaphore to limit number of concurrently open files
	tokens := NewCountingSemaphore(job.Settings.MaxConcurrentWorkers)
	go func() {
		for _, inputFile := range job.InputFiles {
			//FIXME: constructfilename does not work for files with path
			task := job.GetTask(inputFile, NewFileDescriptor(fileutils.ConstructFileName(inputFile.Name, "rw", "", "-converted-v1-2-v2")))
			tokens.Reserve(1)
			go ConvertFile(*task, resChan)
		}
	}()
	for i := 0; i < len(job.InputFiles); i++ {
		res := <-resChan
		tokens.Free(1) //release a reserved worker
		if job.Settings.Debug >= rosewood.DebugUpdates || res.error != nil {
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

func ConvertFile(task task, resChan chan result) {
	var (
		in  *os.File
		out *os.File
		err error
	)
	//TODO: remove inputDescriptor?
	in, err = getValidInputReader(DefaultRwInputDescriptor(task.settings).SetFileName(task.inputFile.Name))
	if err != nil {
		resChan <- task.GetResult(err)
		return
	}
	defer in.Close()
	if _, err = os.Stat(task.outputFile.Name); err == nil && !task.settings.OverWriteOutputFile {
		resChan <- task.GetResult(fmt.Errorf("file already exists: %s", task.outputFile.Name))
		return
	}
	if out, err = getOutputWriter(task.outputFile.Name, task.settings.OverWriteOutputFile); err != nil {
		resChan <- task.GetResult(err)
		return
	}
	defer func() {
		if err == nil { //do not save file if runFile below failed
			if closeErr := fileutils.CloseAndRename(out, task.outputFile.Name, task.settings.OverWriteOutputFile); closeErr != nil {
				resChan <- task.GetResult(closeErr)
			}
		}
	}()
	if err = rosewood.ConvertToCurrentVersion(task.settings, in, out); err != nil {
		resChan <- task.GetResult(errors.ErrorsToError(err))
		return
	}
	resChan <- task.GetResult(err)
}
