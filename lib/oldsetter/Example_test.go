package setter

import (
	"fmt"
	"os"
)

func ExampleLoadJob() {
	const scriptFileName = "/Users/salah/Dropbox/code/go/src/github.com/drgo/rosewood/lib/setter/carpenter.mdson"
	file, err := os.Open(scriptFileName)
	if err != nil {
		fmt.Println("error opening file", err.Error())
	}
	job, err := LoadJob(file)
	if err == nil {
		fmt.Printf("%v\n", job.OutputFileName)
		fmt.Printf("%v\n", job.RosewoodSettings.MaxConcurrentWorkers)
		fmt.Printf("%v\n", job.RosewoodSettings.ReportAllError)
		fmt.Println(job.InputFileNames[0])
	} else {
		fmt.Printf("error: %v\n", err)
	}
	// Output:
	// fromconfig.docx
	// 30
	// true
	// tab1old.html
}
