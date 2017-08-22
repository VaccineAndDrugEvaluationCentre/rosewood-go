package rosewood

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_genAllPossibleRangePoints(t *testing.T) {
	type args struct {
		p1 RwInt
		p2 RwInt
		by RwInt
	}
	tests := []struct {
		name      string
		args      args
		wantPList []RwInt
	}{
		{"1:2:10", args{1, 10, 2}, []RwInt{1, 3, 5, 7, 9}},
		{"0:2:10", args{0, 10, 2}, []RwInt{0, 2, 4, 6, 8, 10}},
		{"1:11:10", args{1, 10, 11}, nil},
		{"1:11:0", args{1, 11, 0}, nil},
		{"1:11:11", args{1, 11, 11}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPList := genAllPossibleRangePoints(tt.args.p1, tt.args.p2, tt.args.by); !reflect.DeepEqual(gotPList, tt.wantPList) {
				t.Errorf("genAllPossibleRangePoints() = %v, want %v", gotPList, tt.wantPList)
			}
		})
	}
}

func Test_expandSpan(t *testing.T) {
	const na = MissingRwInt
	const showOutput = true
	tests := []struct {
		cs           span
		wantListSize int
		wantSpan1    span
		wantErr      bool
	}{
		{span{1, 6, 1, 6, na, na, nil, nil}, 1, span{1, 6, 1, 6, na, na, nil, nil}, false},
		{span{1, 6, 1, 6, 2, 2, nil, nil}, 9, span{1, na, 1, na, na, na, nil, nil}, false},
		{span{1, 6, 1, 6, 2, na, nil, nil}, 3, span{1, na, 1, 6, na, na, nil, nil}, false},
		{span{1, 6, 1, 6, na, 2, nil, nil}, 3, span{1, 6, 1, na, na, na, nil, nil}, false},
		{span{1, 6, 1, 6, na, na, []RwInt{1, 3, 5}, nil}, 3, span{1, na, 1, 6, na, na, nil, nil}, false},
		{span{11, 16, 1, 6, 2, na, []RwInt{1, 3, 5}, nil}, 6, span{11, na, 1, 6, na, na, nil, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.cs.testString(), func(t *testing.T) {
			gotSList, err := expandSpan(tt.cs)
			if showOutput {
				fmt.Printf("%s:\n%v\n", tt.cs.testString(), gotSList)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("expandSpan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotSList) != tt.wantListSize {
				t.Errorf("wrong number of expanded spans = %d, wanted %v", len(gotSList), tt.wantListSize)
				return
			}
			if !reflect.DeepEqual(gotSList[0], tt.wantSpan1) {
				t.Errorf("wrong first expanded span = %v, want %v", gotSList[0], tt.wantSpan1)
			}
		})
	}
}

func Test_deduplicateSpanList(t *testing.T) {
	const na = MissingRwInt
	const showOutput = true
	tests := []struct {
		cs           span
		wantListSize int
		wantSpan1    span
	}{
		{span{11, 16, 1, 6, 2, na, []RwInt{1, 3, 5}, nil}, 6, span{11, na, 1, 6, na, na, nil, nil}},
		{span{1, 6, 1, 6, 2, na, []RwInt{1, 3, 5}, nil}, 3, span{11, na, 1, 6, na, na, nil, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.cs.testString(), func(t *testing.T) {
			gotSList, err := expandSpan(tt.cs)
			if err != nil {
				t.Errorf("expandSpan() error = %v", err)
				return
			}
			if showOutput {
				fmt.Printf("%s:\n%v\n", tt.cs.testString(), gotSList)
			}
			redupList := deduplicateSpanList(gotSList)
			if len(redupList) != tt.wantListSize {
				t.Errorf("wrong number of dedup spans = %d, wanted %v", len(redupList), tt.wantListSize)
				return
			}
			if showOutput {
				fmt.Printf("%s:\n%v\n", "deduplicated", redupList)
				fmt.Println(strings.Repeat("*", 30))
			}
			// if !reflect.DeepEqual(gotSList[0], tt.wantSpan1) {
			// 	t.Errorf("wrong first expanded span = %v, want %v", gotSList[0], tt.wantSpan1)
			// }
		})
	}
}
