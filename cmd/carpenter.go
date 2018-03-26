// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

// carpenter is reference implementation of the Rosewood language
package main

import (
	"bufio"
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
		if err := Run(settings, flg.Args()); err != nil {
			return err
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

func setupCommandFlag(settings *rosewood.Settings) (flgSets []*flag.FlagSet, err error) {
	globals := NewCommand("", []Flag{
		{&settings.Debug, "debug", "d", 0},
	})
	cmdRun := NewCommand("run", []Flag{
		{&settings.StyleSheetName, "style", "s", ""},
		{&settings.OutputFileName, "output", "o", ""},
		{&settings.OverWriteOutputFile, "replace", "r", false},
		{&settings.SectionSeparator, "sep", "S", "+++"},
	})
	cmdCheck := NewCommand("check", []Flag{
		{&settings.SectionSeparator, "sep", "S", "+++"},
	})
	cmdHelp := NewCommand("help", []Flag{})
	cmdVersion := NewCommand("version", []Flag{})
	flgSets = append(flgSets, globals, cmdRun, cmdCheck, cmdHelp, cmdVersion)
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
			return err //fmt.Errorf(ErrRunningFile, "<stdin>", err)
		}
	case 1: //input=single file
		settings.InputFileName = args[0]
		if err := runSingle(settings); err != nil {
			return err // fmt.Errorf(ErrRunningFile, settings.InputFileName, err)
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
		in          io.ReadCloser
		out         io.WriteCloser
		err         error
		minFileSize = settings.SectionsPerTable * len(settings.SectionSeparator)
	)
	annotateError := func(err error) error {
		if err == nil {
			return nil
		}
		return fmt.Errorf("----------\nerror running file [%s]:\n%s", settings.InputFileName, err)
	}
	switch settings.InputFileName {
	case "<stdin>":
		in, _ = getInputFile("", minFileSize)
		//output either stdout if outFileName=="" or outFileName
		if !settings.CheckSyntaxOnly { //do not need an output for only checking syntax
			if out, err = getOutputFile(settings.OutputFileName, settings.OverWriteOutputFile); err != nil {
				return annotateError(err)
			}
			defer out.Close()
		}
	default: //single file
		if in, err = getInputFile(settings.InputFileName, minFileSize); err != nil {
			return annotateError(err)
		}
		defer in.Close()
		//output either outFileName or a new file =inFileName + "ext" if outFileName==""
		if !settings.CheckSyntaxOnly { //do not need an output
			if settings.OutputFileName == "" {
				settings.OutputFileName = fileutils.ReplaceFileExt(settings.InputFileName, "html")
			}
			if out, err = getOutputFile(settings.OutputFileName, settings.OverWriteOutputFile); err != nil {
				return annotateError(err)
			}
			defer out.Close()
		}
	}
	ri := rosewood.NewInterpreter(settings)
	return annotateError(runFile(ri, in, out))
}

func runFile(ri *rosewood.Interpreter, in io.Reader, out io.Writer) error {
	file, err := ri.Parse(bufio.NewReader(in), "")
	if err != nil || ri.Settings().CheckSyntaxOnly {
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
		out         io.WriteCloser
		err         error
		minFileSize = settings.SectionsPerTable * len(settings.SectionSeparator)
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
	for _, f := range inFileNames {
		in, err := getInputFile(f, minFileSize)
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

func getOutputFile(fileName string, overWrite bool) (*os.File, error) {
	if fileName == "" || fileName == "<stdout>" {
		return os.Stdout, nil
	}
	out, err := fileutils.CreateFile(fileName, overWrite)
	if err != nil {
		return nil, fmt.Errorf(ErrOpenOutFile, fileName, err)
	}
	return out, nil
}

func getInputFile(fileName string, minFileSize int) (*os.File, error) {
	if fileName == "" || fileName == "<stdin>" {
		return os.Stdin, nil
	}
	in, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	//this check here rather than in the interpreter because we need access to *File to rewind it
	//where as the interpreter uses io.Reader which does not have a stream
	if err = fileutils.CheckTextStream(in, minFileSize); err != nil {
		return nil, err
	}
	//rewind file stream
	_, err = in.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf(ErrOpenInFile, fileName, err)
	}
	return in, nil
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
