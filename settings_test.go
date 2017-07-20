package carpenter

import (
	"fmt"
	"reflect"
	"testing"
)

//just for testing
type settingsMap map[string]interface{}

func TestSettings_Get(t *testing.T) {
	type fields struct {
		items settingsMap
	}
	tests := []struct {
		name      string
		fields    fields
		args      string
		wantValue interface{}
		wantOk    bool
	}{
		{"string", fields{settingsMap{"filename": "myfilename"}}, "filename", "myfilename", true},
		{"int", fields{settingsMap{"maxopenfiles": 10}}, "maxopenfiles", 10, true},
		{"bool", fields{settingsMap{"isDir": true}}, "isDir", true, true},
		//the following must fail
		{"int", fields{settingsMap{"maxopenfiles": 10}}, "notmaxopenfiles", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Settings{
				items: tt.fields.items,
			}
			gotValue, gotOk := s.Get(tt.args)
			if gotOk != tt.wantOk {
				t.Errorf("Settings.Get() gotOk = %v, want %v", gotOk, tt.wantOk)
				return
			}
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("Settings.Get() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestSettings_Set(t *testing.T) {
	type fields struct {
		items map[string]interface{}
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Settings{
				items: tt.fields.items,
			}
			if got := s.Set(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Settings.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleNewSettings() {
	s := NewSettings()
	printValue := func(key string) {
		value, ok := s.Get(key)
		if ok {
			fmt.Println(value)
		} else {
			fmt.Println("!")
		}
	}
	s.Set("FontName", "Manlo")
	s.Set("FontSize", 12)
	s.Set("FontItalic", true)
	printValue("FontName")
	printValue("FontSize")
	printValue("FontItalic")
	printValue("doesnotexist")
	//Output:
	//Manlo
	//12
	//true
	//!
}

func TestSettings_LoadSettings(t *testing.T) {
	type fields struct {
		items map[string]interface{}
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Settings{
				items: tt.fields.items,
			}
			if err := s.LoadSettings(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Settings.LoadSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSettings_SaveSettings(t *testing.T) {
	type fields struct {
		items map[string]interface{}
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Settings{
				items: tt.fields.items,
			}
			if err := s.SaveSettings(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Settings.SaveSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
