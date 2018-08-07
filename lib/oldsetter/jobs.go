package setter

import (
	"bytes"
	"encoding/json"
	"fmt"
)

//Job holds info about the current run
type Job struct { //IMMUTABLE once initialized; TODO: hide internal details using getter functions
	FileName            string //if not empty, points to source config file that was used to load the config
	Command             string
	Debug               int
	ExecutableVersion   string `json:"-"`
	InputFileNames      []string
	OverWriteOutputFile bool
	OutputFileName      string
	OutputFormat        string
	PreserveWorkFiles   bool
	SaveConvertedFile   bool
	RosewoodSettings    *RosewoodSettings
	StyleSheetName      string
	WorkDirName         string
}

//DefaultJob returns default job
func DefaultJob(settings *RosewoodSettings) *Job {
	return &Job{
		Command: "run",
		//FIXME: update this via ? option if needed.
		// ExecutableVersion: Version,
		RosewoodSettings: settings,
	}
}

//GetNameOfInputFile returns name of input file corresponding to index
func (job *Job) GetNameOfInputFile(index int) string {
	if index < 0 || len(job.InputFileNames) >= index {
		return ""
	}
	return job.InputFileNames[index]
}

func (job Job) String() string {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(job); err != nil {
		panic(fmt.Sprintf("failed to print job configuration: %v", err))
	}
	return buf.String()
}

//GetValidSettings loads configuration from config file if one exists in current dir,
//otherwise returns default Rosewood settings
// func GetValidSettings() (*Settings, error) {
// 	configFileName, err := os.Getwd()
// 	if err != nil {
// 		return nil, err
// 	}
// 	configFileName = filepath.Join(configFileName, ConfigFileBaseName)
// 	configBuf, err := ioutil.ReadFile(configFileName)
// 	if os.IsNotExist(err) { //file not found
// 		return rosewood.DefaultSettings(), nil
// 	}
// 	settings := new(Settings)
// 	if err != nil { //some error other than file does not exist occurred
// 		return nil, fmt.Errorf("failed to load configuration file %s: %v", configFileName, err)
// 	}
// 	if err = json.Unmarshal(configBuf, settings); err != nil {
// 		return nil, fmt.Errorf("failed to parse settings: %v", err)
// 	}
// 	fmt.Println("configuration loaded from " + configFileName)

// 	return settings, nil
// }

// type FileDescriptor struct {
// 	Name    string
// 	Version string `json:",omitempty"`
// }

// func NewFileDescriptor(fileName string) string {
// 	return &FileDescriptor{
// 		Name:    fileName,
// 		Version: "",
// 	}
// }

// func FileDescriptorsToStrings(fds []string) string
// 	for _, f := range fds {
// 		ss = append(ss, f.Name)
// 	}
// 	return ss
// }

//GetValidJob loads configuration from a config file if one exists in current dir,
// //otherwise returns default Job including default Rosewood settings
// func GetValidJob(configFileName string) (*Job, error) {
// 	// configFileName = filepath.Join(configFileName, ConfigFileBaseName)
// 	configBuf, err := ioutil.ReadFile(configFileName)
// 	if os.IsNotExist(err) { //configFile not found
// 		return DefaultJob(rosewood.DefaultSettings()), nil
// 	}
// 	if err != nil { //some error other than file does not exist occurred
// 		return nil, fmt.Errorf("failed to load configuration file %s: %v", configFileName, err)
// 	}
// 	job := new(Job)
// 	if err = json.Unmarshal(configBuf, job); err != nil {
// 		return nil, fmt.Errorf("failed to parse settings: %v", err)
// 	}
// 	fmt.Println("configuration loaded from " + configFileName)
// 	job.FileName = configFileName
// 	return job, nil
// }
