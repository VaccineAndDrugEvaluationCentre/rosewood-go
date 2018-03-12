//  Copyright 2017 Salah Mahmud. All rights reserved.

// package carpenter is reference implementation of the Rosewood language
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/drgo/fileutils"
	"github.com/drgo/rosewood"
)

var (
	//Version holds the exe version initialized in the Makefile
	Version string
	//Build holds the exe build number initialized in the Makefile
	Build string
	//Settings holds arguments/flags passed by the user
	settings *rosewood.Settings
)

func init() {
	settings = rosewood.DefaultSettings()
	settings.ExecutableVersion = Version
	settings.InputFileName = "<stdin>" //default to stdin

	flag.StringVar(&settings.OutputFileName, "o", "", "")
	flag.StringVar(&settings.OutputFileName, "output", "", "")

	flag.StringVar(&settings.StyleSheetName, "css", "", "")

	flag.BoolVar(&settings.Debug, "v", false, "")
	flag.BoolVar(&settings.Debug, "verbose", false, "")

	flag.StringVar(&settings.SectionSeparator, "s", "", "")
	flag.StringVar(&settings.SectionSeparator, "sep", "", "")

	flag.BoolVar(&settings.OverWriteOutputFile, "r", false, "")
	flag.BoolVar(&settings.OverWriteOutputFile, "replace", false, "")

	flag.BoolVar(&settings.CheckSyntaxOnly, "c", false, "")
	flag.BoolVar(&settings.CheckSyntaxOnly, "check", false, "")
}

func main() {
	flag.Usage = helpMessage
	flag.Parse()
	if settings.Debug {
		fmt.Printf("current settings:\n")
		fmt.Printf("%s\n", settings)
	}
	if flag.NArg() == 0 { //no arguments
		if info, _ := os.Stdin.Stat(); info.Size() == 0 { //no input is being piped in
			usage(1)
		}
	}
	if err := Run(settings, flag.Args()); err != nil {
		crash(err)
	}
	os.Exit(0)
}

//Run work-horse and main entry function
func Run(settings *rosewood.Settings, args []string) error {
	switch len(args) {
	case 0: //input=stdin
		if info, _ := os.Stdin.Stat(); info.Size() == 0 {
			return fmt.Errorf("nothing is piped in")
		}
		if err := runSingle(settings); err != nil {
			return fmt.Errorf("error running file [%s]: %s", "<stdin>", err)
		}
	case 1: //input=single file
		settings.InputFileName = args[0]
		if err := runSingle(settings); err != nil {
			return fmt.Errorf("error running file [%s]: %s", settings.InputFileName, err)
		}
	default: //input= > 1 file
		switch settings.OutputFileName {
		case "": //one outputfile for each input file
			var retErr error
			for _, settings.InputFileName = range args {
				if err := runSingle(settings); err != nil {
					retErr = err
				}
			}
			return fmt.Errorf("one or more errors running batch: %s", retErr)
		default: //all output goes into a single file
			if err := runMulti(args, settings); err != nil {
				return fmt.Errorf("one or more errors running batch: %s", err)
			}
		}
	}
	return nil
}

func runMulti(inFileNames []string, settings *rosewood.Settings) error {
	var (
		in          io.ReadCloser
		out         io.WriteCloser
		err, retErr error
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
	for _, f := range inFileNames {
		in, err = getInputFile(f, minFileSize)
		if err != nil {
			fmt.Printf("error opening file %s: %s\n", f, err)
			retErr = err
			continue
		}
		defer in.Close()
		err = runFile(ri, in, out)
		report(settings, err)
		if err != nil {
			retErr = err
		}
	}
	return retErr
}

func runSingle(settings *rosewood.Settings) error {
	var (
		in          io.ReadCloser
		out         io.WriteCloser
		err         error
		minFileSize = settings.SectionsPerTable * len(settings.SectionSeparator)
	)
	switch settings.InputFileName {
	case "<stdin>":
		in, _ = getInputFile("", minFileSize)
		//output either stdout if outFileName=="" or outFileName
		if !settings.CheckSyntaxOnly { //do not need an output
			if out, err = getOutputFile(settings.OutputFileName, settings.OverWriteOutputFile); err != nil {
				return err
			}
			defer out.Close()
		}
	default: //single file
		if in, err = getInputFile(settings.InputFileName, minFileSize); err != nil {
			return err
		}
		defer in.Close()
		//output either outFileName or a new file =inFileName + "ext" if outFileName==""
		if !settings.CheckSyntaxOnly { //do not need an output
			if settings.OutputFileName == "" {
				settings.OutputFileName = fileutils.ReplaceFileExt(settings.InputFileName, "html")
			}
			if out, err = getOutputFile(settings.OutputFileName, settings.OverWriteOutputFile); err != nil {
				return err
			}
			defer out.Close()
		}
	}
	ri := rosewood.NewInterpreter(settings)
	err = runFile(ri, in, out)
	return err
}

func runFile(ri *rosewood.Interpreter, in io.Reader, out io.Writer) error {
	file, err := ri.Parse(bufio.NewReader(in), "")
	if err != nil || ri.Settings().CheckSyntaxOnly {
		return err
	}
	err = ri.RenderTables(out, file.Tables(), rosewood.NewHTMLRenderer())
	report(ri.Settings(), err)
	return err
}

func report(settings *rosewood.Settings, err error) {
	if !settings.Debug {
		return
	}
	out := settings.OutputFileName
	if out == "" {
		out = "<stdout>"
	}
	if err == nil {
		fmt.Printf("File %s processed to %s\n", settings.InputFileName, out)
	} else {
		fmt.Printf("File %s failed\n", settings.InputFileName)
	}
}

func getOutputFile(fileName string, overWrite bool) (*os.File, error) {
	if fileName == "" || fileName == "<stdout>" {
		return os.Stdout, nil
	}
	out, err := fileutils.CreateFile(fileName, overWrite)
	if err != nil {
		return nil, fmt.Errorf("error opening output file: %s", err)
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
		return nil, fmt.Errorf("error opening input file %s: %s", fileName, err)
	}
	return in, nil
}

func helpMessage() {
	fmt.Fprintf(os.Stderr, versionMessage, Version, Build)
	fmt.Fprintln(os.Stderr, usageMessage)
}

func usage(exitCode int) {
	helpMessage()
	if exitCode > -1 {
		os.Exit(exitCode)
	}
}

const (
	redColor = "\x1b[31m"
)

func crash(err error) {
	fmt.Fprintf(os.Stderr, redColor+"%s\n", err)
	os.Exit(1)
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
