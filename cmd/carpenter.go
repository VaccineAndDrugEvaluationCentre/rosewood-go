// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

// carpenter is reference implementation of the Rosewood language
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/drgo/errors"
	"github.com/drgo/fileutils"
	rosewood "github.com/drgo/rosewood/lib"
	"github.com/drgo/rosewood/lib/setter"
	_ "github.com/drgo/rosewood/renderers/html" //include needed renderers
)

type result struct {
	fileName string
	err      error
}

var (
	//Version holds the exe version initialized in the Makefile
	Version string
	//Build holds the exe build number initialized in the Makefile
	Build string
)

//TODO:
//clean up debug and warnings
//Debug=0 silent, 1=warnings only 2= verbose  3=internal debug info
// remove # from syntax error: syntax error line #39 col #23: expected row, col or an argument, found 1
// allow quoted argument in style command
// change version strings into constants eg v0.1 Version0_1
// remove blank line in top of rendered html
// add support to inlined-markdown.
// move all utilities to appropriate packages
// refresh vendor packages
// document new arguments
// add concurrent processing of files--> control access to stdout
// add graceful shutdown
// add reading of settings from a json file
// add support for settings in package types
// clean-up all tests.
// expose setter constant through Rosewood.
// add css into generated html files

func main() {
	if err := RunApp(); err != nil {
		crash(err)
	}
}

