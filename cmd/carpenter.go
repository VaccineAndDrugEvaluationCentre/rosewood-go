// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

// carpenter is reference implementation of the Rosewood language
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
	_ "github.com/drgo/rosewood/renderers/html" //include needed renderers
)

type result struct {
	inputFileName  string
	outputFileName string
	err            error
}

const (
	ConfigFileBaseName = "carpenter.json"
)

//TODO:
//clean up debug and warnings: Debug=0 silent, 1=warnings only 2= verbose  3=internal debug info
// allow quoted argument in style command
// change version strings into constants eg v0.1 Version0_1
// remove blank line in top of rendered html
// add support to inlined-markdown.
// move all utilities to appropriate packages
// refresh vendor packages
// document new arguments
// add concurrent processing of files--> check rosewood routines for shared memory
// add graceful shutdown https://golang.org/pkg/os/signal/
// add cleanup routine
// add reading of settings from a json file specified using config command--see setters.go
// add support for running jobs from a json file
// add support for settings in package types
// clean-up all tests.
// add css into generated html files
// use consistent errors types and constants eg NewError()
// add support for exporting to docx
// add support for processing subfolder if arg==./..

//Docs update
// check identifies v0.1 files as such but does not check them

func main() {
	if err := RunApp(); err != nil {
		crash(err)
	}
}

//TODO: add gracefull shutdown logic here along with a pointer to an optional cleanup function

//RunApp has all program logic; entry point for all tests
func RunApp() error {
	//fmt.Printf("main: len=%d ==>%v\n", len(cmdArgs), cmdArgs)
	if len(os.Args) == 1 { //only app name passed, exit
		return fmt.Errorf(ErrWrongCommand)
	}
	//load configuration from config file if one exists in current dir
	settings, err := GetValidSettings()
	if err != nil {
		return err
	}
	settings.ExecutableVersion = Version
	flgSets, _ := setupCommandFlag(settings)
	flg, err := ParseCommandLine(flgSets[0], flgSets[1:]...)
	if err != nil {
		return err
	}
	if settings.Debug == rosewood.DebugAll {
		fmt.Printf("current settings:\n%s\n", settings)
	}
	inputFileNames := flg.Args()
	switch flg.Name() {
	case "check":
		settings.CheckSyntaxOnly = true
		fallthrough
	case "run":
		//FIXME: this check is not working
		if err = Run(settings, inputFileNames); rosewood.Errors().IsParsingError(err) {
			err = fmt.Errorf("one or more errors occurred during file processing") //do not report again
		}
	case "v1tov2":
		if err = V1toV2(settings, inputFileNames); err != nil { //shadwing err, so check/run/convert errors will not returned
			fmt.Printf("one or more errors occurred during file processing")
		}
	case "init":
		if err = Init(settings, inputFileNames); err != nil {
			fmt.Printf("one or more errors occurred during configuration initialization")
		}
	case "version":
		fmt.Println(getVersion())
	case "help":
		helpMessage(flg.Args(), getVersion())
	default:
		helpMessage(nil, getVersion())
		return fmt.Errorf(ErrWrongCommand)
	}
	return err
}

