// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

// func Test_spanToRangeList(t *testing.T) {
// 	type args struct {
// 		cmdList []*Command
// 		cmdType RwKeyWord
// 	}
// 	tests := []struct {
// 		name      string
// 		args      args
// 		wantRList []Range
// 		wantErr   bool
// 	}{}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotRList, err := spanToRangeList(tt.args.cmdList, tt.args.cmdType)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("spanToRangeList() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotRList, tt.wantRList) {
// 				t.Errorf("spanToRangeList() = %v, want %v", gotRList, tt.wantRList)
// 			}
// 		})
// 	}
// }
