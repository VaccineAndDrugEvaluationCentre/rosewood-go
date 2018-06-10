package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	rosewood "github.com/drgo/rosewood/lib"
	"github.com/drgo/rosewood/lib/setter"
)

//Job holds info about the current run
type Job struct {
	Settings            *setter.Settings
	ExecutableVersion   string `json:"-"`
	OverWriteOutputFile bool
	OutputFileName      string
	OutputFormat        string
	PreserveWorkFiles   bool
	SaveConvertedFile   bool
	StyleSheetName      string
	WorkDirName         string
	InputFileNames      []string
}

func DefaultJob(settings *setter.Settings) *Job {
	return &Job{
		Settings:          settings,
		ExecutableVersion: Version,
	}
}

//GetValidJob loads configuration from a config file if one exists in current dir,
//otherwise returns default Job including default Rosewood settings
func GetValidJob() (*Job, error) {
	configFileName, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configFileName = filepath.Join(configFileName, ConfigFileBaseName)
	configBuf, err := ioutil.ReadFile(configFileName)
	configFileFound := !os.IsNotExist(err)
	var job *Job
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
	settings := new(setter.Settings)
	return job, nil
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
