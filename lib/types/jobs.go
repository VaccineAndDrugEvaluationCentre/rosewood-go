package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/drgo/mdson"
)

//Job holds info about the current run
type Job struct { //IMMUTABLE once initialized; TODO: hide internal details using getter functions
	Command             string `mdson:"-"`
	Debug               int
	ExecutableVersion   string `mdson:"-"`
	RwFileNames         []string
	OverWriteOutputFile bool
	OutputFileName      string
	OutputFormat        string
	PreserveWorkFiles   bool
	SaveConvertedFile   bool
	RosewoodSettings    *RosewoodSettings
	StyleSheetName      string
	WorkDirName         string
	HTMLFileNames       []string `mdson:"-"` //files generated from Rosewood scripts
	HTMLWorkDirName     string   `mdson:"-"` //dir where HTMLFileNames are stored
	ConfigFileName      string   `mdson:"-"` //MDSon file that was used to load the config
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
	if index < 0 || len(job.RwFileNames) >= index {
		return ""
	}
	return job.RwFileNames[index]
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

//LoadFromMDSonFile loads job confgiuration from a file
func (job *Job) LoadFromMDSonFile(FileName string) error {
	configFile, err := os.Open(FileName)
	if err != nil {
		return fmt.Errorf("failed to load configuration file: %v", err)
	}
	defer configFile.Close()
	if err = job.LoadFromMDSon(configFile); err != nil {
		return fmt.Errorf("failed to parse configuration file %s: %v", FileName, err)
	}
	job.ConfigFileName = FileName
	if job.RosewoodSettings.Debug >= DebugUpdates {
		fmt.Println("configuration loaded from " + FileName)
	}
	return nil
}

//LoadFromMDSon loads job confgiuration from an io.Reader
func (job *Job) LoadFromMDSon(r io.Reader) error {
	mdson.SetDebug(job.Debug)
	err := mdson.Unmarshal(r, job)
	if err != nil {
		return err
	}
	return nil
}

//SaveToMDSon saves job configuration to specified MDSon writer
func (job *Job) SaveToMDSon(w io.Writer) error {
	buf, err := mdson.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to save job configuration: %s", err)
	}
	_, err = w.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to save job configuration: %s", err)
	}
	return nil
}

//SaveToMDSonFile save job configuration to specified MDSon file
func (job *Job) SaveToMDSonFile(FileName string, overwrite bool) error {
	err := mdson.MarshalToFile(job, FileName, overwrite)
	if err != nil {
		return fmt.Errorf("failed to save job configuration: %v", err)
	}
	return nil
}

// //ConfigFileName returns path to the MDSon file used to initialize this Job
// //empty if a file was not used
// func (job *Job) ConfigFileName() string {
// 	return job.ConfigFileName
// }

// //SetConfigFileName sets the MDSon file name
// func (job *Job) SetConfigFileName(fileName string) {
// 	job.ConfigFileName = fileName
// }
