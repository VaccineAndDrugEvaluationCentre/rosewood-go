package types

//RosewoodSettings for controlling Rosewood lib
type RosewoodSettings struct {
	CheckSyntaxOnly    bool   `json:"-"`
	ColumnSeparator    string `json:"-"`
	ConvertOldVersions bool
	ConvertFromVersion string
	//controls printing debug info by internal lib routines
	Debug                int
	DoNotInlineCSS       bool
	MandatoryCol         bool
	MaxConcurrentWorkers int
	//	OverWriteOutputFile  bool
	// OutputFileName       string
	// OutputFormat         string
	//FIXME: move to options?
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
func NewRosewoodSettings() *RosewoodSettings {
	return &RosewoodSettings{}
}

//DefaultSettings returns default settings in case no settings were set.
func DefaultRosewoodSettings() *RosewoodSettings {
	settings := NewRosewoodSettings()
	settings.SectionsPerTable = 4
	settings.SectionCapacity = 100
	settings.SectionSeparator = "+++"
	settings.ColumnSeparator = "|"
	settings.RangeOperator = ':'
	settings.MaxConcurrentWorkers = 24
	return settings
}

//FIXME: remove tracing to runOptions

//DebugSettings returns default settings for settings and setup tracing
func DebugRosewoodSettings(Tracing int) *RosewoodSettings {
	settings := DefaultRosewoodSettings()
	settings.Debug = Tracing
	return settings
}

const configFileBaseName = "htmldocx.json"
const scriptFileBaseName = "script.htds"

//Options holds info on a single run
type Options struct {
	Debug                 int
	JobFileName           string
	DefaultConfigFileName string
	DefaultScriptFileName string
	ExecutableVersion     string `json:"-"`
	LibVersion            string `json:"-"`
}

//DefaultOptions returns a default option setting
func DefaultOptions() *Options {
	return &Options{
		Debug: 3,
		DefaultConfigFileName: configFileBaseName,
		DefaultScriptFileName: scriptFileBaseName,
	}
}

//SetDebug sets the debug level
func (opt *Options) SetDebug(level int) *Options {
	opt.Debug = level
	return opt
}

//FIXME: remove once new tracing is implemented
const (
	//DebugSilent print errors only
	DebugSilent int = iota
	//DebugWarning print warnings and errors
	DebugWarning
	//DebugUpdates print execution updates, warnings and errors
	DebugUpdates
	//DebugAll print internal debug messages, execution updates, warnings and errors
	DebugAll
)
