//  Copyright 2013 VDEC. All rights reserved.

// package carpenter is reference implementation of the Rosewood language
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/drgo/rosewood"
)

//initialized in the Makefile
var (
	Version string
	Build   string
)

var (
	verbose     bool
	help        bool
	outFileName string
	cssFileName string
)

func init() {
	flag.BoolVar(&help, "h", false, "")
	flag.StringVar(&outFileName, "o", "", "")
	flag.StringVar(&outFileName, "output", "", "")
	flag.StringVar(&cssFileName, "css", "", "")
	flag.BoolVar(&verbose, "v", false, "")
	flag.BoolVar(&verbose, "verbose", false, "")
}

func main() {
	flag.Usage = helpMessage
	flag.Parse()
	if help {
		usage(0)
	}
	//settings
	settings := rosewood.DefaultSettings()
	settings.Debug = verbose
	settings.StyleSheet = cssFileName

	//setup output
	var err error
	out := os.Stdout
	if outFileName != "" {
		if out, err = os.Create(outFileName); err != nil {
			crash(outFileName, err)
		}
		defer out.Close()
	}
	switch flag.NArg() {
	case 0:
		if info, _ := os.Stdin.Stat(); info.Size() > 0 { //input is being piped in
			if err := run("<stdin>", out, settings); err != nil {
				crash("<stdin>", err)
			}
		} else {
			usage(1)
		}
	default:
		for _, inFileName := range flag.Args() { //skip the command line name
			fmt.Println(inFileName)
			if err := run(inFileName, out, settings); err != nil {
				os.Exit(1)
			}
			fmt.Println(inFileName, " done.")
		}
	}
	os.Exit(0)
}

func run(inFileName string, out io.Writer, settings *rosewood.Settings) error {
	ri := rosewood.NewInterpreter(settings)
	Run := func(in io.Reader) error {
		err := ri.Run(in, out)
		if err != nil {
			fmt.Printf("error running file [%s]: %s\n", inFileName, err)
			eList := ri.Errors()
			//fmt.Println("eList:", eList)
			for _, e := range eList {
				fmt.Println(e)
			}
		}
		return err
	}
	switch inFileName {
	case "<stdin>":
		return Run(bufio.NewReader(os.Stdin))
	default:
		in, err := os.Open(inFileName)
		if err != nil {
			return fmt.Errorf("error opening input file %s: %s", inFileName, err)
		}
		defer in.Close()
		return Run(in)
	}
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

func crash(inFileName string, err error) {
	log.Fatalf("error running file [%s]: %s", inFileName, err)
}

// func runPipe(in io.Reader, out io.Writer) error {
// 	echo := func(s string) { //prints s to out followed by linefeed
// 		io.WriteString(out, s)
// 		io.WriteString(out, OSEOL)
// 	}
// 	p := rosewood.NewCommandParser(nil)
// 	cmdList, err := p.ParseCommands(in)
// 	if err != nil {
// 		echo(p.Errors(-1)) //show all errors
// 		return err
// 	}
// 	echo(cmdList[0].String())
// 	//p.Run(cmdList)
// 	if err != nil {
// 		echo(p.Errors(-1)) //show all errors
// 		return err
// 	}
// 	return nil
// }

const (
	othercolor = "\x1b[39m"
	redColor   = "\x1b[31m"
)

func newecho(w *io.Writer, s string, color string) {
	fmt.Printf("%s: %s", color, s)
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
