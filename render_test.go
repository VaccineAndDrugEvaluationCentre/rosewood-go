package carpenter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

const (
	showOutput = true
	writeOut   = true
)

func TestRender(t *testing.T) {
	type args struct {
		r *HtmlRenderer
		t *Table
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "render table with one row, 3 cols",
			args: args{r: NewHtmlRenderer(),
				t: monadicParseTableData("subtitle1|32423|60%|\nsubtitle2|0|0%|\n"),
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := Render(w, tt.args.r, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if gotW := w.String(); gotW != tt.want {
			// 	t.Errorf("Render() = %v, want %v", gotW, tt.want)
			// }
			if showOutput {
				fmt.Printf("%s: \n %q \n", tt.name, w.String())
			}
			if writeOut {
				if err := ioutil.WriteFile("test.html", w.Bytes(), 0644); err != nil {
					t.Errorf("failed to write to file %s: %v", "test.html", err)
				}

			}
		})
	}
}
