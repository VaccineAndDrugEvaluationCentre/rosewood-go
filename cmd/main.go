// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.
// carpenter is reference implementation of the Rosewood language
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/drgo/fileutils"
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

func main() {
	if err := RunApp(); err != nil {
		crash(err)
	}
}

//RunApp has all program logic; entry point for all tests
//WARNING: not thread-safe; this is the only function allowed to change the job configuration
func RunApp() error {
	if len(os.Args) == 1 { //no command line arguments
		return DoFromConfigFile("")
	}
	job, err := LoadConfigFromCommandLine()
	if err != nil {
		return err
	}
	if job.RosewoodSettings.Debug == rosewood.DebugAll {
		fmt.Printf("current settings:\n%s\n", job)
	}
	switch job.Command { //TODO: check command is case insensitive
	case "do":
		if len(job.RwFileNames) == 0 {
			return fmt.Errorf("must specify an MDSon configuration file")
		}
		if err = DoFromConfigFile(job.RwFileNames[0]); err != nil {
			return err
		}
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
		return fmt.Errorf(ErrWrongCommand)
	}
	return err
}

//DoFromConfigFile runs a job using a config file (not command line options)
func DoFromConfigFile(configFileName string) error {
	var err error
	if configFileName == "" {
		if configFileName, err = fileutils.GetFullPath(ConfigFileBaseName); err != nil {
			return err
		}
	}
	//load configuration from config file
	job := rosewood.DefaultJob(rosewood.DefaultSettings())
	if err = job.LoadFromMDSonFile(configFileName); err != nil {
		return fmt.Errorf("failed to parse configuration file %s: %v", configFileName, err)
	}
	if job.RosewoodSettings.Debug >= rosewood.DebugAll {
		fmt.Printf("current configuration: \n %s\n", job)
	}
	if err = DoRun(job); rosewood.Errors().IsParsingError(err) {
		err = fmt.Errorf("one or more errors occurred during file processing") //do not report again
	}
	return err
}

//LoadConfigFromCommandLine creates a object using command line arguments
func LoadConfigFromCommandLine() (*rosewood.Job, error) {
	job := rosewood.DefaultJob(rosewood.DefaultSettings()) //TODO: ensure all defaults are reasonable
	flgSets, _ := setupCommandFlag(job)
	flg, err := ParseCommandLine(flgSets[0], flgSets[1:]...)

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
