package rosewood

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_createGridTable(t *testing.T) {
	tests := []struct {
		name       string
		tab        string
		source     string
		wantMrlist []Range
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
			Ankylosing spondylitis      |       720        |      720.0      |         M45         |
			`,
			`merge row 1
			merge row 2:2 col 2:3
			merge row 2:3 col 1 
			merge row 2:3 col 4
			style row 1:3 header
		`, nil, false},
	}
	p := NewCommandParser(DefaultSettings()) //use default settings
	tab := newTable()
	var err error
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tab.cmdList, err = p.ParseCommands(strings.NewReader(tt.source))
			tab.contents, err = NewTableContents(tt.tab)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTableContents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tab.normalizeMergeRanges()
			gotMrlist, err := createMergeRangeList(tab.cmdList)
			if (err != nil) != tt.wantErr {
				t.Errorf("createMergeRangeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, mr := range gotMrlist { //todo: better test of the creation of sorted mergelist
				fmt.Printf("%d: %s\n", i, mr.testString())
			}

			//fmt.Printf("%v\n", tab)
			grid, err := createGridTable(tab.contents, gotMrlist)
			if (err != nil) != tt.wantErr {
				t.Errorf("createGridTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("%v\n", grid)
			// if !reflect.DeepEqual(gotMrlist, tt.wantMrlist) {
			// 	t.Errorf("createMergeRangeList() = %v, want %v", gotMrlist, tt.wantMrlist)
			// }
		})
	}
}

func Test_normalizeSpan(t *testing.T) {
	type args struct {
		cs       span
		rowCount RwInt
		colCount RwInt
	}
	tests := []struct {
		name string
		args args
		want span
	}{
		{"no-missing", args{makeSpan(1, 1, 4, 4), 6, 4}, makeSpan(1, 1, 4, 4)},
		{"r1-r2-missing", args{makeSpan(MissingRwInt, MissingRwInt, 2, 3), 6, 4}, makeSpan(1, 6, 2, 3)},
		{"c1-c2-missing", args{makeSpan(1, 1, MissingRwInt, MissingRwInt), 6, 4}, makeSpan(1, 1, 1, 4)},
		{"r2missing", args{makeSpan(1, MissingRwInt, 2, 2), 6, 4}, makeSpan(1, 1, 2, 2)},
		{"c2missing", args{makeSpan(1, 1, 2, MissingRwInt), 6, 4}, makeSpan(1, 1, 2, 2)},
		{"r1-r2-c2-missing", args{makeSpan(MissingRwInt, MissingRwInt, MissingRwInt, 3), 6, 4}, makeSpan(1, 6, 3, 3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeSpan(tt.args.cs, tt.args.rowCount, tt.args.colCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("normalizeSpan() = %v, want %v", got, tt.want)
			}
		})
	}
}
