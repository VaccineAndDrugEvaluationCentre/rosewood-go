package types

import (
	"testing"

	"github.com/drgo/rosewood/utils"
)

func TestParseTableData(t *testing.T) {
	const showOutput = true
	tests := []struct {
		name    string
		args    string
		want    string
		wantErr bool
	}{
		//following tests must produce an error; so wantErr is true
		{name: "empty input",
			args:    "",
			want:    "",
			wantErr: true,
		},
		{name: "input with only LF",
			args:    "\n",
			want:    "",
			wantErr: true,
		},
		{name: "missing field separator",
			args:    "hello\n",
			want:    "",
			wantErr: true,
		},
		{name: "input with 2 rows and no | in first row",
			args:    "merged line one\n subtitle1|32423|60%|\n",
			want:    "",
			wantErr: true,
		},

		{name: "missing field separator and LF",
			args:    "hello",
			want:    "",
			wantErr: true,
		},
		//as above but with CRLF
		{name: "input with only CRLF",
			args:    "\r\n",
			want:    "",
			wantErr: true,
		},
		{name: "missing field separator with CRLF",
			args:    "hello\r\n",
			want:    "",
			wantErr: true,
		},
		//following tests must NOT produce an error; so wantErr is false
		{name: "one row, 3 cols",
			args:    "subtitle1|32423|60%|\n",
			want:    "r1 c1: subtitle1|r1 c2: 32423|r1 c3: 60%|\n",
			wantErr: false,
		},
		{name: "two rows, 3 cols",
			args: "subtitle1|32423|60%|\nsubtitle2|0|0%|\n",
			want: "r1 c1: subtitle1|r1 c2: 32423|r1 c3: 60%|\n" +
				"r2 c1: subtitle2|r2 c2: 0|r2 c3: 0%|\n",
			wantErr: false,
		},
		{name: "two rows, 1 col",
			args: "subtitle1|\nsubtitle2|\n",
			want: "r1 c1: subtitle1|\n" +
				"r2 c1: subtitle2|\n",
			wantErr: false,
		},
		//following tests use CRLF as EOL
		{name: "one row, 3 cols, CRLF as EOL",
			args:    "subtitle1|32423|60%|\r\n",
			want:    "r1 c1: subtitle1|r1 c2: 32423|r1 c3: 60%|\n",
			wantErr: false,
		},
		{name: "two rows, 3 cols, CRLF as EOL",
			args: "subtitle1|32423|60%|\r\nsubtitle2|0|0%|\r\n",
			want: "r1 c1: subtitle1|r1 c2: 32423|r1 c3: 60%|\n" +
				"r2 c1: subtitle2|r2 c2: 0|r2 c3: 0%|\n",
			wantErr: false,
		},
		{name: "two rows, 1 col, CRLF as EOL",
			args: "subtitle1|\r\nsubtitle2|\r\n",
			want: "r1 c1: subtitle1|\n" +
				"r2 c1: subtitle2|\n",
			wantErr: false,
		},
		{name: "one row, 1 col, LF as EOL",
			args:    "subtitle1|\n",
			want:    "r1 c1: subtitle1|\n",
			wantErr: false,
		},
		{name: "one row, 1 col, CRLF as EOL",
			args:    "subtitle1|\r\n",
			want:    "r1 c1: subtitle1|\n",
			wantErr: false,
		},
	}
	trace := utils.NewTrace(true, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTableContents(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTableData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.want {
				t.Errorf("ParseTableData() = [%v], want [%v]", got, tt.want)
			}
			if showOutput {
				trace.Printf("%s (%q): \n %v \n", tt.name, tt.args, got)
			}
		})
	}
}

// func Test_tableContents_ValidateRange(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		tab     string
// 		cmd     string
// 		want    Range
// 		wantErr bool
// 	}{
// 		{"", "subtitle1|32423|60%|\nsubtitle2|0|0%|\n", "merge row 1", makeRange(1, 1, 1, 3), false},
// 		{"", "subtitle1|32423|60%|\nsubtitle2|0|0%|\n", "merge row 1 col 1:3", makeRange(1, 1, 1, 3), false},
// 		{"", "subtitle1|32423|60%|\nsubtitle2|0|0%|\n", "merge row 2 col 1:2", makeRange(1, 1, 1, 3), false},
// 		{"", "subtitle1|32423|60%|\nsubtitle2|0|0%|\n", "merge row 2:3 col 1", makeRange(1, 1, 1, 3), false},
// 		{"", "subtitle1|32423|60%|\nsubtitle2|0|0%|\n", "merge row 2:3 col 4", makeRange(1, 1, 1, 3), false},
// 	}
// 	p := NewCommandParser(DefaultSettings())
// 	for _, tt := range tests {
// 		t.Run(tt.cmd, func(t *testing.T) {
// 			tab := monadicParseTableData(tt.tab)
// 			cmd, err := p.ParseCommands(strings.NewReader(tt.cmd))
// 			got, err := tab.validateRange(cmd[0].cellRange)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("got= %v, want %v", got.testString(), tt.want.testString())
// 			}
// 		})
// 	}
// }

/********* test helpers ************/
//for testing only; it ignores errors
func monadicParseTableData(s string) *TableContents {
	t, _ := NewTableContents(s)
	return t
}
