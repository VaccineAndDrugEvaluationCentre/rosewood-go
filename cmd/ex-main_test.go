package main

import (
	"fmt"
	"os"
)

func ExampleDoFromConfigFile() {
	os.Args = []string{"example"} //simulate passing no arguments
	err := RunApp()
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
