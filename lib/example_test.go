// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package rosewood_test

import (
	"fmt"
	"os"

	"github.com/drgo/rosewood"
)

func Example() {
	const filename = "test-files/correct2tabs.rw"
	ri := rosewood.NewInterpreter(rosewood.DefaultSettings())
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("failed to open file %s\n", err)
		os.Exit(1)
	}
	defer f.Close()
	if err := ri.Run(f, os.Stdout); err != nil {
		fmt.Printf("error parsing file: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
