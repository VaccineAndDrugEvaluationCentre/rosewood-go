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
	MandatoryCol         bool `mdson:"-"`
	MaxConcurrentWorkers int
	//FIXME: move to options?
	PreserveWorkFiles bool
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
	settings.SectionsPerTable = 4
	settings.SectionCapacity = 100
	settings.SectionSeparator = "+++"
	settings.ColumnSeparator = "|"
	settings.RangeOperator = ':'
	settings.MaxConcurrentWorkers = 24
	return settings
}

//FIXME: remove tracing to runOptions

//DebugRosewoodSettings returns default settings for settings and setup tracing
func DebugRosewoodSettings(debug int) *RosewoodSettings {
	settings := DefaultRosewoodSettings()
	settings.Debug = debug
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
	ExecutableVersion     string `mdson:"-"`
	LibVersion            string `mdson:"-"`
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
