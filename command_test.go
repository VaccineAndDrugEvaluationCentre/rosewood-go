package carpenter

import (
	"fmt"
	"strings"
	"testing"
)

func TestCommandParser_ParseOneLineCommands(t *testing.T) {
	tests := []struct {
		source    string
		length    int
		wantError bool
		want      string
	}{
		{"merge row 1:2 col 1 \n", 1, false, "merge row 1:2 col 1:-1 "},
		{"merge row 1:2 col 1 ", 1, false, "merge row 1:2 col 1:-1 "}, //no linefeed
		{"merge row 1:2 col 1:2 \n", 1, false, "merge row 1:2 col 1:2 "},
		{"merge row 1 col 1 \n", 1, false, "merge row 1:-1 col 1:-1 "},
		{"Merge row 1 col 1 \n", 1, false, "merge row 1:-1 col 1:-1 "}, //test case-insensitive
		{"merge row 1 COL 1 \n", 1, false, "merge row 1:-1 col 1:-1 "}, //test case-insensitive
		// syntax errors, hence wantError=true
		{"merge col 1:2 row 1 ", 1, true, "col and row switched"},
		{"merge cl 1:2 row 1 ", 1, true, "misspelled col"},
		{"merge cl 1:2 raw 1 ", 1, true, "misspelled row"},
		{"merge row 1.4 col 1 \n", 1, true, "coordinate not integer"},
		{"merge row x1.4 col 1 \n", 1, true, "coordinate not integer"},
		{"merge row 0 col 1 \n", 1, true, "zero row/col value"},
		{"merge row 1 col 0 \n", 1, true, "zero row/col value"},
		{"merge row 0.0 col 1 \n", 1, true, "zero row/col value"},
		{"merge row col 1 \n", 1, true, "missing row coordinate"},
		{"merge row 1 col \n", 1, true, "missing col number"},
		{"merge row 1 col ", 1, true, "missing col number, no EOL"}, //no linefeed
		{"merge row 1;2 col 1:2 \n", 1, true, ": misstyped"},
		{"merge row 1:2 col 1;2 \n", 1, true, ": misstyped"},
		{"\n", 1, true, "empty line"},
		{"\r\n", 1, true, "empty line Windows"},
		{"\n\n\n", 1, true, "empty lines"},
		{"\r\n\r\n\r\n", 1, true, "empty lines Windows"},
		{"", 1, true, "empty input"}, //no linefeed
		{"row 1:2 col 1:2 \n", 1, true, "missing command identifier"},
		{"trump row 1:2 col 1:2 \n", 1, true, "wrong commander"},
		{`set rangeseparator "-"
			`, 1, false, "set rangeseparator,\"-\""}, //escaping " using \
	}
	p := NewCommandParser(nil) //use default settings
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got, err := p.ParseCommands(strings.NewReader(tt.source))
			fmt.Println(tt.source)
			if tt.wantError != (err != nil) {
				t.Errorf("Error handling failed, wanted %t, got %t \n errors %s:", tt.wantError, err != nil, p.errors.String())
			}
			if err != nil {
				return //if error was correctly reported by the parser do not continue testing
			}
			if got == nil {
				t.Error("ParseCommands() returns nil CommandList")
				return
			}
			if len(got) != tt.length {
				t.Errorf("Length of commands is incorrect, wanted %d, got %d", tt.length, len(got))
			}
			if len(got) == 0 {
				return
			}
			if got[0].String() != tt.want {
				t.Errorf("Commands parsed incorrectly, wanted %q, got %q", tt.want, got[0].String())
			}
		})
	}
}

func TestCommandParser_ParseMultiLineCommands(t *testing.T) {
	tests := []struct {
		source    string
		length    int
		wantError bool
		want      string
	}{
		{`merge row 1:2 col 1
		merge	row 1:2 col 1:2
		`, 2, false, ""},
		{`merge row 1 col 1:2
		row 1:2 col 1:2
		`, 2, true, ""},
		{`merge row 1:2 col 
		merge	row 1:2 col 1:2
		`, 2, true, ""},
	}
	p := NewCommandParser(nil) //use default settings
	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			got, err := p.ParseCommands(strings.NewReader(tt.source))
			fmt.Println(tt.source)
			if tt.wantError != (err != nil) {
				t.Errorf("Error handling failed, wanted %t, got %t \n errors %s:", tt.wantError, err != nil, p.errors.String())
			}
			if err != nil {

				return //if error was correctly reported by the parser do not continue testing
			}
			if got == nil {
				t.Error("ParseCommands() returns nil CommandList")
				return
			}
			if len(got) != tt.length {
				t.Errorf("Length of commands is incorrect, wanted %d, got %d", tt.length, len(got))
			}
			// if len(got) == 0 {
			// 	return
			// }
			// if got[0].String() != tt.want {
			// 	t.Errorf("Commands parsed incorrectly, wanted %q, got %q", tt.want, got[0].String())
			// }
		})
	}
}
