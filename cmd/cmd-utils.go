// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/drgo/core/args"
	rosewood "github.com/drgo/rosewood/lib"
)

var (
	//Version holds the exe version initialized in the Makefile
	Version string
	//Build holds the exe build number initialized in the Makefile
	Build string
)

func setupCommandFlag(job *rosewood.Job) (flgSets []*flag.FlagSet, err error) {
	args.SetOptions(&args.Options{
		Help: func() {
			helpMessage(nil, getVersion())
		},
	})
	globals := args.NewCommand("", []args.Flag{
		args.NewFlag(&job.RosewoodSettings.Debug, "debug", "d"),
		args.NewFlag(&job.RosewoodSettings.MaxConcurrentWorkers, "max-threads", "mt"),
	})
	baseflags := []args.Flag{ //common to several flags
		args.NewFlag(&job.OverWriteOutputFile, "replace", "r"),
	}
	cmdDo := args.NewCommand("do", []args.Flag{})
	runFlags := append(baseflags, []args.Flag{
		args.NewFlag(&job.ConfigFileName, "config", "cfg"),
		args.NewFlag(&job.RosewoodSettings.ConvertOldVersions, "convert-old", "co"),
		args.NewFlag(&job.RosewoodSettings.ConvertFromVersion, "rosewood-version", "rv"),
		args.NewFlag(&job.RosewoodSettings.DoNotInlineCSS, "no-inlined-css", ""),
		args.NewFlag(&job.OutputFileName, "output", "o"),
		args.NewFlag(&job.RosewoodSettings.PreserveWorkFiles, "keep-temp", "k"),
		args.NewFlag(&job.RosewoodSettings.StyleSheetName, "style", "s"),
		args.NewFlag(&job.WorkDirName, "work-dir", "w"),
	}...)
	cmdRun := args.NewCommand("run", runFlags)
	cmdCheck := args.NewCommand("check", []args.Flag{})
	cmdV1tov2 := args.NewCommand("v1tov2", append(baseflags, []args.Flag{
		args.NewFlag(&job.RosewoodSettings.ConvertFromVersion, "rosewood-version", "rv"),
	}...))
	cmdInit := args.NewCommand("init", runFlags)
	cmdHelp := args.NewCommand("help", []args.Flag{})
	cmdQuery := args.NewCommand("h", []args.Flag{})
	cmdVersion := args.NewCommand("version", []args.Flag{})
	flgSets = append(flgSets, globals, cmdDo, cmdRun, cmdCheck, cmdV1tov2, cmdHelp, cmdQuery, cmdVersion, cmdInit)
	for _, fs := range flgSets {
		fs.Usage = nil               //disable internal usage function
		fs.SetOutput(ioutil.Discard) //suppress output from package flag
	}
	return flgSets, nil
}

func getVersion() string {
	return fmt.Sprintf(versionMessage, Version, Build, rosewood.LibVersion())
}

//crash prints errors and exit with code 2. First line is printed in bold red
func crash(err error) {
	lines := strings.Split(err.Error(), "\n")
	if len(lines) > 0 {
		//"\033[31;1;4m turn color red and bold. \033[0m reset colors"
		fmt.Fprintf(os.Stderr, "\033[31;1m%s\n\033[0m", lines[0])
		for i := 1; i < len(lines); i++ {
			fmt.Fprintf(os.Stderr, "%s\n", lines[i])
		}
	}
	os.Exit(ExitWithError)
}

func crashf(format string, a ...interface{}) {
	err := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, "\033[31;1m%s\n\033[0m", err)
	os.Exit(ExitWithError)
}

func helpMessage(topics []string, versionMessage string) {
	if versionMessage != "" {
		fmt.Fprintf(os.Stderr, versionMessage)
	}
	if len(topics) == 0 {
		fmt.Fprintln(os.Stderr, longUsageMessage)
	}
	for _, topic := range topics {
		switch strings.ToLower(strings.TrimSpace(topic)) {
		case "check":
			fmt.Fprintln(os.Stderr, checkUsageMessage)
		case "run":
			fmt.Fprintln(os.Stderr, runUsageMessage)
		case "v1tov2":
			fmt.Fprintln(os.Stderr, v1tov2UsageMessage)
		case "init":
			fmt.Fprintln(os.Stderr, initUsageMessage)
		case "version":
			fmt.Fprintln(os.Stderr, "prints executable version.")
		case "help":
			fmt.Fprintln(os.Stderr, "you got me! forgot to create help message for help.")
		default:
			fmt.Fprintln(os.Stderr, longUsageMessage)
		}
	}
}

// func usage(topics []string, versionMessage string, exitCode int) {
// 	helpMessage(topics, versionMessage)
// 	if exitCode > -1 {
// 		os.Exit(exitCode)
// 	}
// }
