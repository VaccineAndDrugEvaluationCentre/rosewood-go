package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/drgo/fileutils"
	rosewood "github.com/drgo/rosewood/lib"
)

//basedir always points to a valid dir to save output files
func runHTMLFiles(job *Job) ([]string, error) {
	report := func(result result) {
		fmt.Printf("\n--------------\n%sing %s:", job.Command, result.task.inputFile.Name)
		if result.error != nil {
			fmt.Printf("\nErrors: %v\n", result.error)
		} else {
			fmt.Printf("...Done\n")
			if !job.Settings.CheckSyntaxOnly {
				fmt.Printf("output file: %s\n", result.task.outputFile.Name)
			}
		}
	}
	//channel to communicate with; passing result instead of &result because it is small (2 pointers) and to avoid
	//additional memory allocations since &result escapes to the heap and requires GC action
	resCh := make(chan result)
	//A counting semaphore to limit number of concurrently open files
	tokens := NewCountingSemaphore(job.Settings.MaxConcurrentWorkers)
	go func() {
		for _, inputFile := range job.InputFiles {
			task := job.GetTask(inputFile, NewFileDescriptor(filepath.Join(job.WorkDirName,
				fileutils.ReplaceFileExt(filepath.Base(inputFile.Name), "html")))) //assume outputfile is html file with path= job.WorkdirName+ same base name as inputfile
			switch {
			case job.OutputFile.Name == "": //no output file assume html file with the same base name as inputfile
			case job.OutputFormat == "docx": //create temp html files in the basedir
			//in both cases, above assumed outputfile name is good
			case job.OutputFormat == "html": //this happens only if there was a single inputfile
				//AND outfilename with html ext. If so, use the outputfilename
				if filepath.Dir(job.OutputFile.Name) == "." { //no directory, use the work dir
					task.outputFile.Name = filepath.Join(job.WorkDirName, job.OutputFile.Name)
				} else {
					task.outputFile.Name = job.OutputFile.Name
				}
			default:
				panic("unexpected branch in runHTMLFiles()") //should not happen
			}
			//fmt.Printf("Task: %+v\n", task)
			tokens.Reserve(1)           //reserve a worker
			go htmlRunner(*task, resCh) //launch a runSingle worker for each task
		}
	}()
	var err error
	var processedFiles []string
	for i := 0; i < len(job.InputFiles); i++ { //wait for workers to return one by one
		//fmt.Println("inside for loop")
		res := <-resCh
		//fmt.Printf("%+v", res)
		tokens.Free(1) //release a reserved worker
		if job.Settings.Debug >= rosewood.DebugUpdates || res.error != nil {
			report(res)
		}
		if res.error == nil {
			processedFiles = append(processedFiles, res.task.outputFile.Name)
		}
		if err == nil {
			err = res.error
		}
	}
	return processedFiles, err
}

//htmlRunner parses and renders (if in run mode) a single input file into an HTML file
//all errors are returned through resChan channel; only one error per run
func htmlRunner(task task, resChan chan result) {
	var (
		in  *os.File
		out *os.File
		err error //not returned, used to decide whether the output file should be saved or not
	)
	// if inputFileName == "" { //reading from stdin
	// 	inputFileName = "<stdin>"
	// }

	iDesc := DefaultRwInputDescriptor(task.settings).SetFileName(task.inputFile.Name)
	if in, err = getValidInputReader(iDesc); err != nil {
		resChan <- task.GetResult(err) //signal end of task run
		return
	}
	defer in.Close()
	if !task.settings.CheckSyntaxOnly { //do not need an output
		//if the outputFileName already exists and OverWriteOutputFile is false, return an error
		if _, err := os.Stat(task.outputFile.Name); err == nil && !task.settings.OverWriteOutputFile {
			resChan <- task.GetResult(fmt.Errorf("file already exists: %s", task.outputFile.Name))
			return
		}
		//get a temp writer
		if out, err = getOutputWriter(task.outputFile.Name, task.settings.OverWriteOutputFile); err != nil {
			resChan <- task.GetResult(err)
			return
		}
		//define a function that saves the temp output file created below using settings.OutputFileName
		defer func() {
			if err == nil { //only save temp file if runFile() below succeeded
				resChan <- task.GetResult(fileutils.CloseAndRename(out, task.outputFile.Name, task.settings.OverWriteOutputFile))
				return
			}
		}()
	}
	ri := rosewood.NewInterpreter(task.settings).SetScriptIdentifer(task.inputFile.Name)
	err = runFile(ri, in, out)
	resChan <- task.GetResult(err)
}

func runFile(ri *rosewood.Interpreter, in io.ReadSeeker, out io.Writer) error {
	file, err := ri.Parse(in, ri.ScriptIdentifer())
	if err != nil || ri.Setting().CheckSyntaxOnly {
		return ri.ReportError(err)
	}
	hr, err := rosewood.GetRendererByName("html") //TODO: get from settings
	if err != nil {
		return err
	}
	return ri.ReportError(ri.Render(out, file, hr))
}
