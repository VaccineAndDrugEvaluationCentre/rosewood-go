// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package setter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//TODO: split Settings into Rosewood settings and carpenter settings to be included in a jobconfig struct
/*
	jobconfig
	rosewood settings
	outputfilename a filename descriptor
	inputfilenames array of filename descriptors
	others eg workingdir etc
*/

//Settings implements a simple configuration solution.
// `json:"-"` used to exclude certain fields from saving/loading from config files
type Settings struct {
	CheckSyntaxOnly    bool   `json:"-"`
	ColumnSeparator    string `json:"-"`
	ConvertOldVersions bool
	ConvertFromVersion string
	Debug              int
	DoNotInlineCSS     bool
	//	ExecutableVersion    string `json:"-"`
	LibVersion           string `json:"-"`
	MandatoryCol         bool
	MaxConcurrentWorkers int
	//OverWriteOutputFile  bool
	// OutputFileName       string
	// OutputFormat         string
	PreserveWorkFiles bool
	RangeOperator     int32 `json:"-"`
	ReportAllError    bool
	SaveConvertedFile bool
	SectionCapacity   int    `json:"-"`
	SectionSeparator  string `json:"-"`
	SectionsPerTable  int    `json:"-"`
	StyleSheetName    string
	// WorkDirName          string
	TrimCellContents bool
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
	settings.MaxConcurrentWorkers = 24
	return settings
}

//DebugSettings returns default settings for settings and setup tracing
func DebugSettings(Tracing int) *Settings {
	settings := DefaultSettings()
	settings.Debug = Tracing
	return settings
}

//TODO: change path to io.reader
func LoadSettings(path string) (*Settings, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load settings: %v", err)
	}
	var s Settings
	err = json.Unmarshal(file, &s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse settings: %v", err)
	}
	return &s, nil
}

func (settings *Settings) String() string {
	b, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error printing settings %s", err)
	}
	return string(b)
}

//TODO: replace with tempfile with proper close error management
func SaveSettings(settings *Settings, path string, replace bool) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to save settings: %v", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(settings); err != nil {
		return fmt.Errorf("failed to save settings: %v", err)
	}
	return nil
}
