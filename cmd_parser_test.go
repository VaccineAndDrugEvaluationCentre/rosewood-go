package rosewood

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
		{"merge row 1:2 col 1 \n", 1, false, "merge row 1:2 col 1:NA "},
		{"merge row 1:2 col 1 ", 1, false, "merge row 1:2 col 1:NA "}, //no linefeed
		{"merge row 1:2 col 1:2 \n", 1, false, "merge row 1:2 col 1:2 "},
		{"merge row 1 col 1 \n", 1, false, "merge row 1:NA col 1:NA "},
		{"Merge row 1 col 1 \n", 1, false, "merge row 1:NA col 1:NA "},                           //test case-insensitive
		{"merge row 1 COL 1 \n", 1, false, "merge row 1:NA col 1:NA "},                           //test case-insensitive
		{"style row 1 col 1 style1\n", 1, false, "style row 1:NA col 1:NA style1"},               //test args x 1
		{"style row 1 col 1 style1 style2\n", 1, false, "style row 1:NA col 1:NA style1,style2"}, //test args x 2
		{"style row 1 style1\n", 1, false, "style row 1:NA col NA:NA style1"},                    //test args x 1 no col
		{"style row 1 style1 style2\n", 1, false, "style row 1:NA col NA:NA style1,style2"},      //test args x 2 no col
		//{"style row 1,7 style1 \n", 1, false, "style row 1:NA col 0:NA style1"},       //test args x 2
		{`set rangeseparator "-"
			`, 1, false, "set rangeseparator,\"-\""}, //escaping " using \

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
		{"merge row 3:1 col 1\n", 1, true, "row numbers invalid"},
		{"merge row 1:2 col 3:2\n", 1, true, "col numbers invalid"},
		{"\n", 1, true, "empty line"},
		{"\r\n", 1, true, "empty line Windows"},
		{"\n\n\n", 1, true, "empty lines"},
		{"\r\n\r\n\r\n", 1, true, "empty lines Windows"},
		{"", 1, true, "empty input"}, //no linefeed
		{"row 1:2 col 1:2 \n", 1, true, "missing command identifier"},
		{"trump row 1:2 col 1:2 \n", 1, true, "wrong commander"},
		// {`set rangeseparator "-" "onemore"
		// 	`, 1, true, "invalid # args to set"},

	}
	p := NewCommandParser(DefaultSettings()) //use default settings
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got, err := p.ParseCommands(strings.NewReader(tt.source))
			fmt.Println(tt.source)
			if tt.wantError != (err != nil) {
				t.Errorf("Error handling failed, wanted %t, got %t \nerrors %s:\n", tt.wantError, err != nil, p.errors.String())
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
	p := NewCommandParser(DefaultSettings()) //use default settings
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

const script1 = `
/*following commands should parse without error*/
merge row 1:2 col 1
merge row 1:2 col 1
merge row 1:2 col 1:2   //another comment
merge row 1 col 1 
merge row 1 COL 1 
style row 1 col 1 style1
style row 1 col 1 style1
style row 1 style1
style row 1 style1 style2
set rangeseparator "-"
//line comment
`

//set rangeseparator "-"
func TestCommandParser_ParseFullScript(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		length    int
		wantError bool
		want      string
	}{
		{"Script 1", script1, 10, false, ""},
	}
	p := NewCommandParser(DefaultSettings()) //use default settings

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.ParseCommands(strings.NewReader(tt.source))
			// fmt.Println(tt.source)
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