//Run is the main work-horse function;
func Run(settings *rosewood.Settings, inputFileNames []string) error {
	//No arguments: check if the input is coming from stdin
	if len(inputFileNames) == 0 {
		if info, _ := os.Stdin.Stat(); info.Size() == 0 {
			return fmt.Errorf(ErrMissingInFile)
		}
		inputFileNames = append(inputFileNames, "") //empty argument signals stdin
	}
	var (
		err            error
		baseDir        string
		format         string
		outputFileName = strings.TrimSpace(settings.OutputFileName)
		start          time.Time
	)
	verb := "Check"
	if !settings.CheckSyntaxOnly {
		if format, err = GetValidFormat(inputFileNames, outputFileName); err != nil {
			return err
		}
		preserveWorkFiles := strings.TrimSpace(settings.WorkDirName) != "" || (format == "docx" && settings.PreserveWorkFiles) || format == "html"
		if baseDir, err = GetOutputBaseDir(settings.WorkDirName, preserveWorkFiles); err != nil {
			return err
		}
		if !preserveWorkFiles { //baseDir is temp, schedule removing it
			defer os.RemoveAll(baseDir)
		}
		verb = "Process"
	}

	if settings.Debug >= rosewood.DebugUpdates {
		fmt.Printf("%sing %d file(s)\n", verb, len(inputFileNames))
		start = time.Now()
	}

	processedFiles, err := runHTMLFiles(settings, inputFileNames, baseDir, format)
	if err != nil || len(processedFiles) == 0 {
		//	fmt.Printf("runHTMLFiles failed %s", err)
		return err
	}
	//fmt.Printf("format=%s \n %s \n", format, strings.Join(processedFiles, "|"))
	if !settings.CheckSyntaxOnly {
		switch {
		case format == "docx":
			docxOpts := htmldocx.DefaultOptions().SetDebug(settings.Debug)
			if err = htmldocx.Convert(processedFiles, outputFileName, docxOpts); err != nil {
				return fmt.Errorf("failed to convert to docx file: %s", err)
			}
		case format == "html":
			// fmt.Println("about to save file as " + outputFileName)
			// if outputFileName != "" { //we are here, so there must be only 1 file in processedFiles, rename it
			// 	if err = os.Rename(processedFiles[0], outputFileName); err != nil {
			// 		return fmt.Errorf("failed to rename file to %s: %s", outputFileName, err)
			// 	} else {
			// 		fmt.Println("file saved as " + outputFileName)
			// 	}
			// }
		default:
			return fmt.Errorf("unsupported format: %s", format) //should not happen
		}
	}
	if settings.Debug >= rosewood.DebugUpdates {
		fmt.Printf("%sed %d file(s) in %s\n", verb, len(inputFileNames), time.Since(start).String())
	}
	return err
}

//basedir always points to a valid dir to save output files
func runHTMLFiles(settings *rosewood.Settings, inputFileNames []string, baseDir, format string) ([]string, error) {
	report := func(res result) {
		if !settings.CheckSyntaxOnly {
			fmt.Printf("\n--------------\nprocessing %s:", res.inputFileName)
		} else {
			fmt.Printf("\n--------------\nchecking %s:", res.inputFileName)
		}
		if res.err != nil {
			fmt.Printf("\nErrors: %v\n", res.err)
		} else {
			fmt.Printf("...Done\n")
			if !settings.CheckSyntaxOnly {
				fmt.Printf("output file: %s\n", res.outputFileName)
			}
		}
	}
	//channel to communicate with
	resCh := make(chan result)
	//A counting semaphore to limit number of open files
	tokens := NewCountingSemaphore(settings.MaxConcurrentWorkers)
	go func() {
		outputFileName := strings.TrimSpace(settings.OutputFileName)
		tempOutputFileName := ""
		for _, inputFileName := range inputFileNames {
			switch {
			case outputFileName == "": //no output file assume html file with the same base name as inputfile
				tempOutputFileName = filepath.Join(baseDir,
					fileutils.ReplaceFileExt(filepath.Base(inputFileName), "html"))
			case format == "html": //this happens only if there was a single inputfile
				//AND outfilename with html ext, use the outputfilename
				if filepath.Dir(outputFileName) == "." { //no directory, use the baseDir
					tempOutputFileName = filepath.Join(baseDir, outputFileName)
				}
			case format == "docx": //create temp html files in the basedir
				tempOutputFileName = filepath.Join(baseDir,
					fileutils.ReplaceFileExt(filepath.Base(inputFileName), "html"))
			default:
				panic("unexpected branch in runHTMLFiles()") //should not happen
			}
			fmt.Println(outputFileName + "-->" + tempOutputFileName)
			tokens.Reserve(1)                                                 //reserve a worker
			go htmlRunner(settings, inputFileName, tempOutputFileName, resCh) //launch a runSingle worker for each file
		}
	}()
	var err error
	var processedFiles []string
	for i := 0; i < len(inputFileNames); i++ { //wait for workers to return one by one
		//fmt.Println("inside for loop")
		res := <-resCh
		//fmt.Printf("%+v", res)
		tokens.Free(1) //release a reserved worker
		if settings.Debug >= rosewood.DebugUpdates || res.err != nil {
			report(res)
		}
		if res.err == nil {
			processedFiles = append(processedFiles, res.outputFileName)
		}
		if err == nil {
			err = res.err
		}
	}
	return processedFiles, err
}

