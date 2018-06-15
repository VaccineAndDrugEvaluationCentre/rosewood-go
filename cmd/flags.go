package main

import (
	"flag"
	"fmt"

	rosewood "github.com/drgo/rosewood/lib"
)

//WARNING: the default value field in flag struct is not used
func setupCommandFlag(job *Job) (flgSets []*flag.FlagSet, err error) {
	globals := NewCommand("", []Flag{
		{&job.Settings.Debug, "debug", "d", nil},
		{&job.Settings.MaxConcurrentWorkers, "max-threads", "mt", nil},
	})
	cmdDo := NewCommand("do", []Flag{})
	cmdRun := NewCommand("run", []Flag{
		{&job.Settings.ConvertOldVersions, "convert-old", "co", nil},
		{&job.Settings.ConvertFromVersion, "rosewood-version", "rv", nil},
		{&job.Settings.DoNotInlineCSS, "no-inlined-css", "", nil},
		{&job.OutputFile.Name, "output", "o", nil},
		{&job.Settings.OverWriteOutputFile, "replace", "r", nil},
		{&job.Settings.PreserveWorkFiles, "keep-temp", "k", nil},
		{&job.Settings.SaveConvertedFile, "save-converted", "sc", nil},
		{&job.Settings.StyleSheetName, "style", "s", nil},
		{&job.WorkDirName, "work-dir", "w", nil},
	})
	cmdCheck := NewCommand("check", []Flag{})
	cmdV1tov2 := NewCommand("v1tov2", []Flag{
		{&job.Settings.ConvertFromVersion, "rosewood-version", "rv", nil},
		{&job.Settings.OverWriteOutputFile, "replace", "r", nil},
	})
	cmdInit := NewCommand("init", []Flag{
		{&job.Settings.OverWriteOutputFile, "replace", "r", nil},
	})
	cmdHelp := NewCommand("help", []Flag{})
	cmdVersion := NewCommand("version", []Flag{})
	flgSets = append(flgSets, globals, cmdDo, cmdRun, cmdCheck, cmdV1tov2, cmdHelp, cmdVersion, cmdInit)
	for _, fs := range flgSets {
		fs.Usage = func() {}    //disable internal usage function
		fs.SetOutput(devNull{}) //suppress output from package flag
	}
	return flgSets, nil
}

func getVersion() string {
	return fmt.Sprintf(versionMessage, Version, Build, rosewood.Version)
}
