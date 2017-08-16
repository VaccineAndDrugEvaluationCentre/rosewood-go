package rosewood

type ReportStatus int

const (
	Info ReportStatus = iota
	Echo
	Warning
	Error
	Fatal
)

//Settings implements a simple configuration solution.
type Settings struct {
	RangeOperator int32
	MandatoryCol  bool
	StyleSheet    string
	TableFileName string
	LogFileName   string
	Debug         bool
	RunMode       RunMode
	Report        func(string, ReportStatus)
}

//NewSettings returns an empty Settings struct
func NewSettings() *Settings {
	s := Settings{}
	return &s
}

//DefaultSettings returns default settings in case no settings were set.
func DefaultSettings() *Settings {
	settings := NewSettings()
	settings.RangeOperator = ':'
	settings.Debug = false
	// if ri.debug {
	// 	log.Printf("default settings loaded")
	// }
	return settings
}

// func (s Settings) String() string {
// 	return fmt.Sprintf("Settings:\n %#v", s)
// }

// //todo: change path to io.reader
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

// //todo: change path to io.writer
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
