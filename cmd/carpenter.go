// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

// carpenter is reference implementation of the Rosewood language
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/drgo/fileutils"
	"github.com/drgo/htmldocx"
	rosewood "github.com/drgo/rosewood/lib"
	_ "github.com/drgo/rosewood/renderers/html" //include needed renderers
)

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
//WARNING: not safe to call concurrently; this is the only function allowed to change the job configuration
//WARNING: beyond this function, changing settings or job fields is not permitted
func RunApp() error {
	if len(os.Args) == 1 { //no command line arguments
		return DoFromConfigFile()
	}
	job, err := LoadConfigFromCommandLine()
	if err != nil {
		return err
	}
	if job.Settings.Debug == rosewood.DebugAll {
		fmt.Printf("current settings:\n%s\n", job)
	}
	switch job.Command { //TODO: check command is case insensitive
	case "do":
		if err = DoFromConfigFile(); err != nil {
			return err
		}
	case "check":
		job.Settings.CheckSyntaxOnly = true //TODO: rationalize; not needed any more here
		fallthrough
	case "run":
		//FIXME: this check is not working
		if err = DoRun(job); rosewood.Errors().IsParsingError(err) {
			err = fmt.Errorf("one or more errors occurred during file processing") //do not report again
		}
	case "v1tov2":
		if err = V1toV2(job); err != nil { //shadwing err, so check/run/convert errors will not returned
			fmt.Printf("one or more errors occurred during file processing")
		}
	case "init": //FIXME
		if err = DoInit(job); err != nil {
			fmt.Printf("one or more errors occurred during configuration initialization")
		}
	case "version":
		fmt.Println(getVersion())
	case "help": //FIXME:
		//helpMessage(job.InputFiles, getVersion())
	default:
		helpMessage(nil, getVersion())
		return fmt.Errorf(ErrWrongCommand)
	}
	return err
}

func DoFromConfigFile() error {
	var (
		configFileName string
		err            error
	)
	if len(os.Args) == 1 { //only app name passed, use ConfigFileBaseName in current folder
		if configFileName, err = fileutils.GetFullPath(ConfigFileBaseName); err != nil {
			return err
		}
	} else {
		//2 or more arguments, is it app name + do + a json file
		if strings.TrimSpace(strings.ToLower(os.Args[1])) != "do" {
			return fmt.Errorf("invalid command %s", os.Args[1])
		}
		if len(os.Args) < 3 {
			return fmt.Errorf("must specify a json configuration file")
		}
		configFileName = os.Args[2]
		if strings.ToLower(filepath.Ext(configFileName)) != ".json" {
			return fmt.Errorf("invalid config file name [%s] passed, ext must be json", configFileName)
		}
		if configFileName, err = fileutils.GetFullPath(configFileName); err != nil {
			return err
		}
	}
	//load configuration from config file
	configBuf, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return fmt.Errorf("failed to load configuration file: %v", err)
	}
	job := DefaultJob(rosewood.DefaultSettings())
	if err = json.Unmarshal(configBuf, job); err != nil {
		return fmt.Errorf("failed to parse configuration file %s: %v", configFileName, err)
	}
	if job.Settings.Debug >= rosewood.DebugUpdates {
		fmt.Println("configuration loaded from " + configFileName)
	}
	job.FileName = configFileName
	if job.Settings.Debug >= rosewood.DebugAll {
		fmt.Printf("current configuration: \n %s\n", job)
	}
	if err = DoRun(job); rosewood.Errors().IsParsingError(err) {
		err = fmt.Errorf("one or more errors occurred during file processing") //do not report again
	}
	return err
}

