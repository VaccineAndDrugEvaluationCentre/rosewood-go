// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.
// carpenter is reference implementation of the Rosewood language
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/drgo/fileutils"
	rosewood "github.com/drgo/rosewood/lib"
	_ "github.com/drgo/rosewood/renderers/html" //include needed renderers
)

const (
	ConfigFileBaseName = "carpenter.json"
)

//TODO:
//clean up debug and warnings: Debug=0 silent, 1=warnings only 2= verbose  3=internal debug info
// allow quoted argument in style command
// move all utilities to appropriate packages
// refresh vendor packages
// document new arguments
// add support for settings in package types
// clean-up all tests.
// use consistent errors types and constants eg NewError()
//add gracefull shutdown https://golang.org/pkg/os/signal/ along with a pointer to an optional cleanup function
// add support for processing subfolder if arg==./..
//?? add support for automerge; merged cells proceesed correctly even if there were no merge commands
// add word section struct to hold html inputfiles and section settings include headers and footers
// expand doInit to create word sections from input files in command line or current folder
// add command run to run an external file on the table; useful for formatting many similar tables

func main() {
	if err := RunApp(); err != nil {
		crash(err)
	}
}

//RunApp has all program logic; entry point for all tests
//WARNING: not safe to call concurrently; this is the only function allowed to change the job configuration
func RunApp() error {
	if len(os.Args) == 1 { //no command line arguments
		return DoFromConfigFile()
	}
	job, err := LoadConfigFromCommandLine()
	if err != nil {
		return err
	}
	if job.Settings.Debug == rosewood.DebugAll {
		fmt.Printf("current settings:\n%s\n", job)
	}
	switch job.Command { //TODO: check command is case insensitive
	case "do":
		if err = DoFromConfigFile(); err != nil {
			return err
		}
	case "check":
		job.Settings.CheckSyntaxOnly = true
		fallthrough
	case "run":
		job.Command = "process" //change to print nicer messages
		//FIXME: this check is not working
		if err = DoRun(job); rosewood.Errors().IsParsingError(err) {
			err = fmt.Errorf("one or more errors occurred during file processing") //do not report again
		}
	case "v1tov2":
		if err = V1toV2(job); err != nil { //shadwing err, so check/run/convert errors will not returned
			err = fmt.Errorf("one or more errors occurred during file processing")
		}
	case "init": //FIXME
		if err = DoInit(job); err != nil {
			err = fmt.Errorf("one or more errors occurred during configuration initialization")
		}
	case "version":
		fmt.Println(getVersion())
	case "help":
		helpMessage(FileDescriptorsToStrings(job.InputFiles), getVersion())
	default:
		helpMessage(nil, getVersion())
		return fmt.Errorf(ErrWrongCommand)
	}
	return err
}

func DoFromConfigFile() error {
	var (
		configFileName string
		err            error
	)
	if len(os.Args) == 1 { //only app name passed, use ConfigFileBaseName in current folder
		if configFileName, err = fileutils.GetFullPath(ConfigFileBaseName); err != nil {
			return err
		}
	} else {
		//  we must have been called with app name + do + a json file
		if strings.TrimSpace(strings.ToLower(os.Args[1])) != "do" {
			return fmt.Errorf("invalid command %s", os.Args[1])
		}
		if len(os.Args) < 3 {
			return fmt.Errorf("must specify a json configuration file")
		}
		configFileName = os.Args[2]
		if strings.ToLower(filepath.Ext(configFileName)) != ".json" {
			return fmt.Errorf("invalid config file name [%s] passed, ext must be json", configFileName)
		}
		if configFileName, err = fileutils.GetFullPath(configFileName); err != nil {
			return err
		}
	}
	//load configuration from config file
	configBuf, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return fmt.Errorf("failed to load configuration file: %v", err)
	}
	job := DefaultJob(rosewood.DefaultSettings())
	if err = json.Unmarshal(configBuf, job); err != nil {
		return fmt.Errorf("failed to parse configuration file %s: %v", configFileName, err)
	}
	if job.Settings.Debug >= rosewood.DebugUpdates {
		fmt.Println("configuration loaded from " + configFileName)
	}
	job.FileName = configFileName
	if job.Settings.Debug >= rosewood.DebugAll {
		fmt.Printf("current configuration: \n %s\n", job)
	}
	if err = DoRun(job); rosewood.Errors().IsParsingError(err) {
		err = fmt.Errorf("one or more errors occurred during file processing") //do not report again
	}
	return err
}

func LoadConfigFromCommandLine() (*Job, error) {
	job := DefaultJob(rosewood.DefaultSettings()) //TODO: ensure all defaults are reasonable
	flgSets, _ := setupCommandFlag(job)
	flg, err := ParseCommandLine(flgSets[0], flgSets[1:]...)
	if err != nil {
		return nil, err
	}
	job.Command = flg.Name()
	//TODO: validate command line inputs; use in RunfromConfigFile too?!
	for _, fileName := range flg.Args() {
		job.InputFiles = append(job.InputFiles, NewFileDescriptor(fileName))
	}
	return job, nil
}

// func interactive() {
// 	usage(-1)
// 	in := bufio.NewScanner(os.Stdin)
// 	out := os.Stdout
// 	echo := func(s string, status rosewood.ReportStatus) { //prints s to out followed by linefeed
// 		io.WriteString(out, s)
// 		io.WriteString(out, EOL)
// 	}
// 	settings := rosewood.DefaultSettings()
// 	settings.Report = echo
// 	rwi := rosewood.NewInterpreter(settings)
// 	for {
// 		trace.Printf("\n>>")
// 		if !in.Scan() || strings.ToLower(in.Text()) == "q" {
// 			return
// 		}
// 		err := rwi.Parse(strings.NewReader(in.Text()), "stdin")
// 		if err != nil {
// 			continue
// 		}
// 		//echo(cmdList[0].String())
// 		//		err = p.Run(cmdList)
// 		if err != nil {
// 			//			echo(fmt.Sprintf("runtime error: %s", err)) //show the first error only
// 			continue
// 		}
// 	}
// }
