// carpenter
//  standrad implementation of RoseWood
//  Salah Mahmud
//  21 Aug 2016

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

// func init() {
// 	if cpu := runtime.NumCPU(); cpu == 1 {
// 		runtime.GOMAXPROCS(2)
// 	} else {
// 		runtime.GOMAXPROCS(cpu)
// 	}
// }

func main() {
	interactive()
	// pt := pt.PlatinumSearcher{Out: os.Stdout, Err: os.Stderr}
	// exitCode := pt.Run(os.Args[1:])
	// os.Exit(exitCode)
}

func interactive() {
	usage(true)
	in := bufio.NewScanner(os.Stdin)
	out := os.Stdin
	p := carpenter.NewCommandParser(nil)
	for {
		io.WriteString(out, "\n")
		fmt.Printf(">>")
		if !in.Scan() || strings.ToLower(in.Text()) == "q" {
			return
		}
		got, err := p.ParseCommands(strings.NewReader(in.Text()))
		if err != nil {
			io.WriteString(out, p.Errors(0)) //show the first error only
			continue
		}
		io.WriteString(out, got[0].String())
	}
}

func run() {
	// var fileName string
	// flag.StringVar(&fileName, "c", "", "Path to RoseWood file to parse")
	// flag.Parse()
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
