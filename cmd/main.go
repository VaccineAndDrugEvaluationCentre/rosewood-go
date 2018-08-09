// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.
// carpenter is reference implementation of the Rosewood language
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/drgo/fileutils"
	rosewood "github.com/drgo/rosewood/lib"
	_ "github.com/drgo/rosewood/renderers/html" //include needed renderers
)

const (
	//ConfigFileBaseName default config file name
	ConfigFileBaseName = "carpenter.mdson"
)

//FIXME: create htmldocx config when invoked from command line
//FIXME: ensure workdirname is used to find rw files
//FIXME: ensure overwriteoutputfile flag is respected

//TODO:
//clean up debug and warnings: Debug=0 silent, 1=warnings only 2= verbose  3=internal debug info
// allow quoted argument in style command
// move all utilities to appropriate packages
// document new arguments
// add support for settings in package types
// clean-up all tests.
// use consistent errors types and constants eg NewError()
//add gracefull shutdown https://golang.org/pkg/os/signal/ along with a pointer to an optional cleanup function
// add support for processing subfolder if arg==./..
//?? add support for automerge; merged cells proceesed correctly even if there were no merge commands
// expand doInit to create word sections from input files in command line or current folder
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
		return DoFromConfigFile()
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
		if err = DoFromConfigFile(); err != nil {
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
		if _, err = DoInit(job); err != nil {
			err = fmt.Errorf("one or more errors occurred during configuration initialization: %s", err)
		}
	case "version":
		fmt.Println(getVersion())
	case "help":
		helpMessage(job.RwFileNames, getVersion())
	default:
		helpMessage(nil, getVersion())
		return fmt.Errorf(ErrWrongCommand)
	}
	return err
}

//DoFromConfigFile runs a job using a config file (not command line options)
func DoFromConfigFile() error {
	var (
		configFileName string
		err            error
	)
	if len(os.Args) == 1 { //only app name passed, use ConfigFileBaseName in current folder
		//FIXME: replace getfullpath with os.Abs() if os.Abs can handle empty argument
		if configFileName, err = fileutils.GetFullPath(ConfigFileBaseName); err != nil {
			return err
		}
	} else {
		//  we must have been called with app name + do + a json file
		if strings.TrimSpace(strings.ToLower(os.Args[1])) != "do" {
			return fmt.Errorf("invalid command %s", os.Args[1])
		}
		if len(os.Args) < 3 {
			return fmt.Errorf("must specify an MDSon configuration file")
		}
		configFileName = os.Args[2]
		if ext := strings.ToLower(filepath.Ext(configFileName)); ext != ".mdson" {
			return fmt.Errorf("invalid config file name [%s] passed", configFileName)
		}
		//FIXME: replace getfullpath with os.Abs()
		if configFileName, err = fileutils.GetFullPath(configFileName); err != nil {
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
	//TODO: validate command line inputs; use in RunfromConfigFile too?!
	for _, fileName := range flg.Args() {
		job.RwFileNames = append(job.RwFileNames, fileName)
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
