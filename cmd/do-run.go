package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/drgo/fileutils"
	"github.com/drgo/htmldocx"
	rosewood "github.com/drgo/rosewood/lib"
)

//DoRun is the main work-horse function;
//WARNING: not safe to call concurrently
func DoRun(job *Job) error {
	var (
		err   error
		start time.Time
	)
	//if no input files: check if the input is coming from stdin
	if len(job.InputFiles) == 0 {
		if info, _ := os.Stdin.Stat(); info.Size() == 0 {
			return fmt.Errorf(ErrMissingInFile)
		}
		job.InputFiles = append(job.InputFiles, NewFileDescriptor("")) //empty argument signals stdin
	}
	if !job.Settings.CheckSyntaxOnly {
		if job.OutputFormat, err = GetValidFormat(job); err != nil {
			return err
		}
		job.WorkDirName = strings.TrimSpace(job.WorkDirName)
		//fmt.Printf("job.WorkDirName =%v \n", job.WorkDirName)
		job.PreserveWorkFiles = job.OutputFormat == "html" || job.WorkDirName != "" || (job.OutputFormat == "docx" && job.Settings.PreserveWorkFiles)
		//fmt.Printf("job.PreserveWorkFiles=%v, %v, %v \n", job.OutputFormat == "html", job.WorkDirName != "", (job.OutputFormat == "docx" && job.Settings.PreserveWorkFiles))
		if job.WorkDirName, err = GetOutputBaseDir(job.WorkDirName, job.PreserveWorkFiles); err != nil {
			return err
		}
		//Dangerous: could remove an entire folder if there is a bug somewhere
		if !job.PreserveWorkFiles && strings.Contains(job.WorkDirName, "rw-temp101") { //baseDir is temp, schedule removing it but only if we created it
			defer os.RemoveAll(job.WorkDirName)

		}
	}

	if job.Settings.Debug >= rosewood.DebugUpdates {
		fmt.Printf("%sing %d file(s) in work dir=%s\n", job.Command, len(job.InputFiles), job.WorkDirName)
		start = time.Now()
	}

	processedFiles, err := runHTMLFiles(job)
	if err != nil {
		return err
	}
	if len(processedFiles) == 0 {
		panic("unexpected failure in DoRun(): len(processedFiles) == 0")
	}
	if job.Settings.Debug >= rosewood.DebugAll {
		fmt.Println("inside DoRun() after returning from runhtmlfiles()")
		if err = fileutils.PrintFileStat(processedFiles[0]); err != nil {
			fmt.Printf("inside DoRun() failed to print stats of: %v\n", err)
		}
	}
	if !job.Settings.CheckSyntaxOnly {
		switch {
		case job.OutputFormat == "docx":
			docxOpts := htmldocx.DefaultOptions().SetDebug(job.Settings.Debug)
			if err = htmldocx.Convert(processedFiles, job.OutputFile.Name, docxOpts); err != nil {
				return err
			}
			if job.Settings.Debug >= rosewood.DebugUpdates {
				fmt.Printf("saved to docx %s\n", job.OutputFile.Name)
			}

		case job.OutputFormat == "html":
		default:
			return fmt.Errorf("unsupported format: %s", job.OutputFormat) //should not happen
		}
	}
	if job.Settings.Debug >= rosewood.DebugUpdates {
		fmt.Printf("%sed %d file(s) in %s\n", job.Command, len(job.InputFiles), time.Since(start).String())
	}
	return err
}

//basedir always points to a valid dir to save output files
func runHTMLFiles(job *Job) ([]string, error) {
	report := func(result result) {
		fmt.Printf("--------------\n%sing %s:", job.Command, result.task.inputFile.Name)
		if result.error != nil {
			fmt.Printf("\nErrors: %v\n", result.error)
		} else {
			fmt.Printf("...success\n")
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
			if job.Settings.Debug >= rosewood.DebugAll {
				fmt.Println("inside runhtmlfiles result processing loop")
				if err = fileutils.PrintFileStat(res.task.outputFile.Name); err != nil {
					fmt.Printf("inside runhtmlfiles failed to print stats of: %v\n", err)
				}
			}
			processedFiles = append(processedFiles, res.task.outputFile.Name)
		}
		if err == nil { //capture the first error
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
	//TODO: check stdin processing
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
	}
	ri := rosewood.NewInterpreter(task.settings).SetScriptIdentifer(task.inputFile.Name)
	if err = runFile(ri, in, out); err == nil && !task.settings.CheckSyntaxOnly {
		err = fileutils.CloseAndRename(out, task.outputFile.Name, task.settings.OverWriteOutputFile)
		// if err = fileutils.PrintFileStat(task.outputFile.Name); err != nil {
		// 	fmt.Printf("failed to print stats of: %v\n", err)
		// }
	}
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