//htmlRunner parses and renders (if in run mode) a single input file into an HTML file
//all errors are returned through resChan channel; only one error per run
func htmlRunner(settings *rosewood.Settings, inputFileName, outputFileName string, resChan chan result) {
	var (
		in  *os.File
		out *os.File
		err error //not returned, used to decide whether the output file should be saved or not
	)
	if inputFileName == "" { //reading from stdin
		inputFileName = "<stdin>"
	}
	iDesc := DefaultRwInputDescriptor(settings)
	if in, err = getValidInputReader(iDesc.SetFileName(inputFileName)); err != nil {
		resChan <- result{inputFileName, outputFileName, err}
		return
	}
	defer in.Close()
	if !settings.CheckSyntaxOnly { //do not need an output
		//if the outputFileName already exists and OverWriteOutputFile is false, return an error
		if _, err := os.Stat(outputFileName); err == nil && !settings.OverWriteOutputFile {
			resChan <- result{inputFileName, outputFileName, fmt.Errorf("file already exists: %s", outputFileName)}
			return
		}
		//output either settings.OutputFileName or a new file =inFileName + "ext" if settings.OutputFileName==""
		if out, err = getOutputWriter(outputFileName, settings.OverWriteOutputFile); err != nil {
			resChan <- result{inputFileName, outputFileName, err}
			return
		}
		//define a function that saves the temp output file created below using settings.OutputFileName
		defer func() {
			if err == nil { //only save temp file if runFile() below succeeded
				resChan <- result{inputFileName, outputFileName, fileutils.CloseAndRename(out, outputFileName, settings.OverWriteOutputFile)}
				return
			}
		}()
	}
	ri := rosewood.NewInterpreter(settings).SetScriptIdentifer(inputFileName)
	err = runFile(ri, in, out)
	resChan <- result{inputFileName, outputFileName, err}
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

// func runMulti(settings *rosewood.Settings, inFileNames []string, resChan chan result) {
// 	var (
// 		in  *os.File
// 		out *os.File
// 		err error //not returned, used to decide whether the output file should be saved or not
// 	)
// 	outputFileName := settings.OutputFileName //can't be empty
// 	//define a function that saves the temp output file created below using settings.OutputFileName
// 	onOutputFileClose := func() {
// 		if err == nil { //only save temp file if runFile() below succeeded
// 			resChan <- result{outputFileName, fileutils.CloseAndRename(out, outputFileName, settings.OverWriteOutputFile)}
// 			return
// 		}
// 	}
// 	//open output file
// 	if out, err = getOutputFile(outputFileName, settings.OverWriteOutputFile); err != nil {
// 		resChan <- result{outputFileName, err}
// 		return
// 	}
// 	defer onOutputFileClose()

// 	ri := rosewood.NewInterpreter(settings)   //concurrent-safe
// 	for _, inputFileName := range inFileNames {
// 		in, err = getValidInputReader(DefaultRwInputDescriptor(settings).SetFileName(inputFileName))
// 		if err != nil {
// 			resChan <- result{inputFileName, err}
// 			return
// 		}
// 		defer in.Close()
// 		if err = runFile(ri, in, out); err != nil {
// 			resChan <- result{inputFileName, err}
// 			return
// 			//errs.Add(fmt.Errorf("error running file %s:\n%s", inputFileName, errors.ErrorsToError(err)))
// 		}
// 		in.Close()
// 	}
// }

// func interactive() {
// 	usage(-1)
// 	in := bufio.NewScanner(os.Stdin)
// 	out := os.Stdout
// 	echo := func(s string, status rosewood.ReportStatus) { //prints s to out followed by linefeed
// 		io.WriteString(out, s)
// 		io.WriteString(out, EOL)
// 	}
// 	settings := rosewood.DefaultSettings()
// 	settings.Report = echo
// 	rwi := rosewood.NewInterpreter(settings)
// 	for {
// 		trace.Printf("\n>>")
// 		if !in.Scan() || strings.ToLower(in.Text()) == "q" {
// 			return
// 		}
// 		err := rwi.Parse(strings.NewReader(in.Text()), "stdin")
// 		if err != nil {
// 			continue
// 		}
// 		//echo(cmdList[0].String())
// 		//		err = p.Run(cmdList)
// 		if err != nil {
// 			//			echo(fmt.Sprintf("runtime error: %s", err)) //show the first error only
// 			continue
// 		}
// 	}
// }

// if settings.OutputFileName != "" { // && len(args) > 1
//multiple files with one single output file
// if settings.Debug >= rosewood.DebugUpdates {
// 	fmt.Println("in a multi-input single-output mode")
// }
// if err = runMulti(settings, args); err != nil {
// 	return fmt.Errorf(ErrRunningBatch, errors.ErrorsToError(err))
// }
// } else {
//either one input file with
//this signals that we need to create one outputfile for each input file
// if settings.Debug >= rosewood.DebugUpdates {
// 	fmt.Println("in a multi-input multi-output mode")
// }