func LoadConfigFromCommandLine() (*Job, error) {
	job := DefaultJob(rosewood.DefaultSettings()) //TODO: ensure all defaults are reasonable
	flgSets, _ := setupCommandFlag(job)
	flg, err := ParseCommandLine(flgSets[0], flgSets[1:]...)
	if err != nil {
		return nil, err
	}
	job.Command = flg.Name()
	//TODO: validate command line inputs; use in RunfromConfigFile too?!
	for _, fileName := range flg.Args() {
		job.InputFiles = append(job.InputFiles, NewFileDescriptor(fileName))
	}
	return job, nil
}

//Run is the main work-horse function;
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
		fmt.Printf("job.WorkDirName =%v \n", job.WorkDirName)
		job.PreserveWorkFiles = job.OutputFormat == "html" || job.WorkDirName != "" || (job.OutputFormat == "docx" && job.Settings.PreserveWorkFiles)
		fmt.Printf("job.PreserveWorkFiles=%v, %v, %v \n", job.OutputFormat == "html", job.WorkDirName != "", (job.OutputFormat == "docx" && job.Settings.PreserveWorkFiles))
		if job.WorkDirName, err = GetOutputBaseDir(job.WorkDirName, job.PreserveWorkFiles); err != nil {
			return err
		}
		if !job.PreserveWorkFiles { //baseDir is temp, schedule removing it
			defer os.RemoveAll(job.WorkDirName)
		}
	}

	if job.Settings.Debug >= rosewood.DebugUpdates {
		fmt.Printf("%sing %d file(s) in work dir=%s\n", job.Command, len(job.InputFiles), job.WorkDirName)
		start = time.Now()
	}

	processedFiles, err := runHTMLFiles(job)
	if err != nil || len(processedFiles) == 0 {
		return err
	}
	fmt.Printf("format=%s \n %s \n", job.OutputFormat, strings.Join(processedFiles, "|\n"))
	if !job.Settings.CheckSyntaxOnly {
		switch {
		case job.OutputFormat == "docx":
			docxOpts := htmldocx.DefaultOptions().SetDebug(job.Settings.Debug)
			if err = htmldocx.Convert(processedFiles, job.OutputFile.Name, docxOpts); err != nil {
				return fmt.Errorf("failed to convert to docx file: %s", err)
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

// func runMulti(settings *rosewood.Settings, inFileNames []string, resChan chan task) {
// 	var (
// 		in  *os.File
// 		out *os.File
// 		err error //not returned, used to decide whether the output file should be saved or not
// 	)
// 	outputFileName := settings.OutputFileName //can't be empty
// 	//define a function that saves the temp output file created below using settings.OutputFileName
// 	onOutputFileClose := func() {
// 		if err == nil { //only save temp file if runFile() below succeeded
// 			resChan <- task{outputFileName, fileutils.CloseAndRename(out, outputFileName, settings.OverWriteOutputFile)}
// 			return
// 		}
// 	}
// 	//open output file
// 	if out, err = getOutputFile(outputFileName, settings.OverWriteOutputFile); err != nil {
// 		resChan <- task{outputFileName, err}
// 		return
// 	}
// 	defer onOutputFileClose()

// 	ri := rosewood.NewInterpreter(settings)   //concurrent-safe
// 	for _, inputFileName := range inFileNames {
// 		in, err = getValidInputReader(DefaultRwInputDescriptor(settings).SetFileName(inputFileName))
// 		if err != nil {
// 			resChan <- task{inputFileName, err}
// 			return
// 		}
// 		defer in.Close()
// 		if err = runFile(ri, in, out); err != nil {
// 			resChan <- task{inputFileName, err}
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
// if job.Settings.Debug >= rosewood.DebugUpdates {
// 	fmt.Println("in a multi-input single-output mode")
// }
// if err = runMulti(settings, args); err != nil {
// 	return fmt.Errorf(ErrRunningBatch, errors.ErrorsToError(err))
// }
// } else {
//either one input file with
//this signals that we need to create one outputfile for each input file
// if job.Settings.Debug >= rosewood.DebugUpdates {
// 	fmt.Println("in a multi-input multi-output mode")
// }
