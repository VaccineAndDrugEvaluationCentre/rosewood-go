package rosewood_test

import (
	"fmt"
	"os"

	"github.com/drgo/rosewood"
)

func Example() {
	const filename = "test-files/correct2tabs.rw"
	ri := rosewood.NewInterpreter(rosewood.DefaultSettings())
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("failed to parse file %s\n", err)
		os.Exit(1)
	}
	defer file.Close()
	if err := ri.Run(file, os.Stdout); err != nil {
		fmt.Printf("error parsing file: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
