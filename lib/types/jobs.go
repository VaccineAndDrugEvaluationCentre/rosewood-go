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
	Command             string
	Debug               int
	ExecutableVersion   string `json:"-"`
	RwFileNames         []string
	OverWriteOutputFile bool
	OutputFileName      string
	OutputFormat        string
	PreserveWorkFiles   bool
	SaveConvertedFile   bool
	RosewoodSettings    *RosewoodSettings
	StyleSheetName      string
	WorkDirName         string
	configFileName      string //if not empty, points to source config file that was used to load the config
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
	job.configFileName = FileName
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

//SaveToMDSon save job configuration to specified MDSon file
func (job *Job) SaveToMDSon(FileName string, overwrite bool) error {
	err := mdson.MarshalToFile(job, FileName, overwrite)
	if err != nil {
		return fmt.Errorf("failed to save job configuration: %v", err)
	}
	return nil
}

//ConfigFileName path to the MDSon file used to initialize this Job
//empty if a file was not used
func (job *Job) ConfigFileName() string {
	return job.configFileName
}
