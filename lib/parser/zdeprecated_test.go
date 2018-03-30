// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package parser

import (
	"fmt"
	"testing"

	"./types"
)

func Test_createGridTable(t *testing.T) {
	tests := []struct {
		name       string
		tab        string
		source     string
		wantMrlist []types.Range
		wantErr    bool
	}{
		// {"", "row1|\nr2c1|r2c2|r2c3|\nr3c1|r3c2|r3c3|r3c4|\nr4c1|r4c2|r4c3|\n",
		// 	"", nil, false}, //no merge commands
		// {"", "row1|\nr2c1|r2c2|r2c3|\nr3c1|r3c2|r3c3|r3c4|\nr4c1|r4c2|r4c3|\n",
		// 	`merge	row 4 col 1:2
		// 	merge	row 1
		// 	merge	row 2:2 col 2:3
		// 	merge	row 2:3 col 1:2
		// 	merge row 1:2 col 1
		// `, nil, false},
		{"", `>= 1 admission (ICD-10 codes as below) OR >= 2 physician claims (ICD-9 codes as below) |
			Disease                     |             ICD9 codes             |     ICD10 codes     |
										| Physician claims |  Hospital data  |                     |
			Pernicious anemia           |       281        |      281.0      |        D51.0        |
			Autoimmune hemolytic anemia |       283        |      283.0      |        D59.1        |
			Ankylosing spondylitis      |       720        |      720.0      |         M45         |`,
			`merge row 1
			merge row 2:2 col 2:3
			merge row 2:3 col 1 
			merge row 2:3 col 4
			style row 1:3 header
		`, nil, false},
	}
	fmt.Println(tests)
	// trace := utils.NewTrace(true, nil)
	// p := NewCommandParser(utils.DebugSettings(true)) //set to false to silence tracing
	// tab := types.NewTable()
	// var err error
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		//TODO: parse sources before passing control commands to p.ParseCommandLines
	// 		tab.cmdList, err = p.ParseCommandLines(types.NewControlSection(strings.Split(tt.source, "\n")))
	// 		tab.contents, err = types.NewTableContents(tt.tab)
	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("NewTableContents() error = %v, wantErr %v", err, tt.wantErr)
	// 			return
	// 		}
	// 		tab.normalizeMergeRanges()
	// 		gotMrlist, err := createMergeRangeList(tab.cmdList)
	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("createMergeRangeList() error = %v, wantErr %v", err, tt.wantErr)
	// 			return
	// 		}
	// 		for i, mr := range gotMrlist { //TODO: better test of the creation of sorted mergelist
	// 			trace.Printf("%d: %s\n", i, mr.testString())
	// 		}

	// 		//trace.Printf("%v\n", tab)
	// 		grid, err := createGridTable(tab.contents, gotMrlist)
	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("createGridTable() error = %v, wantErr %v", err, tt.wantErr)
	// 			return
	// 		}
	// 		trace.Printf("%v\n", grid)
	// 		// if !reflect.DeepEqual(gotMrlist, tt.wantMrlist) {
	// 		// 	t.Errorf("createMergeRangeList() = %v, want %v", gotMrlist, tt.wantMrlist)
	// 		// }
	// 	})
	// }
}
