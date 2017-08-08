// carpenter
//  Reference implementation of RoseWood
//  Copyright Salah Mahmud, 2017

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/drgo/carpenter"
)

var (
	Version string
	Build   string
)

func main() {
	// fileName := flag.String("pattern", "", "Path to RoseWood file to parse")
	// flag.Parse()

	// if *pattern == "" {
	// 	fmt.Println("Pattern argument is missing.")
	// 	fmt.Println("Usage:")
	// 	flag.PrintDefaults()
	// 	return
	// }

	exitCode := 0
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		interactive()
	} else if info.Size() > 0 { //input is being piped in
		exitCode = run(bufio.NewReader(os.Stdin), os.Stdout)
	}

	// pt := pt.PlatinumSearcher{Out: os.Stdout, Err: os.Stderr}
	os.Exit(exitCode)
}

func interactive() {
	usage(true)
	in := bufio.NewScanner(os.Stdin)
	out := os.Stdout
	echo := func(s string) { //prints s to out followed by linefeed
		io.WriteString(out, s)
		io.WriteString(out, "\n")
	}
	p := carpenter.NewCommandParser(nil)
	for {
		fmt.Printf("\n>>")
		if !in.Scan() || strings.ToLower(in.Text()) == "q" {
			return
		}
		cmdList, err := p.ParseCommands(strings.NewReader(in.Text()))
		if err != nil {
			echo(p.Errors(0)) //show the first error only
			continue
		}
		echo(cmdList[0].String())
		err = p.Run(cmdList)
		if err != nil {
			echo(fmt.Sprintf("runtime error: %s", err)) //show the first error only
			continue
		}
	}
}

func run(in io.Reader, out io.Writer) int {
	echo := func(s string) { //prints s to out followed by linefeed
		io.WriteString(out, s)
		io.WriteString(out, "\n")
	}
	p := carpenter.NewCommandParser(nil)
	cmdList, err := p.ParseCommands(in)
	if err != nil {
		echo(p.Errors(-1)) //show all errors
		return 1
	}
	echo(cmdList[0].String())
	p.Run(cmdList)
	if err != nil {
		echo(p.Errors(-1)) //show all errors
		return 1
	}
	return 0
}

func runfromFile(fileName string) {
	// rs := carpenter.NewRwScript()
	// err := rs.ParseFile(fileName)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // log.Printf("%s \n", table)
	// log.Println("Start scanning...")
	// //	cmdList := ParseRosewoodScript(strings.NewReader(table.sections[TableControl].String()))
	// //	log.Printf("This is the entire command list: \n %+v", cmdList)
}

func usage(interactive bool) {
	fmt.Printf("Carpenter %s (%s)\nCopyRight VDEC 2017\n\n", Version, Build)
	if interactive {
		fmt.Printf("Enter a Rosewood command or press 'q' to exit\n")
	}
}

const (
	othercolor = "\x1b[39m"
	redColor   = "\x1b[31m"
)

func newecho(w *io.Writer, s string, color string) {
	fmt.Printf("%s: %s", color, s)
}
