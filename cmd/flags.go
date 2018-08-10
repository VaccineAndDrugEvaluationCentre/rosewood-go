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
	cmdDo := NewCommand("do", []Flag{})
	cmdRun := NewCommand("run", []Flag{
		{&job.ConfigFileName, "config", "cfg", nil},
		{&job.RosewoodSettings.ConvertOldVersions, "convert-old", "co", nil},
		{&job.RosewoodSettings.ConvertFromVersion, "rosewood-version", "rv", nil},
		{&job.RosewoodSettings.DoNotInlineCSS, "no-inlined-css", "", nil},
		{&job.OutputFileName, "output", "o", nil},
		//FIXME: read overwriteoutputfile and PreserveWorkFiles into options
		{&job.OverWriteOutputFile, "replace", "r", nil},
		{&job.RosewoodSettings.PreserveWorkFiles, "keep-temp", "k", nil},
		{&job.RosewoodSettings.StyleSheetName, "style", "s", nil},
		{&job.WorkDirName, "work-dir", "w", nil},
	})
	cmdCheck := NewCommand("check", []Flag{})
	cmdV1tov2 := NewCommand("v1tov2", []Flag{
		{&job.RosewoodSettings.ConvertFromVersion, "rosewood-version", "rv", nil},
		{&job.OverWriteOutputFile, "replace", "r", nil},
	})
	cmdInit := NewCommand("init", []Flag{
		//FIXME: read overwriteoutputfile and PreserveWorkFiles into options
		{&job.ConfigFileName, "config", "cfg", nil},
		{&job.RosewoodSettings.ConvertOldVersions, "convert-old", "co", nil},
		{&job.RosewoodSettings.ConvertFromVersion, "rosewood-version", "rv", nil},
		{&job.RosewoodSettings.DoNotInlineCSS, "no-inlined-css", "", nil},
		{&job.OutputFileName, "output", "o", nil},
		{&job.OverWriteOutputFile, "replace", "r", nil},
		{&job.RosewoodSettings.PreserveWorkFiles, "keep-temp", "k", nil},
		{&job.RosewoodSettings.SaveConvertedFile, "save-converted", "sc", nil},
		{&job.RosewoodSettings.StyleSheetName, "style", "s", nil},
		{&job.WorkDirName, "work-dir", "w", nil},
	})
	cmdHelp := NewCommand("help", []Flag{})
	cmdQuery := NewCommand("h", []Flag{})
	cmdVersion := NewCommand("version", []Flag{})
	flgSets = append(flgSets, globals, cmdDo, cmdRun, cmdCheck, cmdV1tov2, cmdHelp, cmdQuery, cmdVersion, cmdInit)
	for _, fs := range flgSets {
		fs.Usage = func() {}    //disable internal usage function
		fs.SetOutput(devNull{}) //suppress output from package flag
	}
	return flgSets, nil
}

func getVersion() string {
	return fmt.Sprintf(versionMessage, Version, Build, rosewood.LibVersion())
}
