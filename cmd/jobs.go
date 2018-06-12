package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	rosewood "github.com/drgo/rosewood/lib"
	"github.com/drgo/rosewood/lib/setter"
)

type task struct {
	settings   *setter.Settings
	inputFile  *FileDescriptor
	outputFile *FileDescriptor
}

func (task *task) GetResult(err error) result {
	return result{
		task:  task,
		error: err,
	}
}

type result struct {
	task  *task
	error error
}

func NewResult(task *task, err error) *result {
	return &result{
		task:  task,
		error: err,
	}
}

func (result *result) SetError(err error) *result {
	result.error = err
	return result
}

type FileDescriptor struct {
	Name    string
	Version string `json:",omitempty"`
}

func NewFileDescriptor(fileName string) *FileDescriptor {
	return &FileDescriptor{
		Name:    fileName,
		Version: "",
	}
}

//Job holds info about the current run
type Job struct { //IMMUTABLE once initialized; TODO: hide internal details but getter functions
	FileName            string //if not empty, points to source config file that was used to load the config
	Command             string
	ExecutableVersion   string `json:"-"`
	InputFiles          []*FileDescriptor
	OverWriteOutputFile bool
	OutputFile          *FileDescriptor
	OutputFormat        string
	PreserveWorkFiles   bool
	SaveConvertedFile   bool
	Settings            *setter.Settings
	StyleSheetName      string
	WorkDirName         string
}

func DefaultJob(settings *setter.Settings) *Job {
	dir, _ := os.Getwd()
	return &Job{
		Command:           "run",
		ExecutableVersion: Version,
		OutputFile:        NewFileDescriptor(""),
		Settings:          settings,
		WorkDirName:       dir,
	}
}

//GetValidJob loads configuration from a config file if one exists in current dir,
//otherwise returns default Job including default Rosewood settings
func GetValidJob(configFileName string) (*Job, error) {
	// configFileName, err := os.Getwd()
	// if err != nil {
	// 	return nil, err
	// }
	// configFileName = filepath.Join(configFileName, ConfigFileBaseName)
	configBuf, err := ioutil.ReadFile(configFileName)
	job := new(Job)
	if !os.IsNotExist(err) { //configFile found
		if err != nil { //some error other than file does not exist occurred
			return nil, fmt.Errorf("failed to load configuration file %s: %v", configFileName, err)
		}
		err = json.Unmarshal(configBuf, job)
		if err != nil {
			return nil, fmt.Errorf("failed to parse settings: %v", err)
		}
		fmt.Println("configuration loaded from " + configFileName)
		job.FileName = configFileName
	} else {
		job = DefaultJob(rosewood.DefaultSettings())
	}
	return job, nil
}

func (job *Job) GetTask(inputFile, outputFile *FileDescriptor) *task {
	return &task{
		settings:   job.Settings,
		inputFile:  inputFile,
		outputFile: outputFile,
	}
}

func (job *Job) GetDefaultTask() *task {
	return &task{
		settings: job.Settings,
	}
}

func (job *Job) GetNameOfInputFile(index int) string {
	if index < 0 || len(job.InputFiles) >= index {
		return ""
	}
	return job.InputFiles[index].Name
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

//TODO: replace with tempfile with proper close error management
func SaveJob(job *Job, path string, replace bool) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to save job configuration: %v", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(job); err != nil {
		return fmt.Errorf("failed to save job configuration: %v", err)
	}
	return nil
}

//GetValidSettings loads configuration from config file if one exists in current dir,
//otherwise returns default Rosewood settings
func GetValidSettings() (*setter.Settings, error) {
	configFileName, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configFileName = filepath.Join(configFileName, ConfigFileBaseName)
	configBuf, err := ioutil.ReadFile(configFileName)
	configFileFound := !os.IsNotExist(err)
	settings := new(setter.Settings)
	if configFileFound {
		if err != nil { //some error other than file does not exist occurred
			return nil, fmt.Errorf("failed to load configuration file %s: %v", configFileName, err)
		}
		err = json.Unmarshal(configBuf, settings)
		if err != nil {
			return nil, fmt.Errorf("failed to parse settings: %v", err)
		}
		fmt.Println("configuration loaded from " + configFileName)
	} else {
		settings = rosewood.DefaultSettings()
	}
	return settings, nil
}
