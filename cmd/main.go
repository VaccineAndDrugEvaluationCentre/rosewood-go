// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.
// carpenter is reference implementation of the Rosewood language
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/drgo/core/args"
	"github.com/drgo/core/files"
	"github.com/drgo/core/ui"
	rosewood "github.com/drgo/rosewood/lib"
	_ "github.com/drgo/rosewood/renderers/html" //include needed renderers
)

const (
	//ConfigFileBaseName default config file name
	ConfigFileBaseName = "carpenter.mdson"
)

//FIXME: ensure workdirname is used to find rw files

//TODO:
//clean up debug and warnings: Debug=0 silent, 1=warnings only 2= verbose  3=internal debug info
// allow quoted argument in style command
// move all utilities to appropriate packages
// add support for settings in package types
// clean-up all tests.
// use consistent errors types and constants eg NewError()
//add gracefull shutdown https://golang.org/pkg/os/signal/ along with a pointer to an optional cleanup function
// add support for processing subfolder if arg==./..
// add command run to run an external rw file on the table; useful for formatting many similar tables

var (
	exeName string
	ux      ui.UI
)

func main() {
	exeName = os.Args[0]
	if err := RunApp(); err != nil {
		crash(err)
	}
}

//RunApp entry point for all tests
func RunApp() error {
	var (
		job *rosewood.Job
		err error
	)
	if len(os.Args) == 1 { //no command line arguments
		job, err = LoadConfigFromFile("")
	} else {
		job, err = LoadConfigFromCommandLine()
	}
	if err != nil {
		return err
	}
	ux = ui.NewUI(job.RosewoodSettings.Debug)
	return RunJob(job)
}

//RunJob the workforce
//WARNING: not thread-safe; this is the only function allowed to change the job configuration
func RunJob(job *rosewood.Job) error {
	var err error
	ux.Log("Started on ", time.Now())
	ux.Log("current settings:", job)
	switch job.Command {
	case "do":
		if len(job.RwFileNames) == 0 {
			return fmt.Errorf("must specify an MDSon configuration file")
		}
		job, err = LoadConfigFromFile(job.RwFileNames[0])
		if err != nil {
			return err
		}
		return RunJob(job)
	case "check":
		job.RosewoodSettings.CheckSyntaxOnly = true
		fallthrough
	case "run":
		job.Command = "process" //change to print nicer messages
		//FIXME: this check is not working
		if err = DoRun(job); rosewood.Errors().IsParsingError(err) {
			err = fmt.Errorf("one or more errors occurred during file processing") //do not report again
		}
	case "v1tov2":
		if err = V1toV2(job); err != nil {
			err = fmt.Errorf("one or more errors occurred during file processing: %s", err)
		}
	case "init":
		configFilename, err := DoInit(job)
		if err == nil && job.RosewoodSettings.Debug >= rosewood.DebugUpdates {
			fmt.Printf("configuration saved as '%s'\n", configFilename)
		}
		return err
	case "version":
		fmt.Println(getVersion())
	case "help", "h":
		helpMessage(job.RwFileNames, getVersion())
	default:
		helpMessage(nil, getVersion())
		return fmt.Errorf(ErrWrongCommand, exeName)
	}
	return err
}

//LoadConfigFromFile loads a job from a config file
func LoadConfigFromFile(configFileName string) (job *rosewood.Job, err error) {
	if configFileName == "" {
		configFileName = ConfigFileBaseName
	}
	if configFileName, err = files.GetFullPath(ConfigFileBaseName); err != nil {
		return nil, err
	}
	job = rosewood.DefaultJob(rosewood.DefaultSettings())
	if err = job.LoadFromMDSonFile(configFileName); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file %s: %v", configFileName, err)
	}
	job.Command = "run" // config file jobs are equivalent to calling carpenter run
	return job, nil
}

//LoadConfigFromCommandLine creates a job object using command line arguments
func LoadConfigFromCommandLine() (*rosewood.Job, error) {
	//TODO: ensure all defaults are reasonable
	job := rosewood.DefaultJob(rosewood.DefaultSettings())
	flgSets, _ := setupCommandFlag(job)
	flg, err := args.ParseCommandLine(flgSets[0], flgSets[1:]...)

	if err != nil {
		return nil, err
	}
	job.Command = flg.Name()
	if len(flg.Args()) == 0 {
		return job, nil
	}
	switch runtime.GOOS {
	case "windows":
		job.RwFileNames, err = WinArgsToFileNames(flg.Args()[0])
		if err != nil {
			return nil, err
		}
	default:
		for _, fileName := range flg.Args() {
			job.RwFileNames = append(job.RwFileNames, fileName)
		}
	}
	return job, nil
}

func WinArgsToFileNames(args string) ([]string, error) {
	if !strings.ContainsAny(args, "*?") {
		return []string{args}, nil
	}
	return filepath.Glob(args)
}
