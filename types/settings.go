package types

//RosewoodSettings for controlling Rosewood lib
type RosewoodSettings struct {
	CheckSyntaxOnly    bool   `mdson:"-"`
	ColumnSeparator    string `mdson:"-"`
	ConvertOldVersions bool
	ConvertFromVersion string
	//controls printing debug info by internal lib routines
	Debug                int
	DoNotInlineCSS       bool
	MandatoryCol         bool   `mdson:"-"`
	MarkdownRender       string //"disabled", "strict", "standard"
	MaxConcurrentWorkers int
	// PreserveWorkFiles    bool
	RangeOperator     int32 `mdson:"-"`
	ReportAllError    bool
	SaveConvertedFile bool
	SectionCapacity   int    `mdson:"-"`
	SectionSeparator  string `mdson:"-"`
	SectionsPerTable  int    `mdson:"-"`
	StyleSheetName    string
	TrimCellContents  bool
}

//NewRosewoodSettings returns an empty Settings struct
func NewRosewoodSettings() *RosewoodSettings {
	return &RosewoodSettings{}
}

//DefaultRosewoodSettings returns default settings in case no settings were set.
func DefaultRosewoodSettings() *RosewoodSettings {
	settings := NewRosewoodSettings()
	settings.MarkdownRender = "strict"
	settings.SectionsPerTable = 4
	settings.SectionCapacity = 100
	settings.SectionSeparator = "+++"
	settings.ColumnSeparator = "|"
	settings.RangeOperator = ':'
	settings.MaxConcurrentWorkers = 24
	return settings
}

//DebugRosewoodSettings returns default settings for settings and setup tracing
func DebugRosewoodSettings(debug int) *RosewoodSettings {
	settings := DefaultRosewoodSettings()
	settings.Debug = debug
	return settings
}
