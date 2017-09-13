package rosewood

import (
	"bytes"
	"testing"

	"github.com/drgo/rosewood/types"
	"github.com/drgo/rosewood/utils"
)

//arrays in Go cannot be declared constant
var (
	tabr3c4 = []*types.Row{ //simple table 3X4 no merging
		types.MakeRow(
			types.MakeCell("cell11", 1, 1, types.CsUndefined, 0, 0),
			types.MakeCell("cell12", 1, 2, types.CsUndefined, 0, 0),
			types.MakeCell("cell13", 1, 3, types.CsUndefined, 0, 0),
			types.MakeCell("cell14", 1, 4, types.CsUndefined, 0, 0),
		),
		types.MakeRow(
			types.MakeCell("cell21", 2, 1, types.CsUndefined, 0, 0),
			types.MakeCell("cell22", 2, 2, types.CsUndefined, 0, 0),
			types.MakeCell("cell23", 2, 3, types.CsUndefined, 0, 0),
			types.MakeCell("cell24", 2, 4, types.CsUndefined, 0, 0),
		),
		types.MakeRow(
			types.MakeCell("cell31", 3, 1, types.CsUndefined, 0, 0),
			types.MakeCell("cell32", 3, 2, types.CsUndefined, 0, 0),
			types.MakeCell("cell33", 3, 3, types.CsUndefined, 0, 0),
			types.MakeCell("cell34", 3, 4, types.CsUndefined, 0, 0),
		)}
)

func TestRender(t *testing.T) {
	tests := []struct {
		name    string
		tab     *types.Table
		want    string
		wantErr bool
	}{
		{name: "table 3X4 no merging",
			tab:     types.MakeTableContents(tabr3c4, 4),
			want:    ``,
			wantErr: false,
		},
		{name: "table 3X4 merge row 1:2 col 3:4",
			tab: types.MakeTableContents([]*types.Row{
				types.MakeRow(
					types.MakeCell("cell11", 1, 1, types.CsUndefined, 0, 0), //text, cell, row, hidden, rowspan, colspan
					types.MakeCell("cell12", 1, 2, types.CsUndefined, 0, 0),
					types.MakeCell("cell13", 1, 3, types.CsUndefined, 2, 2),
					types.MakeCell("cell14", 1, 4, types.CsMerged, 0, 0),
				),
				types.MakeRow(
					types.MakeCell("cell21", 2, 1, types.CsUndefined, 0, 0),
					types.MakeCell("cell22", 2, 2, types.CsUndefined, 0, 0),
					types.MakeCell("cell23", 2, 3, types.CsMerged, 0, 0),
					types.MakeCell("cell24", 2, 4, types.CsMerged, 0, 0),
				),
				types.MakeRow(
					types.MakeCell("cell31", 3, 1, types.CsUndefined, 0, 0),
					types.MakeCell("cell32", 3, 2, types.CsUndefined, 0, 0),
					types.MakeCell("cell33", 3, 3, types.CsUndefined, 0, 0),
					types.MakeCell("cell34", 3, 4, types.CsUndefined, 0, 0),
				),
			}, 4),
			want:    `<td rowspan="2" colspan="2">cell13</td></tr>`,
			wantErr: false,
		},
		{name: "table 3X4 merge row 1:1 col 1:4",
			tab: types.MakeTableContents([]*types.Row{
				types.MakeRow(
					types.MakeCell("cell11", 1, 1, types.CsUndefined, 0, 4),
					types.MakeCell("cell12", 1, 2, types.CsMerged, 0, 0),
					types.MakeCell("cell13", 1, 3, types.CsMerged, 0, 0),
					types.MakeCell("cell14", 1, 4, types.CsMerged, 0, 0),
				),
				types.MakeRow(
					types.MakeCell("cell21", 2, 1, types.CsUndefined, 0, 0),
					types.MakeCell("cell22", 2, 2, types.CsUndefined, 0, 0),
					types.MakeCell("cell23", 2, 3, types.CsUndefined, 0, 0),
					types.MakeCell("cell24", 2, 4, types.CsUndefined, 0, 0),
				),
				types.MakeRow(
					types.MakeCell("cell31", 3, 1, types.CsUndefined, 0, 0),
					types.MakeCell("cell32", 3, 2, types.CsUndefined, 0, 0),
					types.MakeCell("cell33", 3, 3, types.CsUndefined, 0, 0),
					types.MakeCell("cell34", 3, 4, types.CsUndefined, 0, 0),
				),
			}, 4),
			want:    `<tr><td colspan="4">cell11</td></tr>`,
			wantErr: false,
		},
		{name: "table 3X4 merge row 1:4 col 1:1",
			tab: types.MakeTableContents([]*types.Row{
				types.MakeRow(
					types.MakeCell("cell11", 1, 1, types.CsUndefined, 4, 0),
					types.MakeCell("cell12", 1, 2, types.CsUndefined, 0, 0),
					types.MakeCell("cell13", 1, 3, types.CsUndefined, 0, 0),
					types.MakeCell("cell14", 1, 4, types.CsUndefined, 0, 0),
				),
				types.MakeRow(
					types.MakeCell("cell21", 2, 1, types.CsMerged, 0, 0),
					types.MakeCell("cell22", 2, 2, types.CsUndefined, 0, 0),
					types.MakeCell("cell23", 2, 3, types.CsUndefined, 0, 0),
					types.MakeCell("cell24", 2, 4, types.CsUndefined, 0, 0),
				),
				types.MakeRow(
					types.MakeCell("cell31", 3, 1, types.CsMerged, 0, 0),
					types.MakeCell("cell32", 3, 2, types.CsUndefined, 0, 0),
					types.MakeCell("cell33", 3, 3, types.CsUndefined, 0, 0),
					types.MakeCell("cell34", 3, 4, types.CsUndefined, 0, 0),
				),
			}, 4),
			want:    `<tr><td rowspan="4">cell11</td><td>`,
			wantErr: false,
		},
	}
	// fileBuffer := &bytes.Buffer{}
	trace := utils.NewTrace(true, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			re := NewHTMLRenderer()
			re.SetWriter(w)
			re.SetSettings(utils.DebugSettings(true))
			re.SetTables([]*types.Table{tt.tab})
			tt.tab.Run()
			// err := tt.tab.Render(w, re)
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			// if !bytes.Contains(w.Bytes(), []byte(tt.want)) {
			// 	t.Errorf("Wanted string [%s] was not found", tt.want)
			// }
			if showOutput {
				trace.Printf("%s: \n %q \n", tt.name, w.String())
			}
			// if writeOut {
			// 	w.WriteString("<p><p>")
			// 	w.WriteTo(fileBuffer)
			// }
		})
	}
	// if writeOut {
	// 	outFileName := path.Join(testDirName, "rendertest."+testFileExt)
	// 	if err := ioutil.WriteFile(outFileName, fileBuffer.Bytes(), 0644); err != nil {
	// 		t.Errorf("failed to write to file %s: %v", outFileName, err)
	// 	}
	// 	trace.Printf("Results saved to file://%s\n", outFileName)
	// }
}
