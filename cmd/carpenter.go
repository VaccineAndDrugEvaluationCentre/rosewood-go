// carpenter
//  Reference implementation of RoseWood
//  Copyright Salah Mahmud, 2017

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/drgo/rosewood"
)

var (
	Version string
	Build   string
	OSEOL   string
)

func init() {
	OSEOL = "\n"
	if runtime.GOOS == "windows" {
		OSEOL = "\r\n"
	}
}

func main() {
	verbose := flag.Bool("v", false, "verbose")
	help := flag.Bool("h", false, "prints this screen")
	flag.Usage = helpMessage
	flag.Parse()
	if *help {
		usage(0)
	}
	settings := rosewood.DefaultSettings()
	settings.Debug = *verbose
	var (
		input    io.ReadCloser
		out      io.WriteCloser
		filename string
		err      error
	)
	switch flag.NArg() {
	case 0:
		if info, _ := os.Stdin.Stat(); info.Size() > 0 { //input is being piped in
			input = ioutil.NopCloser(os.Stdin)
			filename = "<stdin>"
		} else {
			usage(1)
		}
	case 1:
		filename = flag.Arg(0)
		if input, err = os.Open(filename); err != nil {
			log.Fatalf("error processing input file %s: %s", filename, err)
		}
		defer input.Close()
	default:
		log.Fatalf("wrong parameters %s", strings.Join(flag.Args(), ","))
	}
	out = os.Stdout
	if err := run(bufio.NewReader(input), out, settings); err != nil {
		log.Fatalf("error processing input file %s: %s", filename, err)
	}
	os.Exit(0)
}

func run(in io.Reader, out io.Writer, settings *rosewood.Settings) error {
	ri := rosewood.NewInterpreter(settings)
	if err := ri.Run(in, out); err != nil {
		return err
	}

	return nil
}

func helpMessage() {
	fmt.Printf("Carpenter %s (%s)\nCopyRight VDEC 2017\n", Version, Build)
	io.WriteString(os.Stderr, "Usage: Carpenter <input Rosewood file> \n")
	flag.PrintDefaults()
}

func usage(exitCode int) {
	helpMessage()
	if exitCode > -1 {
		os.Exit(exitCode)
	}
}

func runPipe(in io.Reader, out io.Writer) error {
	echo := func(s string) { //prints s to out followed by linefeed
		io.WriteString(out, s)
		io.WriteString(out, OSEOL)
	}
	p := rosewood.NewCommandParser(nil)
	cmdList, err := p.ParseCommands(in)
	if err != nil {
		echo(p.Errors(-1)) //show all errors
		return err
	}
	echo(cmdList[0].String())
	//p.Run(cmdList)
	if err != nil {
		echo(p.Errors(-1)) //show all errors
		return err
	}
	return nil
}

const (
	othercolor = "\x1b[39m"
	redColor   = "\x1b[31m"
)

func newecho(w *io.Writer, s string, color string) {
	fmt.Printf("%s: %s", color, s)
}

//ParseFile takes path to a file containing RoseWood script and parses it possibly returning an error

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
// 		fmt.Printf("\n>>")
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

/*
// Package flg implements a getopt-compatible command line flag parser.
    package flg

    // Parse parses the command-line flags defined in package flag.
    // In contrast to flag.Parse, Parse imposes getopt semantics:
    //	- single letter flag names must be specified with a single dash: -x
    //	- longer names must be specified with a double dash: --long
    //	- the argument to a single-letter flag can follow it immediately:
    //	  -xfoo means -x foo when -x takes an argument.
    //	- multiple short flags can be combined: -xyz means -x -y -z,
    //	  when neither -x nor -y takes an argument.
    //	- name aliases can be introduced by calling Alias before Parse
    func Parse()

    // Alias records new as an alias for the flag named old.
    // Typically old is a long name and new is a single-letter name or vice versa.
    // For example, Alias("r", "recursive").
	func Alias(new, old string)

*/
