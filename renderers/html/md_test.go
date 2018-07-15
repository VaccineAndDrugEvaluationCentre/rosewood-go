package html

import (
	"testing"
)

func TestInlinedMdToHTML(t *testing.T) {
	tests := []struct {
		md      string
		want    string
		wantErr bool
	}{
		{"x is > y", "x is &gt; y", false},
		{"**strong**", "<strong>strong</strong>", false},
		{"*em*", "<em>em</em>", false},
		{"__strong__", "<strong>strong</strong>", false},
		{"_em_", "<em>em</em>", false},
		{"**The rates and *95%CI* of anemia^1^ among olive **oil drinkers ~see definition in appendix~",
			"<strong>The rates and <em>95%CI</em> of anemia<sup>1</sup> among olive </strong>oil drinkers <sub>see definition in appendix</sub>", false},
		{"table 1 model^1^", "table 1 model<sup>1</sup>", false},
		{"table 1 model~1~", "table 1 model<sub>1</sub>", false},
		{"**strong*em* **", "<strong>strong<em>em</em> </strong>", false},
		{"** *strongem* **", "<strong> <em>strongem</em> </strong>", false},
		{"**_strongem_**", "<strong><em>strongem</em></strong>", false},
		{"\\*em\\*", "*em*", false},
		{"\\~tilde\\~", "~tilde~", false},
		//errors
		{"**strong*", "<strong>strong<em>", true},
		{"**strong*em***", "<strong>strong<em>em<strong><em>", true},
	}
	for _, tt := range tests {
		t.Run(tt.md, func(t *testing.T) {
			got, err := InlinedMdToHTML(tt.md, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("InlinedMdToHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("InlinedMdToHTML() result = %v, want %v", string(got), tt.want)
			}
		})
	}
}
