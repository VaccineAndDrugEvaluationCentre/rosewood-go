package carpenter

import (
	"fmt"
	"testing"
)

func TestParseTableData(t *testing.T) {
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
		{name: "input with no | but with LF",
			args:    "hello\n",
			want:    "",
			wantErr: true,
		},
		{name: "input with no | and no LF",
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
		{name: "input with no | but with CRLF",
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
	}
	const showOutput = true //false to suppress printing returned table structs
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTableData(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTableData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.want {
				t.Errorf("ParseTableData() = [%v], want [%v]", got, tt.want)
			}
			if showOutput {
				fmt.Printf("%s (%q): \n %v \n", tt.name, tt.args, got)
			}
		})
	}
}