//RunApp has all program logic; easier to test than func main
func RunApp() error {
	if len(os.Args) == 1 { //only app name passed, exit
		return fmt.Errorf(ErrWrongCommand)
	}
	settings := rosewood.DefaultSettings()
	settings.ExecutableVersion = Version
	settings.InputFileName = "<stdin>" //default to stdin
	flgSets, _ := setupCommandFlag(settings)
	flg, err := ParseCommandLine(flgSets[0], flgSets[1:]...)
	if err != nil {
		return err
	}
	if settings.Debug == setter.DebugAll {
		fmt.Printf("current settings:\n%s\n", settings)
	}
	switch flg.Name() {
	case "check":
		settings.CheckSyntaxOnly = true
		fallthrough
	case "run":
		err = Run(settings, flg.Args())
	case "v1tov2":
		//settings.ConvertFromVersion = "v0.1"
		err = V1toV2(settings, flg.Args())
		if err != nil {
			err = fmt.Errorf("one or more errors occurred during file conversion %s", err)
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

//Run is the main work-horse function
func Run(settings *rosewood.Settings, args []string) error {
	if settings.Debug >= setter.DebugUpdates {
		fmt.Printf("Processing %d files\n", len(args))
	}
	switch len(args) {
	case 0: //input=stdin
		if info, _ := os.Stdin.Stat(); info.Size() == 0 {
			return fmt.Errorf(ErrMissingInFile)
		}
		settings.InputFileName = "<stdin>"
		if err := runSingle(settings); err != nil {
			return err
		}
	case 1: //input=single file
		settings.InputFileName = args[0]
		if err := runSingle(settings); err != nil {
			return err
		}
	default: //input= > 1 file
		switch settings.OutputFileName {
		case "": //one outputfile for each input file
			if settings.Debug >= setter.DebugUpdates {
				fmt.Println("in a multi-input multi-output mode\n")
			}
			errs := errors.NewErrorList()
			//TODO: convert to concurrent
			for _, settings.InputFileName = range args {
				fmt.Println(settings.InputFileName)
				if err := runSingle(settings); err != nil {
					errs.Add(errors.ErrorsToError(err))
				}
			}
			return fmt.Errorf(ErrRunningBatch, errors.ErrorsToError(errs))
		default: //all output goes into a single file
			if settings.Debug >= setter.DebugUpdates {
				fmt.Println("in a multi-input single-output mode\n")
			}
			if err := runMulti(args, settings); err != nil {
				return fmt.Errorf(ErrRunningBatch, errors.ErrorsToError(err))
			}
		}
	}
	return nil
}

//runSingle parses and render (if in run mode) a single input file
func runSingle(settings *rosewood.Settings) (err error) {
	var (
		in  *os.File
		out *os.File
	)
	if settings.Debug >= setter.DebugUpdates {
		fmt.Printf("Processing file %s\n", settings.InputFileName)
	}
	//define a function that saves the temp output file created below using settings.OutputFileName
	onOutputFileClose := func(outputFileName string) {
		if err == nil { //do not save file if runFile() below failed
			if closeErr := fileutils.CloseAndRename(out, outputFileName, settings.OverWriteOutputFile); closeErr != nil {
				err = annotateError(settings.InputFileName, closeErr)
			}
		}
	}
	iDesc := DefaultRwInputDescriptor(settings)
	switch settings.InputFileName {
	case "<stdin>":
		in, _ = getValidInputReader(iDesc.SetFileName("")) //setFileName called for consistency and clarity
		//output either stdout if outFileName=="" or outFileName
		if !settings.CheckSyntaxOnly { //do not need an output for only checking syntax
			//verify that settings.OutputFileName is never empty
			if out, err = getOutputFile(settings.OutputFileName, settings.OverWriteOutputFile); err != nil {
				return annotateError(settings.InputFileName, err)
			}
			defer onOutputFileClose(settings.OutputFileName)
		}
	default: //other file
		if in, err = getValidInputReader(iDesc.SetFileName(settings.InputFileName)); err != nil {
			return annotateError(settings.InputFileName, err)
		}
		defer in.Close()
		//output either settings.OutputFileName or a new file =inFileName + "ext" if settings.OutputFileName==""
		if !settings.CheckSyntaxOnly { //do not need an output
			outputFileName := settings.OutputFileName
			if outputFileName == "" {
				outputFileName = fileutils.ReplaceFileExt(settings.InputFileName, "html")
			}
			if out, err = getOutputFile(outputFileName, settings.OverWriteOutputFile); err != nil {
				return annotateError(settings.InputFileName, err)
			}
			defer onOutputFileClose(outputFileName)
		}
	}
	ri := rosewood.NewInterpreter(settings).SetScriptIdentifer(settings.InputFileName)
	err = annotateError(settings.InputFileName, runFile(ri, in, out))
	return
}

func runFile(ri *rosewood.Interpreter, in io.ReadSeeker, out io.Writer) error {
	file, err := ri.Parse(in, ri.ScriptIdentifer())
	if err != nil || ri.Setting().CheckSyntaxOnly {
		return ri.ReportError(err)
	}
	hr, err := rosewood.GetRendererByName("html")
	if err != nil {
		return err
	}
	return ri.ReportError(ri.Render(out, file, hr))
}

func runMulti(inFileNames []string, settings *rosewood.Settings) (err error) {
	var (
		in  *os.File
		out *os.File
	)
	errs := errors.NewErrorList()
	//define a function that saves the temp output file created below using settings.OutputFileName
	onOutputFileClose := func() {
		if errs == nil { //do not save file if runFile below failed
			if closeErr := fileutils.CloseAndRename(out, settings.OutputFileName, settings.OverWriteOutputFile); closeErr != nil {
				errs.Add(closeErr)
			}
		}
	}
	//open output file if needed
	if !settings.CheckSyntaxOnly {
		if out, err = getOutputFile(settings.OutputFileName, settings.OverWriteOutputFile); err != nil {
			return err
		}
		defer onOutputFileClose()
	}
	ri := rosewood.NewInterpreter(settings)
	for _, f := range inFileNames {
		in, err = getValidInputReader(DefaultRwInputDescriptor(settings).SetFileName(f))
		if err != nil {
			errs.Add(err)
			continue
		}
		if err = runFile(ri, in, out); err != nil {
			errs.Add(fmt.Errorf("error running file %s:\n%s", f, errors.ErrorsToError(err)))
		}
		in.Close()
	}
	return errs
}

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
