package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	for i, x := range flag.Args() {
		fmt.Printf("the argument #%d is %s\n", i, x)
	}
}
