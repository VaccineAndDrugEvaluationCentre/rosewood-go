package main

import (
	"flag"
	"fmt"

	rosewood "github.com/drgo/rosewood/lib"
)

//WARNING: the default value field in flag struct is not used
func setupCommandFlag(job *rosewood.Job) (flgSets []*flag.FlagSet, err error) {
	globals := NewCommand("", []Flag{
		{&job.RosewoodSettings.Debug, "debug", "d", nil},
		{&job.RosewoodSettings.MaxConcurrentWorkers, "max-threads", "mt", nil},
	})
	baseflags := []Flag{ //common to several flags
		{&job.OverWriteOutputFile, "replace", "r", nil},
	}
	cmdDo := NewCommand("do", []Flag{})
	runFlags := append(baseflags, []Flag{
		{&job.ConfigFileName, "config", "cfg", nil},
		{&job.RosewoodSettings.ConvertOldVersions, "convert-old", "co", nil},
		{&job.RosewoodSettings.ConvertFromVersion, "rosewood-version", "rv", nil},
		{&job.RosewoodSettings.DoNotInlineCSS, "no-inlined-css", "", nil},
		{&job.OutputFileName, "output", "o", nil},
		{&job.RosewoodSettings.PreserveWorkFiles, "keep-temp", "k", nil},
		{&job.RosewoodSettings.StyleSheetName, "style", "s", nil},
		{&job.WorkDirName, "work-dir", "w", nil},
	}...)
	cmdRun := NewCommand("run", runFlags)
	cmdCheck := NewCommand("check", []Flag{})
	cmdV1tov2 := NewCommand("v1tov2", append(baseflags, []Flag{
		{&job.RosewoodSettings.ConvertFromVersion, "rosewood-version", "rv", nil},
	}...))
	cmdInit := NewCommand("init", runFlags)
	cmdHelp := NewCommand("help", []Flag{})
	cmdQuery := NewCommand("h", []Flag{})
	cmdVersion := NewCommand("version", []Flag{})
	flgSets = append(flgSets, globals, cmdDo, cmdRun, cmdCheck, cmdV1tov2, cmdHelp, cmdQuery, cmdVersion, cmdInit)
	for _, fs := range flgSets {
		fs.Usage = nil          //disable internal usage function
		fs.SetOutput(devNull{}) //suppress output from package flag
	}
	return flgSets, nil
}

func getVersion() string {
	return fmt.Sprintf(versionMessage, Version, Build, rosewood.LibVersion())
}
