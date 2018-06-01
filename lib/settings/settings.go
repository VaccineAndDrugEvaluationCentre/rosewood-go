// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package settings

import (
	"encoding/json"
	"fmt"
)

//Settings implements a simple configuration solution.
type Settings struct {
	CheckSyntaxOnly     bool
	ColumnSeparator     string
	ConvertOldVersions  bool
	ConvertFromVersion  string
	Debug               int
	DoNotInlineCSS      bool
	ExecutableVersion   string
	LibVersion          string
	InputFileName       string
	MandatoryCol        bool
	OverWriteOutputFile bool
	OutputFileName      string
	RangeOperator       int32
	ReportAllError      bool
	SaveConvertedFile   bool
	SectionCapacity     int
	SectionSeparator    string
	SectionsPerTable    int
	StyleSheetName      string
	TrimCellContents    bool
}

//NewSettings returns an empty Settings struct
func NewSettings() *Settings {
	return &Settings{}
}

//DefaultSettings returns default settings in case no settings were set.
func DefaultSettings() *Settings {
	settings := NewSettings()
	settings.SectionsPerTable = 4
	settings.SectionCapacity = 100
	settings.SectionSeparator = "+++"
	settings.ColumnSeparator = "|"
	settings.RangeOperator = ':'
	return settings
}

//DebugSettings returns default settings for settings and setup tracing
func DebugSettings(Tracing int) *Settings {
	settings := DefaultSettings()
	settings.Debug = Tracing
	return settings
}

func (settings *Settings) String() string {
	b, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error printing settings %s", err)
	}
	return string(b)
}

// func (s Settings) String() string {
// 	return fmt.Sprintf("Settings:\n %#v", s)
// }

// //TODO: change path to io.reader
// func (s *Settings) LoadSettings(path string) error {
// 	file, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		return fmt.Errorf("failed to load settings: %v", err)
// 	}
// 	err = json.Unmarshal(file, &s.items)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse settings: %v", err)
// 	}
// 	return nil
// }

// //TODO: change path to io.writer
// func (s *Settings) SaveSettings(path string, replace true) error {
// 	file, err := os.Create(path)
// 	if err != nil {
// 		return fmt.Errorf("failed to save settings: %v", err)
// 	}
// 	defer file.Close()
// 	e := json.NewEncoder(file)
// 	if err := e.Encode(s.items); err != nil {
// 		return fmt.Errorf("failed to save settings: %v", err)
// 	}
// 	return nil
// }
