package main

import (
	"flag"
	"fmt"

	rosewood "github.com/drgo/rosewood/lib"
)

//WARNING: the default value field in flag struct is not used
func setupCommandFlag(settings *rosewood.Settings) (flgSets []*flag.FlagSet, err error) {
	globals := NewCommand("", []Flag{
		{&settings.Debug, "debug", "d", 0},
		{&settings.MaxConcurrentWorkers, "max-threads", "mt", 0},
	})
	cmdRun := NewCommand("run", []Flag{
		{&settings.ConvertOldVersions, "convert-old", "co", false},
		{&settings.ConvertFromVersion, "rosewood-version", "rv", settings.ConvertFromVersion},
		{&settings.DoNotInlineCSS, "no-inlined-css", "", false},
		{&settings.OutputFileName, "output", "o", ""},
		{&settings.OverWriteOutputFile, "replace", "r", false},
		{&settings.PreserveWorkFiles, "keep-temp", "k", false},
		{&settings.SaveConvertedFile, "save-converted", "sc", false},
		//		{&settings.SectionSeparator, "sep", "S", "+++"},
		{&settings.StyleSheetName, "style", "s", ""},
		{&settings.WorkDirName, "work-dir", "w", ""},
	})
	cmdCheck := NewCommand("check", []Flag{
		//		{&settings.SectionSeparator, "sep", "S", "+++"},
	})
	cmdV1tov2 := NewCommand("v1tov2", []Flag{
		{&settings.ConvertFromVersion, "rosewood-version", "rv", ""},
		{&settings.OverWriteOutputFile, "replace", "r", false},
	})
	cmdInit := NewCommand("init", []Flag{
		{&settings.OverWriteOutputFile, "replace", "r", false},
	})
	cmdHelp := NewCommand("help", []Flag{})
	cmdVersion := NewCommand("version", []Flag{})
	flgSets = append(flgSets, globals, cmdRun, cmdCheck, cmdV1tov2, cmdHelp, cmdVersion, cmdInit)
	for _, fs := range flgSets {
		fs.Usage = func() {}    //disable internal usage function
		fs.SetOutput(devNull{}) //suppress output from package flag
	}
	return flgSets, nil
}

func getVersion() string {
	return fmt.Sprintf(versionMessage, Version, Build, rosewood.Version)
}
