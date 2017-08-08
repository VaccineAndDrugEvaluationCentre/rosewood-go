package carpenter

import "fmt"

func ExampleNewRwInterpreter() {
	rwi := NewRwInterpreter()
	if err := rwi.ParseFile("test-files/fulltab1.txt"); err != nil {
		fmt.Printf("error parsing file: %s", err)
	}
	fmt.Println(rwi.SectionCount())
	if rwi.SectionCount() == 4 {
		fmt.Printf("%d\n", rwi.sections[0].offset)
		fmt.Printf("%d\n", rwi.sections[2].offset)
		fmt.Printf("%d\n", rwi.sections[0].LineCount())
		fmt.Printf("%d\n", rwi.sections[3].LineCount())
	}
	// Output:
	//4
	//2
	//19
	//1
	//19
}
