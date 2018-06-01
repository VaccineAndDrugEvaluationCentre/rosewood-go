// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

// carpenter is reference implementation of the Rosewood language
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/drgo/errors"
	"github.com/drgo/fileutils"
	rosewood "github.com/drgo/rosewood/lib"
	_ "github.com/drgo/rosewood/renderers/html" //include needed renderers
)

var (
	//Version holds the exe version initialized in the Makefile
	Version string
	//Build holds the exe build number initialized in the Makefile
	Build string
)

//TODO:
// remove # from syntax error syntax error line #39 col #23: expected row, col or an argument, found 1
// allow quoted argument in style command
// change version strings into constants eg v0.1 Version0_1
// remove blank line in top of rendered html
// add support to inlined-markdown.
// use settings.Debug to control printing debug info; add no-warning default false
// move all utilities to appropriate packages
// refresh vendor packages

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
	if settings.Debug > 0 {
		fmt.Printf("current settings:\n%s\n", settings)
	}
	switch flg.Name() {
	case "check":
		settings.CheckSyntaxOnly = true
		fallthrough
	case "run":
		//TODO: remove if stmt
		if err := Run(settings, flg.Args()); err != nil {
			return err
		}
	case "v1tov2":
		settings.ConvertFromVersion = "v0.1"
		err = V1toV2(settings, flg.Args())
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

func setupCommandFlag(settings *rosewood.Settings) (flgSets []*flag.FlagSet, err error) {
	globals := NewCommand("", []Flag{
		{&settings.Debug, "debug", "d", 0},
		{&settings.OverWriteOutputFile, "replace", "r", false}, //ignored by some commands
	})
	cmdRun := NewCommand("run", []Flag{
		{&settings.ConvertOldVersions, "convert-old", "co", false},
		{&settings.DoNotInlineCSS, "no-inlined-css", "", false},
		{&settings.OutputFileName, "output", "o", ""},
		{&settings.SaveConvertedFile, "save-converted", "sc", false},
		{&settings.SectionSeparator, "sep", "S", "+++"},
		{&settings.StyleSheetName, "style", "s", ""},
		{&settings.ConvertFromVersion, "rosewood-version", "rv", ""},
	})
	cmdCheck := NewCommand("check", []Flag{
		{&settings.SectionSeparator, "sep", "S", "+++"},
	})
	cmdV1tov2 := NewCommand("v1tov2", []Flag{
		{&settings.ConvertFromVersion, "rosewood-version", "rv", ""},
	})
	cmdHelp := NewCommand("help", []Flag{})
	cmdVersion := NewCommand("version", []Flag{})
	flgSets = append(flgSets, globals, cmdRun, cmdCheck, cmdV1tov2, cmdHelp, cmdVersion)
	for _, fs := range flgSets {
		fs.Usage = func() {}    //disable internal usage function
		fs.SetOutput(devNull{}) //suppress output from package flag
	}
	return flgSets, nil
}

//Run work-horse and main entry function
func Run(settings *rosewood.Settings, args []string) error {
	if settings.Debug > 0 {
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
			errs := errors.NewErrorList()
			for _, settings.InputFileName = range args {
				if err := runSingle(settings); err != nil {
					errs.Add(errors.ErrorsToError(err))
				}
			}
			return fmt.Errorf(ErrRunningBatch, errors.ErrorsToError(errs))
		default: //all output goes into a single file
			if err := runMulti(args, settings); err != nil {
				return fmt.Errorf(ErrRunningBatch, errors.ErrorsToError(err))
			}
		}
	}
	return nil
}

//runSingle parses and render (if in run mode) a single input file
func runSingle(settings *rosewood.Settings) error {
	var (
		in  *os.File
		out *os.File
		err error
	)
	iDesc := DefaultRwInputDescriptor(settings)
	switch settings.InputFileName {
	case "<stdin>":
		in, _ = getValidInputReader(iDesc.SetFileName("")) //setFileName called for consistency and clarity
		//output either stdout if outFileName=="" or outFileName
		if !settings.CheckSyntaxOnly { //do not need an output for only checking syntax
			if out, err = getOutputFile(settings.OutputFileName, settings.OverWriteOutputFile); err != nil {
				return annotateError(settings.InputFileName, err)
			}
			defer out.Close()
		}
	default: //single file
		if in, err = getValidInputReader(iDesc.SetFileName(settings.InputFileName)); err != nil {
			return annotateError(settings.InputFileName, err)
		}
		defer in.Close()
		//output either outFileName or a new file =inFileName + "ext" if outFileName==""
		if !settings.CheckSyntaxOnly { //do not need an output
			if settings.OutputFileName == "" {
				settings.OutputFileName = fileutils.ReplaceFileExt(settings.InputFileName, "html")
			}
			if out, err = getOutputFile(settings.OutputFileName, settings.OverWriteOutputFile); err != nil {
				return annotateError(settings.InputFileName, err)
			}
			defer out.Close()
		}
	}
	ri := rosewood.NewInterpreter(settings).SetScriptIdentifer(settings.InputFileName)
	return annotateError(settings.InputFileName, runFile(ri, in, out))
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

func runMulti(inFileNames []string, settings *rosewood.Settings) error {
	var (
		out io.WriteCloser
		err error
	)
	//open output file if needed
	if !settings.CheckSyntaxOnly {
		if out, err = getOutputFile(settings.OutputFileName, settings.OverWriteOutputFile); err != nil {
			return err
		}
		defer out.Close()
	}
	ri := rosewood.NewInterpreter(settings)
	errs := errors.NewErrorList()
	iDesc := DefaultRwInputDescriptor(settings)
	for _, f := range inFileNames {
		in, err := getValidInputReader(iDesc.SetFileName(f))
		if err != nil {
			errs.Add(err)
			continue
		}
		defer in.Close()
		if err = runFile(ri, in, out); err != nil {
			errs.Add(fmt.Errorf("error running file %s:\n%s", f, errors.ErrorsToError(err)))
		}
	}
	return errs
}

func getVersion() string {
	return fmt.Sprintf(versionMessage, Version, Build, rosewood.Version)
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
