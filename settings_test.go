package carpenter

//just for testing
type settingsMap map[string]interface{}

// func TestSettings_Get(t *testing.T) {
// 	type fields struct {
// 		items settingsMap
// 	}
// 	tests := []struct {
// 		name      string
// 		fields    fields
// 		args      string
// 		wantValue interface{}
// 		wantOk    bool
// 	}{
// 		{"string", fields{settingsMap{"filename": "myfilename"}}, "filename", "myfilename", true},
// 		{"int", fields{settingsMap{"maxopenfiles": 10}}, "maxopenfiles", 10, true},
// 		{"bool", fields{settingsMap{"isDir": true}}, "isDir", true, true},
// 		//the following must fail
// 		{"int", fields{settingsMap{"maxopenfiles": 10}}, "notmaxopenfiles", nil, false},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &Settings{
// 				items: tt.fields.items,
// 			}
// 			gotValue, gotOk := s.Get(tt.args)
// 			if gotOk != tt.wantOk {
// 				t.Errorf("Settings.Get() gotOk = %v, want %v", gotOk, tt.wantOk)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotValue, tt.wantValue) {
// 				t.Errorf("Settings.Get() gotValue = %v, want %v", gotValue, tt.wantValue)
// 			}
// 		})
// 	}
// }

// func ExampleNewSettings() {
// 	s := NewSettings()
// 	printValue := func(key string) {
// 		value, ok := s.Get(key)
// 		if ok {
// 			fmt.Println(value)
// 		} else {
// 			fmt.Println("!")
// 		}
// 	}
// 	s.Set("FontName", "Manlo")
// 	s.Set("FontSize", 12)
// 	s.Set("FontItalic", true)
// 	printValue("FontName")
// 	printValue("FontSize")
// 	printValue("FontItalic")
// 	printValue("doesnotexist")
// 	//Output:
// 	//Manlo
// 	//12
// 	//true
// 	//!
// }
