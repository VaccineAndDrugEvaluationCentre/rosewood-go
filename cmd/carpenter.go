// carpenter
//  Reference implementation of RoseWood
//  Copyright Salah Mahmud, 2017

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

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
		input      io.ReadCloser
		out        io.WriteCloser
		inFileName string
		err        error
	)
	switch flag.NArg() {
	case 2:
		outFileName := flag.Arg(1)
		if out, err = os.Open(outFileName); err != nil {
			log.Fatalf("error opening output file %s: %s", outFileName, err)
		}
		defer out.Close()
		fallthrough
	case 1:
		inFileName = flag.Arg(0)
		if input, err = os.Open(inFileName); err != nil {
			log.Fatalf("error opening input file %s: %s", inFileName, err)
		}
		defer input.Close()
	case 0:
		if info, _ := os.Stdin.Stat(); info.Size() > 0 { //input is being piped in
			inFileName = "<stdin>"
			input = os.Stdin
			out = os.Stdout
		} else {
			usage(1)
		}
	default:
		usage(1)
	}
	if err := run(bufio.NewReader(input), out, settings); err != nil {
		log.Fatalf("error running file %s: %s", inFileName, err)
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
	fmt.Fprintf(os.Stderr, "Carpenter %s (%s)\nCopyRight VDEC 2017\n", Version, Build)
	fmt.Fprintf(os.Stderr,
		`Usage: Carpenter <input Rosewood file> <output file>
	if only 1 file is specified, it will be used for input and output will be printed to standard output <stdout>.
	if both files are omitted, input will be read from standard input <stdin> and output printed to standard output <stdout>.%s
	`, "\n")
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
