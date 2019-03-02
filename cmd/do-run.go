package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/drgo/core/files"
	"github.com/drgo/htmldocx"
	rosewood "github.com/drgo/rosewood/lib"
)

//DoRun is the main work-horse function;
//WARNING: not thread safe
func DoRun(job *rosewood.Job) (err error) {
	//if no input files: return an error
	if len(job.RwFileNames) == 0 {
		return fmt.Errorf(ErrMissingInFile)
	}
	// prepare to preserve temp files if requested
	if !job.RosewoodSettings.CheckSyntaxOnly {
		if job.OutputFormat, err = GetValidFormat(job); err != nil {
			return err
		}
		job.WorkDirName = strings.TrimSpace(job.WorkDirName)
		job.PreserveWorkFiles = job.OutputFormat == "html" || job.WorkDirName != "" || (job.OutputFormat == "docx" && job.RosewoodSettings.PreserveWorkFiles)
		if job.WorkDirName, err = GetOutputBaseDir(job.WorkDirName, job.PreserveWorkFiles); err != nil {
			return err
		}
		if !job.PreserveWorkFiles && strings.Contains(job.WorkDirName, "rw-temp101") { //baseDir is temp, schedule removing it but only if we created it
			//Dangerous: could remove an entire folder if there is a bug somewhere
			defer os.RemoveAll(job.WorkDirName)
		}
	}

	ux.Info("processing ", len(job.RwFileNames), "file(s) in work dir=", job.WorkDirName)
	start := time.Now()

	processedFiles, err := runHTMLFiles(job)
	if err != nil {
		return err
	}
	if len(processedFiles) == 0 {
		return fmt.Errorf("unexpected failure in DoRun(): len(processedFiles) == 0")
	}
	ux.Log("inside DoRun() after returning from runhtmlfiles()")

	if !job.RosewoodSettings.CheckSyntaxOnly {
		switch {
		case job.OutputFormat == "docx":
			if err = outputAsDocx(processedFiles, job); err != nil {
				return err
			}
		case job.OutputFormat == "html":
		default:
			return fmt.Errorf("unsupported format: %s", job.OutputFormat) //should not happen
		}
	}
	ux.Info("processed ", len(job.RwFileNames), " file(s) in ", time.Since(start).String())
	// FIXME: enable code but allow for not saving config file automatically
	// if job.ConfigFileName != "" {
	// 	configFilename, err := DoInit(job)
	// 	if err == nil {
	// 		ux.Info("configuration saved as '%s'\n", configFilename)
	// 	}
	// }
	return err
}

//basedir always points to a valid dir to save output files
func runHTMLFiles(job *rosewood.Job) ([]string, error) {
	report := func(result result) {
		if job.RosewoodSettings.Debug < rosewood.DebugUpdates {
			return
		}
		fmt.Printf("--------------\nprocessing %s:", result.task.inputFileName)
		if result.error == nil {
			fmt.Printf("...success\n")
			if !job.RosewoodSettings.CheckSyntaxOnly {
				fmt.Printf("output file: %s\n", result.task.outputFileName)
			}
			return
		}
		fmt.Println("")
	}
	//channel to communicate with; passing result instead of &result because it is small (2 pointers) and to avoid
	//additional memory allocations since &result escapes to the heap and requires GC action
	resCh := make(chan result)
	//A counting semaphore to limit number of concurrently open files
	tokens := NewCountingSemaphore(job.RosewoodSettings.MaxConcurrentWorkers)
	go func() {
		for _, inputFile := range job.RwFileNames {
			task := getTask(inputFile, filepath.Join(job.WorkDirName,
				files.ReplaceFileExt(filepath.Base(inputFile), "html")), job) //assume outputfile is html file with path= job.WorkdirName+ same base name as inputfile
			switch {
			case job.OutputFileName == "": //no output file assume html file with the same base name as inputfile
			case job.OutputFormat == "docx": //create temp html files in the basedir
			//in both cases, above assumed outputfile name is good
			case job.OutputFormat == "html": //this happens only if there was a single inputfile
				//AND outfilename with html ext. If so, use the outputfilename
				if filepath.Dir(job.OutputFileName) == "." { //no directory, use the work dir
					task.outputFileName = filepath.Join(job.WorkDirName, job.OutputFileName)
				} else {
					task.outputFileName = job.OutputFileName
				}
			default:
				panic("unexpected branch in runHTMLFiles()") //should not happen
			}
			tokens.Reserve(1)           //reserve a worker
			go htmlRunner(*task, resCh) //launch a runSingle worker for each task
		}
	}()
	var err error
	var processedFiles []string
	for i := 0; i < len(job.RwFileNames); i++ { //wait for workers to return one by one
		res := <-resCh
		tokens.Free(1) //release a reserved worker
		report(res)
		processedFiles = append(processedFiles, res.task.outputFileName)
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
	iDesc := DefaultRwInputDescriptor(task.settings).SetFileName(task.inputFileName)
	if in, err = getValidInputReader(iDesc); err != nil {
		resChan <- task.getResult(err) //signal end of task run
		return
	}
	defer in.Close()
	if !task.settings.CheckSyntaxOnly { //do not need an output
		//if the outputFileName already exists and OverWriteOutputFile is false, return an error
		if _, err := os.Stat(task.outputFileName); err == nil && !task.OverWriteOutputFile {
			resChan <- task.getResult(fmt.Errorf("file already exists: %s", task.outputFileName))
			return
		}
		//get a temp writer
		if out, err = getOutputWriter(task.outputFileName, task.OverWriteOutputFile); err != nil {
			resChan <- task.getResult(err)
			return
		}
	}
	ri := rosewood.NewInterpreter(task.settings).SetScriptIdentifer(task.inputFileName)
	if err = runFile(ri, in, out); err == nil && !task.settings.CheckSyntaxOnly {
		err = files.CloseAndRename(out, task.outputFileName, task.OverWriteOutputFile)
	}
	resChan <- task.getResult(err)
}

func runFile(ri *rosewood.Interpreter, in io.ReadSeeker, out io.Writer) error {
	file, err := ri.Parse(in, ri.ScriptIdentifer())
	if err != nil || ri.Settings().CheckSyntaxOnly {
		return ri.ReportError(err)
	}
	hr, err := rosewood.GetRendererByName("html") //TODO: get from settings
	if err != nil {
		return err
	}
	return ri.ReportError(ri.Render(out, file, hr))
}

//outputAsDocx caution: modifies job, not thread-safe
func outputAsDocx(processedFiles []string, job *rosewood.Job) error {
	var err error
	configFileName := job.ConfigFileName
	ux.Log("current configFileName", configFileName)
	job.HTMLFileNames = processedFiles
	if configFileName == "" { //from commandline, create temp config file
		if configFileName, err = genConfigFile(job, ""); err != nil {
			return err
		}
		ux.Log("created htmldocx config file ", configFileName)
		defer os.Remove(configFileName)
	}
	docxOpts := htmldocx.DefaultOptions().SetDebug(job.RosewoodSettings.Debug)
	if err := htmldocx.MakeDocxFromMdsonFile(configFileName, docxOpts); err != nil {
		return err
	}
	ux.Log("saved to docx", job.OutputFileName)
	return nil
}
