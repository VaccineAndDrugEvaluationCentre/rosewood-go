package rosewood

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"testing"
)

const (
	showOutput = true
	writeOut   = true
	//	developDir   = "/Users/salah/Dropbox/code/go/src/github.com/drgo/carpenter" //not needed remove
	testDirName  = "test-files" //relative to the executable directory
	testFileName = "rendertest"
	testFileExt  = "html"
)

func genTestTableData(rows []*Row) *table {
	return &table{contents: &tableContents{
		rows:        rows,
		maxFldCount: 5,
	}}
}

//arrays in Go cannot be declared constant
var (
	tabr3c4 = []*Row{ //simple table 3X4 no merging
		&Row{[]*Cell{
			&Cell{"cell11", 1, 1, csUndefined, 0, 0}, //text, cell, row, state, rowspan, colspan
			&Cell{"cell12", 1, 2, csUndefined, 0, 0},
			&Cell{"cell13", 1, 3, csUndefined, 0, 0},
			&Cell{"cell14", 1, 4, csUndefined, 0, 0}},
		},
		&Row{[]*Cell{
			&Cell{"cell21", 2, 1, csUndefined, 0, 0},
			&Cell{"cell22", 2, 2, csUndefined, 0, 0},
			&Cell{"cell23", 2, 3, csUndefined, 0, 0},
			&Cell{"cell24", 2, 4, csUndefined, 0, 0}},
		},
		&Row{[]*Cell{
			&Cell{"cell31", 3, 1, csUndefined, 0, 0},
			&Cell{"cell32", 3, 2, csUndefined, 0, 0},
			&Cell{"cell33", 3, 3, csUndefined, 0, 0},
			&Cell{"cell34", 3, 4, csUndefined, 0, 0}},
		}}
)

func TestRender(t *testing.T) {
	type args struct {
		r *HtmlRenderer
		s *Settings
		t *table
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// {name: "table 2X3",
		// 	args: args{r: NewHtmlRenderer(),
		// 		t: monadicParseTableData("subtitle1|32423|60%|\nsubtitle2|0|0%|\nsubtitle3|1.5|1.5%|\n"),
		// 	},
		// 	want:    "",
		// 	wantErr: false,
		// },
		{name: "table 3X4 no merging",
			args: args{r: NewHtmlRenderer(),
				t: genTestTableData(tabr3c4),
			},
			want:    ``,
			wantErr: false,
		},
		{name: "table 3X4 merge row 1:2 col 3:4",
			args: args{r: NewHtmlRenderer(),
				t: genTestTableData([]*Row{
					&Row{[]*Cell{
						&Cell{"cell11", 1, 1, csUndefined, 0, 0}, //text, cell, row, hidden, rowspan, colspan
						&Cell{"cell12", 1, 2, csUndefined, 0, 0},
						&Cell{"cell13", 1, 3, csUndefined, 2, 2},
						&Cell{"cell14", 1, 4, csMerged, 0, 0}},
					},
					&Row{[]*Cell{
						&Cell{"cell21", 2, 1, csUndefined, 0, 0},
						&Cell{"cell22", 2, 2, csUndefined, 0, 0},
						&Cell{"cell23", 2, 3, csMerged, 0, 0},
						&Cell{"cell24", 2, 4, csMerged, 0, 0}},
					},
					&Row{[]*Cell{
						&Cell{"cell31", 3, 1, csUndefined, 0, 0},
						&Cell{"cell32", 3, 2, csUndefined, 0, 0},
						&Cell{"cell33", 3, 3, csUndefined, 0, 0},
						&Cell{"cell34", 3, 4, csUndefined, 0, 0}},
					},
				}),
			},
			want:    `<td rowspan="2" colspan="2">cell13</td></tr>`,
			wantErr: false,
		},
		{name: "table 3X4 merge row 1:1 col 1:4",
			args: args{r: NewHtmlRenderer(),
				t: genTestTableData([]*Row{
					&Row{[]*Cell{
						&Cell{"cell11", 1, 1, csUndefined, 0, 4}, //text, cell, row, hidden, rowspan, colspan
						&Cell{"cell12", 1, 2, csMerged, 0, 0},
						&Cell{"cell13", 1, 3, csMerged, 0, 0},
						&Cell{"cell14", 1, 4, csMerged, 0, 0}},
					},
					&Row{[]*Cell{
						&Cell{"cell21", 2, 1, csUndefined, 0, 0},
						&Cell{"cell22", 2, 2, csUndefined, 0, 0},
						&Cell{"cell23", 2, 3, csUndefined, 0, 0},
						&Cell{"cell24", 2, 4, csUndefined, 0, 0}},
					},
					&Row{[]*Cell{
						&Cell{"cell31", 3, 1, csUndefined, 0, 0},
						&Cell{"cell32", 3, 2, csUndefined, 0, 0},
						&Cell{"cell33", 3, 3, csUndefined, 0, 0},
						&Cell{"cell34", 3, 4, csUndefined, 0, 0}},
					},
				}),
			},
			want:    `<tr><td colspan="4">cell11</td></tr>`,
			wantErr: false,
		},
		{name: "table 3X4 merge row 1:4 col 1:1",
			args: args{r: NewHtmlRenderer(),
				t: genTestTableData([]*Row{
					&Row{[]*Cell{
						&Cell{"cell11", 1, 1, csUndefined, 4, 0}, //text, cell, row, hidden, rowspan, colspan
						&Cell{"cell12", 1, 2, csUndefined, 0, 0},
						&Cell{"cell13", 1, 3, csUndefined, 0, 0},
						&Cell{"cell14", 1, 4, csUndefined, 0, 0}},
					},
					&Row{[]*Cell{
						&Cell{"cell21", 2, 1, csMerged, 0, 0},
						&Cell{"cell22", 2, 2, csUndefined, 0, 0},
						&Cell{"cell23", 2, 3, csUndefined, 0, 0},
						&Cell{"cell24", 2, 4, csUndefined, 0, 0}},
					},
					&Row{[]*Cell{
						&Cell{"cell31", 3, 1, csMerged, 0, 0},
						&Cell{"cell32", 3, 2, csUndefined, 0, 0},
						&Cell{"cell33", 3, 3, csUndefined, 0, 0},
						&Cell{"cell34", 3, 4, csUndefined, 0, 0}},
					},
				}),
			},
			want:    `<tr><td rowspan="4">cell11</td><td>`,
			wantErr: false,
		},
	}
	fileBuffer := &bytes.Buffer{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := render(w, tt.args.r, DefaultSettings(), tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Contains(w.Bytes(), []byte(tt.want)) {
				t.Errorf("Wanted string [%s] was not found", tt.want)
			}
			if showOutput {
				fmt.Printf("%s: \n %q \n", tt.name, w.String())
			}
			if writeOut {
				w.WriteString("<p><p>")
				w.WriteTo(fileBuffer)
			}
		})
	}
	if writeOut {
		outFileName := path.Join(testDirName, testFileName+"."+testFileExt)
		if err := ioutil.WriteFile(outFileName, fileBuffer.Bytes(), 0644); err != nil {
			t.Errorf("failed to write to file %s: %v", outFileName, err)
		}
		fmt.Printf("Results saved to file://%s\n", outFileName)
	}
}

///Users/salah/Dropbox/code/go/src/github.com/drgo/carpenter/test-files/rendertest.html
