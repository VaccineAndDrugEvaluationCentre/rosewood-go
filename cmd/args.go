package main

import (
	"flag"
	"fmt"

	rosewood "github.com/drgo/rosewood/lib"
)

//WARNING: the default value field in flag struct is not used
func setupCommandFlag(job *Job) (flgSets []*flag.FlagSet, err error) {
	globals := NewCommand("", []Flag{
		{&job.Settings.Debug, "debug", "d", 0},
		{&job.Settings.MaxConcurrentWorkers, "max-threads", "mt", 0},
	})
	cmdDo := NewCommand("do", []Flag{})
	cmdRun := NewCommand("run", []Flag{
		{&job.Settings.ConvertOldVersions, "convert-old", "co", false},
		{&job.Settings.ConvertFromVersion, "rosewood-version", "rv", job.Settings.ConvertFromVersion},
		{&job.Settings.DoNotInlineCSS, "no-inlined-css", "", false},
		{&job.OutputFile.Name, "output", "o", ""},
		{&job.Settings.OverWriteOutputFile, "replace", "r", false},
		{&job.Settings.PreserveWorkFiles, "keep-temp", "k", false},
		{&job.Settings.SaveConvertedFile, "save-converted", "sc", false},
		//	S{&job.Settings.SectionSeparator, "sep", "S", "+++"},
		{&job.Settings.StyleSheetName, "style", "s", ""},
		{&job.WorkDirName, "work-dir", "w", ""},
	})
	cmdCheck := NewCommand("check", []Flag{
		//		{&job.Settings.SectionSeparator, "sep", "S", "+++"},
	})
	cmdV1tov2 := NewCommand("v1tov2", []Flag{
		{&job.Settings.ConvertFromVersion, "rosewood-version", "rv", ""},
		{&job.Settings.OverWriteOutputFile, "replace", "r", false},
	})
	cmdInit := NewCommand("init", []Flag{
		{&job.Settings.OverWriteOutputFile, "replace", "r", false},
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
