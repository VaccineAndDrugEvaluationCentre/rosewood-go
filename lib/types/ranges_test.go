// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

import (
	"reflect"
	"testing"
)

func Test_genAllPossibleRangePoints(t *testing.T) {
	type args struct {
		p1 int
		p2 int
		by int
	}
	tests := []struct {
		name      string
		args      args
		wantPList []int
	}{
		{"1:2:10", args{1, 10, 2}, []int{1, 3, 5, 7, 9}},
		{"0:2:10", args{0, 10, 2}, []int{0, 2, 4, 6, 8, 10}},
		{"1:11:10", args{1, 10, 11}, []int{1}},
		{"1:11:11", args{1, 11, 11}, []int{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPList := genAllPossibleRangePoints(tt.args.p1, tt.args.p2, tt.args.by); !reflect.DeepEqual(gotPList, tt.wantPList) {
				t.Errorf("genAllPossibleRangePoints() = %v, want %v", gotPList, tt.wantPList)
			}
		})
	}
}

// func Test_ExpandSpan(t *testing.T) {
// 	const na = RwMissing
// 	const showOutput = true
// 	tests := []struct {
// 		cs           Span
// 		wantListSize int
// 		wantSpan1    Span
// 		wantErr      bool
// 	}{
// 		{Span{1, 6, 1, 6, na, na, nil, nil}, 1, Span{1, 6, 1, 6, na, na, nil, nil}, false},
// 		{Span{1, 6, 1, 6, 2, 2, nil, nil}, 9, Span{1, na, 1, na, na, na, nil, nil}, false},
// 		{Span{1, 6, 1, 6, 2, na, nil, nil}, 3, Span{1, na, 1, 6, na, na, nil, nil}, false},
// 		{Span{1, 6, 1, 6, na, 2, nil, nil}, 3, Span{1, 6, 1, na, na, na, nil, nil}, false},
// 		{Span{1, 6, 1, 6, na, na, []int{1, 3, 5}, nil}, 3, Span{1, na, 1, 6, na, na, nil, nil}, false},
// 		{Span{11, 16, 1, 6, 2, na, []int{1, 3, 5}, nil}, 6, Span{11, na, 1, 6, na, na, nil, nil}, false},
// 	}
// 	trace := trace.NewTrace(showOutput, nil)
// 	for _, tt := range tests {
// 		t.Run(tt.cs.testString(), func(t *testing.T) {
// 			gotSList, err := tt.cs.ExpandSpan()
// 			if showOutput {
// 				trace.Printf("%s:\n%v\n", tt.cs.testString(), gotSList)
// 			}
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ExpandSpan() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if len(gotSList) != tt.wantListSize {
// 				t.Errorf("wrong number of expanded Spans = %d, wanted %v", len(gotSList), tt.wantListSize)
// 				return
// 			}
// 			if !reflect.DeepEqual(*gotSList[0], tt.wantSpan1) {
// 				t.Errorf("wrong first expanded Span = %v, want %v", gotSList[0], tt.wantSpan1)
// 			}
// 		})
// 	}
// }

// func Test_deduplicateSpanList(t *testing.T) {
// 	const na = RwMissing
// 	const showOutput = true
// 	tests := []struct {
// 		cs           Span
// 		wantListSize int
// 		wantSpan1    Span
// 	}{
// 		{Span{11, 16, 1, 6, 2, na, []int{1, 3, 5}, nil}, 6, Span{11, na, 1, 6, na, na, nil, nil}},
// 		{Span{1, 6, 1, 6, 2, na, []int{1, 3, 5}, nil}, 3, Span{11, na, 1, 6, na, na, nil, nil}},
// 	}
// 	trace := trace.NewTrace(showOutput, nil)
// 	for _, tt := range tests {
// 		t.Run(tt.cs.testString(), func(t *testing.T) {
// 			gotSList, err := tt.cs.ExpandSpan()
// 			if err != nil {
// 				t.Errorf("ExpandSpan() error = %v", err)
// 				return
// 			}
// 			if showOutput {
// 				trace.Printf("%s:\n%v\n", tt.cs.testString(), gotSList)
// 			}
// 			redupList := DeduplicateSpanList(gotSList)
// 			if len(redupList) != tt.wantListSize {
// 				t.Errorf("wrong number of dedup Spans = %d, wanted %v", len(redupList), tt.wantListSize)
// 				return
// 			}
// 			if showOutput {
// 				trace.Printf("%s:\n%v\n", "deduplicated", redupList)
// 				trace.Println(strings.Repeat("*", 30))
// 			}
// 			// if !reflect.DeepEqual(gotSList[0], tt.wantSpan1) {
// 			// 	t.Errorf("wrong first expanded Span = %v, want %v", gotSList[0], tt.wantSpan1)
// 			// }
// 		})
// 	}
// }

func Test_normalizeSpan(t *testing.T) {
	type args struct {
		cs       *Span
		rowCount int
		colCount int
	}
	tests := []struct {
		name string
		args args
		want *Span
	}{
		{"no-missing", args{MakeSpan(1, 1, 4, 4), 6, 4}, MakeSpan(1, 1, 4, 4)},
		{"r1-r2-missing", args{MakeSpan(RwMissing, RwMissing, 2, 3), 6, 4}, MakeSpan(1, 6, 2, 3)},
		{"c1-c2-missing", args{MakeSpan(1, 1, RwMissing, RwMissing), 6, 4}, MakeSpan(1, 1, 1, 4)},
		{"r2missing", args{MakeSpan(1, RwMissing, 2, 2), 6, 4}, MakeSpan(1, 1, 2, 2)},
		{"c2missing", args{MakeSpan(1, 1, 2, RwMissing), 6, 4}, MakeSpan(1, 1, 2, 2)},
		{"r1-r2-c2-missing", args{MakeSpan(RwMissing, RwMissing, RwMissing, 3), 6, 4}, MakeSpan(1, 6, 3, 3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.cs.Normalize(tt.args.rowCount, tt.args.colCount)
			if !reflect.DeepEqual(tt.args.cs, tt.want) {
				t.Errorf("normalizeSpan() = %v, want %v", tt.args.cs, tt.want)
			}
		})
	}
}
