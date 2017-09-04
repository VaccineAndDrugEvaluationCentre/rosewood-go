package rosewood

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestGetWorkingDir(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetWorkingDir(); got != tt.want {
				t.Errorf("GetWorkingDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_iif(t *testing.T) {
	type args struct {
		condition bool
		t         interface{}
		f         interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iif(tt.args.condition, tt.args.t, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("iif() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveIfExists(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RemoveIfExists(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("RemoveIfExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFlush(t *testing.T) {
	bw := bufio.NewWriter(os.Stdout)
	fmt.Fprintf(bw, "test")
	bw.Flush()
}
